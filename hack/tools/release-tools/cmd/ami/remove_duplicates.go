/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ami

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"

	duplicates "sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/internal/ami"
	"sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/printer"
)

// defaultAMIOwnerID is the account CAPA's own AMI build pipeline
// (.github/workflows/build-ami.yml, role gh-image-builder) publishes AMIs
// into.
const defaultAMIOwnerID = "819546954734"

// defaultAMIRegions are the AWS regions CAPA publishes AMIs to.
var defaultAMIRegions = []string{
	"ap-south-1", "eu-west-3", "eu-west-2", "eu-west-1",
	"ap-northeast-2", "ap-northeast-1", "sa-east-1", "ca-central-1",
	"ap-southeast-1", "ap-southeast-2", "eu-central-1",
	"us-east-1", "us-east-2", "us-west-1", "us-west-2",
}

// ec2API is the subset of the EC2 client this command depends on. It exists
// so tests can inject a fake instead of making real AWS calls; *ec2.Client
// satisfies it.
type ec2API interface {
	DescribeImages(ctx context.Context, params *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error)
	DeregisterImage(ctx context.Context, params *ec2.DeregisterImageInput, optFns ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error)
}

// ec2ClientFactory returns an ec2API client configured for region.
type ec2ClientFactory func(ctx context.Context, region string) (ec2API, error)

// removeDuplicatesCmd returns the `remove-duplicates` cobra command.
func removeDuplicatesCmd() *cobra.Command {
	var (
		region             string
		ownerID            string
		output             string
		deleteFlag         bool
		includeUngroupable bool
	)

	cmd := &cobra.Command{
		Use:   "remove-duplicates",
		Short: "Find and remove duplicate AMIs owned by the CAPA AMI build pipeline",
		Long: `Find AMIs that share the same distribution, distribution_version and
kubernetes_version tags within a region. For each such group, the AMI with
the highest build_timestamp tag is kept and the rest are reported as
duplicates.

Only AMIs tagged with kubernetes_version are considered, and only if they
are x86_64, available, and hvm - the properties every CAPA AMI build has.
This keeps unrelated images in the same AWS account out of scope entirely.

By default this command is dry-run: it reports every AMI (kept, duplicate,
and ungroupable) and previews what removal would do, without changing
anything. Pass --delete to actually deregister the AMIs and delete their
backing snapshots.

Only DUPLICATE AMIs are removed by default. Pass --include-ungroupable to
also remove AMIs that couldn't be evaluated for duplication (missing tags).

By default all CAPA regions are searched; use --region to restrict to one
or more (comma-separated).`,
		Example: `
  # Dry run across all CAPA regions
  release-tool ami remove-duplicates

  # Dry run in a single region, JSON report
  release-tool ami remove-duplicates --region us-west-2 -o json

  # Actually remove duplicate AMIs
  release-tool ami remove-duplicates --delete

  # Also remove AMIs that couldn't be evaluated for duplication
  release-tool ami remove-duplicates --delete --include-ungroupable`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return assertNoDuplicateFlags(os.Args,
				[]string{"region"},
				[]string{"owner-id"},
				[]string{"output", "o"},
			)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			regions := defaultAMIRegions
			if requested := splitVersions(region); len(requested) > 0 {
				regions = requested
			}

			amis, err := describeOwnedAMIs(cmd.Context(), newEC2Client, regions, ownerID)
			if err != nil {
				return err
			}

			entries := duplicates.FindDuplicateAMIs(amis).Entries()

			p, err := printer.New(output, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			if err := p.Print(entriesPrinterInput(output, entries)); err != nil {
				return err
			}

			targets := removalTargets(entries, includeUngroupable)
			failed := removeAMIs(cmd.Context(), cmd.OutOrStdout(), newEC2Client, targets, !deleteFlag)
			return removalError(failed)
		},
	}

	cmd.Flags().StringVar(&region, "region", "",
		"Comma-separated list of AWS regions to search (default: all CAPA regions)")
	cmd.Flags().StringVar(&ownerID, "owner-id", defaultAMIOwnerID,
		"AWS account ID that owns the AMIs to search")
	cmd.Flags().StringVarP(&output, "output", "o", string(printer.TypeTable),
		"Output format: table, json, or yaml")
	cmd.Flags().BoolVar(&deleteFlag, "delete", false,
		"Actually remove the AMIs (default is dry-run)")
	cmd.Flags().BoolVar(&includeUngroupable, "include-ungroupable", false,
		"Also remove AMIs that couldn't be evaluated for duplication (default: only DUPLICATE)")

	return cmd
}

// removalError turns a removeAMIs failure count into an error, so the
// command exits non-zero when any AMI failed to be removed.
func removalError(failed int) error {
	if failed > 0 {
		return fmt.Errorf("%d AMI(s) failed to remove", failed)
	}
	return nil
}

// removalTargets returns the entries eligible for removal: always
// StatusDuplicate, plus StatusUngroupable when includeUngroupable is set.
// StatusKeep is never a target.
func removalTargets(entries []duplicates.Entry, includeUngroupable bool) []duplicates.Entry {
	var targets []duplicates.Entry
	for _, e := range entries {
		switch e.Status {
		case duplicates.StatusDuplicate:
			targets = append(targets, e)
		case duplicates.StatusUngroupable:
			if includeUngroupable {
				targets = append(targets, e)
			}
		}
	}
	return targets
}

// removeAMIs previews (dryRun) or performs removal of each target AMI,
// deregistering it and deleting its backing snapshots in one call. Failures
// are printed and skipped rather than aborting the run; the number of
// failures is returned so the caller can decide whether to exit non-zero.
// One client is created per distinct target region and reused across its
// targets.
func removeAMIs(ctx context.Context, w io.Writer, newClient ec2ClientFactory, targets []duplicates.Entry, dryRun bool) int {
	if len(targets) == 0 {
		fmt.Fprintln(w, "\nNo AMIs to remove")
		return 0
	}

	verb := "Would remove"
	if !dryRun {
		verb = "Removing"
	}
	fmt.Fprintf(w, "\n%s %d AMI(s):\n", verb, len(targets))

	clients := make(map[string]ec2API)
	failed := 0
	for _, target := range targets {
		label := fmt.Sprintf("%s %s (%s)", target.Region, target.AMI.ID, target.Status)

		if dryRun {
			fmt.Fprintf(w, "  [DRY RUN] %s\n", label)
			continue
		}

		client, ok := clients[target.Region]
		if !ok {
			var err error
			client, err = newClient(ctx, target.Region)
			if err != nil {
				fmt.Fprintf(w, "  ERROR %s: %v\n", label, err)
				failed++
				continue
			}
			clients[target.Region] = client
		}

		_, err := client.DeregisterImage(ctx, &ec2.DeregisterImageInput{
			ImageId:                   aws.String(target.AMI.ID),
			DeleteAssociatedSnapshots: aws.Bool(true),
		})
		if err != nil {
			fmt.Fprintf(w, "  ERROR %s: %v\n", label, err)
			failed++
			continue
		}
		fmt.Fprintf(w, "  Removed %s\n", label)
	}

	if dryRun {
		fmt.Fprintln(w, "\nDry run — nothing removed. Pass --delete to remove them.")
	}
	return failed
}

// newEC2Client is the production ec2ClientFactory: it returns a real EC2
// client configured for region.
func newEC2Client(ctx context.Context, region string) (ec2API, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("loading AWS config for region %q: %w", region, err)
	}
	return ec2.NewFromConfig(cfg), nil
}

// describeOwnedAMIs fetches all available x86_64/hvm AMIs owned by ownerID
// across regions, tagged with their originating region. The tag-key filter
// restricts results to images carrying the kubernetes_version tag that
// every CAPA AMI build sets, so images unrelated to the CAPA pipeline (but
// owned by the same account) are never returned in the first place.
func describeOwnedAMIs(ctx context.Context, newClient ec2ClientFactory, regions []string, ownerID string) ([]duplicates.AMI, error) {
	var amis []duplicates.AMI

	for _, region := range regions {
		client, err := newClient(ctx, region)
		if err != nil {
			return nil, err
		}

		out, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{
			Filters: []types.Filter{
				{Name: aws.String("owner-id"), Values: []string{ownerID}},
				{Name: aws.String("tag-key"), Values: []string{duplicates.TagKubernetesVersion}},
				{Name: aws.String("architecture"), Values: []string{"x86_64"}},
				{Name: aws.String("state"), Values: []string{"available"}},
				{Name: aws.String("virtualization-type"), Values: []string{"hvm"}},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("describing images in region %q: %w", region, err)
		}

		for _, image := range out.Images {
			amis = append(amis, duplicates.AMI{
				Region: region,
				ID:     aws.ToString(image.ImageId),
				Name:   aws.ToString(image.Name),
				Tags:   tagsToMap(image.Tags),
			})
		}
	}

	return amis, nil
}

func tagsToMap(tags []types.Tag) map[string]string {
	if len(tags) == 0 {
		return nil
	}
	m := make(map[string]string, len(tags))
	for _, tag := range tags {
		m[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	return m
}

// entriesReport wraps a flattened entry list for json/yaml output, keeping
// the top-level shape a named object rather than a bare array.
type entriesReport struct {
	Entries []duplicates.Entry `json:"entries"`
}

// entriesPrinterInput returns a *printer.Table for table output or an
// entriesReport for json/yaml.
func entriesPrinterInput(format string, entries []duplicates.Entry) interface{} {
	normalized := printer.Type(format)
	if normalized == printer.TypeTable || normalized == "" {
		return entriesToTable(entries)
	}
	return entriesReport{Entries: entries}
}

// entriesToTable converts a flattened entry list into a *printer.Table.
func entriesToTable(entries []duplicates.Entry) *printer.Table {
	rows := make([][]string, 0, len(entries))
	for _, e := range entries {
		groupKey := e.GroupKey
		if groupKey == "" {
			groupKey = "-"
		}
		reasonOrTimestamp := e.BuildTimestamp
		if e.Status == duplicates.StatusUngroupable {
			reasonOrTimestamp = e.Reason
		}
		rows = append(rows, []string{e.Region, groupKey, string(e.Status), e.AMI.ID, reasonOrTimestamp})
	}
	return &printer.Table{
		Columns: []string{"REGION", "GROUP", "STATUS", "AMI ID", "BUILD TIMESTAMP / REASON"},
		Rows:    rows,
	}
}
