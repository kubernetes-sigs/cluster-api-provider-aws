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
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	duplicates "sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/internal/ami"
	"sigs.k8s.io/cluster-api-provider-aws/hack/tools/release-tools/printer"
)

// fakeEC2 is a minimal ec2API stub for tests.
type fakeEC2 struct {
	describeImagesFunc  func(ctx context.Context, params *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error)
	deregisterImageFunc func(ctx context.Context, params *ec2.DeregisterImageInput, optFns ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error)
}

func (f *fakeEC2) DescribeImages(ctx context.Context, params *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error) {
	if f.describeImagesFunc != nil {
		return f.describeImagesFunc(ctx, params, optFns...)
	}
	return &ec2.DescribeImagesOutput{}, nil
}

func (f *fakeEC2) DeregisterImage(ctx context.Context, params *ec2.DeregisterImageInput, optFns ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error) {
	if f.deregisterImageFunc != nil {
		return f.deregisterImageFunc(ctx, params, optFns...)
	}
	return &ec2.DeregisterImageOutput{}, nil
}

// factoryFor returns an ec2ClientFactory that always returns client.
func factoryFor(client ec2API) ec2ClientFactory {
	return func(_ context.Context, _ string) (ec2API, error) {
		return client, nil
	}
}

func TestRemovalTargets(t *testing.T) {
	entries := []duplicates.Entry{
		{AMI: duplicates.AMI{ID: "keep-1"}, Status: duplicates.StatusKeep},
		{AMI: duplicates.AMI{ID: "dup-1"}, Status: duplicates.StatusDuplicate},
		{AMI: duplicates.AMI{ID: "ungroupable-1"}, Status: duplicates.StatusUngroupable},
	}

	t.Run("default targets only DUPLICATE", func(t *testing.T) {
		targets := removalTargets(entries, false)
		if len(targets) != 1 || targets[0].AMI.ID != "dup-1" {
			t.Fatalf("targets = %+v, want only dup-1", targets)
		}
	})

	t.Run("include-ungroupable adds UNGROUPABLE but never KEEP", func(t *testing.T) {
		targets := removalTargets(entries, true)
		if len(targets) != 2 {
			t.Fatalf("got %d targets, want 2", len(targets))
		}
		for _, tgt := range targets {
			if tgt.Status == duplicates.StatusKeep {
				t.Errorf("KEEP entry %q became a removal target", tgt.AMI.ID)
			}
		}
	})
}

func TestRemoveAMIsDryRun(t *testing.T) {
	fake := &fakeEC2{
		deregisterImageFunc: func(context.Context, *ec2.DeregisterImageInput, ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error) {
			t.Fatal("DeregisterImage must not be called during a dry run")
			return nil, nil
		},
	}
	targets := []duplicates.Entry{{Region: "us-east-1", AMI: duplicates.AMI{ID: "ami-1"}, Status: duplicates.StatusDuplicate}}

	var buf bytes.Buffer
	failed := removeAMIs(context.Background(), &buf, factoryFor(fake), targets, true)

	if failed != 0 {
		t.Errorf("failed = %d, want 0", failed)
	}
	if !strings.Contains(buf.String(), "[DRY RUN]") {
		t.Errorf("output missing [DRY RUN]: %q", buf.String())
	}
}

func TestRemoveAMIsDeletesWithSnapshots(t *testing.T) {
	var gotInput *ec2.DeregisterImageInput
	fake := &fakeEC2{
		deregisterImageFunc: func(_ context.Context, params *ec2.DeregisterImageInput, _ ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error) {
			gotInput = params
			return &ec2.DeregisterImageOutput{}, nil
		},
	}
	targets := []duplicates.Entry{{Region: "us-east-1", AMI: duplicates.AMI{ID: "ami-1"}, Status: duplicates.StatusDuplicate}}

	var buf bytes.Buffer
	failed := removeAMIs(context.Background(), &buf, factoryFor(fake), targets, false)

	if failed != 0 {
		t.Fatalf("failed = %d, want 0", failed)
	}
	if gotInput == nil {
		t.Fatal("DeregisterImage was not called")
	}
	if aws.ToString(gotInput.ImageId) != "ami-1" {
		t.Errorf("ImageId = %q, want ami-1", aws.ToString(gotInput.ImageId))
	}
	if gotInput.DeleteAssociatedSnapshots == nil || !*gotInput.DeleteAssociatedSnapshots {
		t.Errorf("DeleteAssociatedSnapshots = %v, want true", gotInput.DeleteAssociatedSnapshots)
	}
}

func TestRemoveAMIsReusesClientPerRegion(t *testing.T) {
	fake := &fakeEC2{}
	calls := map[string]int{}
	factory := func(_ context.Context, region string) (ec2API, error) {
		calls[region]++
		return fake, nil
	}
	targets := []duplicates.Entry{
		{Region: "us-east-1", AMI: duplicates.AMI{ID: "ami-1"}, Status: duplicates.StatusDuplicate},
		{Region: "us-east-1", AMI: duplicates.AMI{ID: "ami-2"}, Status: duplicates.StatusDuplicate},
		{Region: "eu-west-1", AMI: duplicates.AMI{ID: "ami-3"}, Status: duplicates.StatusDuplicate},
	}

	var buf bytes.Buffer
	if failed := removeAMIs(context.Background(), &buf, factory, targets, false); failed != 0 {
		t.Fatalf("failed = %d, want 0", failed)
	}

	if calls["us-east-1"] != 1 {
		t.Errorf("us-east-1 client created %d time(s), want 1", calls["us-east-1"])
	}
	if calls["eu-west-1"] != 1 {
		t.Errorf("eu-west-1 client created %d time(s), want 1", calls["eu-west-1"])
	}
}

func TestRemoveAMIsPartialFailure(t *testing.T) {
	var deregistered []string
	fake := &fakeEC2{
		deregisterImageFunc: func(_ context.Context, params *ec2.DeregisterImageInput, _ ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error) {
			id := aws.ToString(params.ImageId)
			if id == "ami-bad" {
				return nil, errors.New("boom")
			}
			deregistered = append(deregistered, id)
			return &ec2.DeregisterImageOutput{}, nil
		},
	}
	targets := []duplicates.Entry{
		{Region: "us-east-1", AMI: duplicates.AMI{ID: "ami-bad"}, Status: duplicates.StatusDuplicate},
		{Region: "us-east-1", AMI: duplicates.AMI{ID: "ami-good"}, Status: duplicates.StatusDuplicate},
	}

	var buf bytes.Buffer
	failed := removeAMIs(context.Background(), &buf, factoryFor(fake), targets, false)

	if failed != 1 {
		t.Fatalf("failed = %d, want 1", failed)
	}
	if len(deregistered) != 1 || deregistered[0] != "ami-good" {
		t.Errorf("deregistered = %v, want [ami-good] (a failure must not stop the remaining targets)", deregistered)
	}
	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("output missing an ERROR line: %q", buf.String())
	}

	if err := removalError(failed); err == nil {
		t.Error("removalError(1) = nil, want a non-nil error")
	}
}

func TestRemovalError(t *testing.T) {
	if err := removalError(0); err != nil {
		t.Errorf("removalError(0) = %v, want nil", err)
	}
	if err := removalError(2); err == nil {
		t.Error("removalError(2) = nil, want an error")
	}
}

func TestDescribeOwnedAMIsFilters(t *testing.T) {
	var gotFilters []types.Filter
	fake := &fakeEC2{
		describeImagesFunc: func(_ context.Context, params *ec2.DescribeImagesInput, _ ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error) {
			gotFilters = params.Filters
			return &ec2.DescribeImagesOutput{}, nil
		},
	}

	if _, err := describeOwnedAMIs(context.Background(), factoryFor(fake), []string{"us-east-1"}, "819546954734"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	byName := make(map[string][]string, len(gotFilters))
	for _, f := range gotFilters {
		byName[aws.ToString(f.Name)] = f.Values
	}

	tagKeyValues, ok := byName["tag-key"]
	if !ok {
		t.Fatal("DescribeImages filters are missing a tag-key filter to scope discovery to CAPA-tagged AMIs")
	}
	if len(tagKeyValues) != 1 || tagKeyValues[0] != duplicates.TagKubernetesVersion {
		t.Errorf("tag-key filter values = %v, want [%s]", tagKeyValues, duplicates.TagKubernetesVersion)
	}
}

func TestEntriesToTable(t *testing.T) {
	entries := []duplicates.Entry{
		{Region: "us-east-1", GroupKey: "ubuntu|22.04|v1.36.3", Status: duplicates.StatusKeep, AMI: duplicates.AMI{ID: "ami-keep"}, BuildTimestamp: "200"},
		{Region: "us-east-1", GroupKey: "ubuntu|22.04|v1.36.3", Status: duplicates.StatusDuplicate, AMI: duplicates.AMI{ID: "ami-dup"}, BuildTimestamp: "100"},
		{Region: "us-east-1", Status: duplicates.StatusUngroupable, AMI: duplicates.AMI{ID: "ami-ungroupable"}, Reason: `missing "distribution" tag`},
	}

	table := entriesToTable(entries)

	if len(table.Rows) != 3 {
		t.Fatalf("got %d rows, want 3", len(table.Rows))
	}
	if got := table.Rows[0]; got[1] != "ubuntu|22.04|v1.36.3" || got[2] != "KEEP" || got[4] != "200" {
		t.Errorf("KEEP row = %v", got)
	}
	if got := table.Rows[2]; got[1] != "-" || got[2] != "UNGROUPABLE" || got[4] != `missing "distribution" tag` {
		t.Errorf("UNGROUPABLE row = %v", got)
	}
}

func TestEntriesPrinterInput(t *testing.T) {
	entries := []duplicates.Entry{{AMI: duplicates.AMI{ID: "ami-1"}, Status: duplicates.StatusKeep}}

	if _, ok := entriesPrinterInput("table", entries).(*printer.Table); !ok {
		t.Error("table format did not return *printer.Table")
	}
	if _, ok := entriesPrinterInput("", entries).(*printer.Table); !ok {
		t.Error("empty format did not default to *printer.Table")
	}

	got, ok := entriesPrinterInput("json", entries).(entriesReport)
	if !ok {
		t.Fatal("json format did not return entriesReport")
	}
	if len(got.Entries) != 1 {
		t.Errorf("entriesReport has %d entries, want 1", len(got.Entries))
	}
}
