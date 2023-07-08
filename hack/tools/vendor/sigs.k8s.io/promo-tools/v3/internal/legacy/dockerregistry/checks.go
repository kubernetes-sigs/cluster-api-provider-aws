/*
Copyright 2020 The Kubernetes Authors.

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

package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"

	"sigs.k8s.io/promo-tools/v3/internal/legacy/stream"
)

// MKImageVulnCheck returns an instance of ImageVulnCheck which
// checks against images that have known vulnerabilities.
func MKImageVulnCheck(
	syncContext *SyncContext,
	newPullEdges map[PromotionEdge]interface{},
	severityThreshold int,
	fakeVulnProducer ImageVulnProducer,
) *ImageVulnCheck {
	return &ImageVulnCheck{
		syncContext,
		newPullEdges,
		severityThreshold,
		fakeVulnProducer,
	}
}

// Run is a function of ImageVulnCheck and checks that none of the
// images to be promoted have any severe vulnerabilities.
func (check *ImageVulnCheck) Run() error {
	var populateRequests PopulateRequests = func(
		sc *SyncContext,
		reqs chan<- stream.ExternalRequest,
		wg *sync.WaitGroup,
	) {
		srcImages := make(map[PromotionEdge]interface{})
		for edge := range check.PullEdges {
			srcImage := PromotionEdge{
				Digest: edge.Digest,
			}
			// Only check the vulnerability for the source image if it
			// hasn't been checked already since multiple promotion
			// edges can contain the same source image
			if _, found := srcImages[srcImage]; found {
				continue
			}
			srcImages[srcImage] = nil
			var req stream.ExternalRequest
			req.RequestParams = edge
			wg.Add(1)
			reqs <- req
		}
	}

	// If no custom ImageVulnProducer is provided, we use the default producer
	// which is simply a call to the Container Analysis API which lists out
	// the vulnerability occurrences for a given image
	var vulnProducer ImageVulnProducer
	if check.FakeVulnProducer != nil {
		vulnProducer = check.FakeVulnProducer
	} else {
		ctx := context.Background()
		client, err := containeranalysis.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("NewClient: %v", err)
		}
		defer client.Close()
		vulnProducer = mkRealVulnProducer(client)
	}

	vulnerableImages := make([]string, 0)
	var processRequest ProcessRequest = func(
		sc *SyncContext,
		reqs chan stream.ExternalRequest,
		requestResults chan<- RequestResult,
		wg *sync.WaitGroup,
		mutex *sync.Mutex,
	) {
		for req := range reqs {
			reqRes := RequestResult{Context: req}
			errs := make(Errors, 0)
			edge, ok := req.RequestParams.(PromotionEdge)
			if !ok {
				logrus.Errorf("invalid type for promotion edge: %v", edge)
			}

			occurrences, err := vulnProducer(edge)
			if err != nil {
				errs = append(
					errs,
					Error{
						Context: "error getting vulnerabilities",
						Error:   err,
					},
				)
			}

			fixableSevereOccurrences := 0
			for _, occ := range occurrences {
				vuln := occ.GetVulnerability()
				vulnErr := ImageVulnError{
					edge.SrcImageTag.Name,
					edge.Digest,
					occ.GetName(),
					vuln,
				}
				// The vulnerability check should only reject a PR if it finds
				// vulnerabilities that are both fixable and severe
				if vuln.GetFixAvailable() &&
					IsSevereOccurrence(vuln, check.SeverityThreshold) {
					errs = append(errs, Error{
						Context: "Vulnerability Occurrence w/ Fix Available",
						Error:   vulnErr,
					})
					fixableSevereOccurrences++
				} else {
					logrus.Error(vulnErr)
				}
			}

			if fixableSevereOccurrences > 0 {
				vulnerableImages = append(vulnerableImages,
					fmt.Sprintf("%v@%v [%v fixable severe vulnerabilities, "+
						"%v total]",
						edge.SrcImageTag.Name,
						edge.Digest,
						fixableSevereOccurrences,
						len(occurrences)))
			}

			reqRes.Errors = errs
			requestResults <- reqRes
		}
	}

	err := check.SyncContext.ExecRequests(
		populateRequests,
		processRequest,
	)
	if err != nil {
		sort.Strings(vulnerableImages)
		return fmt.Errorf("VulnerabilityCheck: "+
			"The following vulnerable images were found:\n    %v",
			strings.Join(vulnerableImages, "\n    "))
	}
	return nil
}

// Error is a function of ImageVulnError and implements the error interface.
func (err ImageVulnError) Error() string {
	// TODO: Why are we not checking errors here?
	//nolint:errcheck,errchkjson
	vulnJSON, _ := json.MarshalIndent(err, "", "  ")
	return string(vulnJSON)
}

// IsSevereOccurrence checks if a vulnerability is a high enough severity to
// fail the ImageVulnCheck.
func IsSevereOccurrence(
	vuln *grafeaspb.VulnerabilityOccurrence,
	severityThreshold int,
) bool {
	severityLevel := vuln.GetSeverity()
	return int(severityLevel) >= severityThreshold
}

func parseImageProjectID(edge *PromotionEdge) (string, error) {
	const projectIDIndex = 1
	splitName := strings.Split(string(edge.SrcRegistry.Name), "/")
	if len(splitName) <= projectIDIndex {
		return "", fmt.Errorf("could not parse project ID from image name: %q",
			string(edge.SrcRegistry.Name))
	}

	return splitName[projectIDIndex], nil
}

// mkRealVulnProducer returns an ImageVulnProducer that gets all vulnerability
// Occurrences associated with the image represented in the PromotionEdge
// using the Container Analysis Service client library.
func mkRealVulnProducer(client *containeranalysis.Client) ImageVulnProducer {
	return func(
		edge PromotionEdge,
	) ([]*grafeaspb.Occurrence, error) {
		// resourceURL is of the form https://gcr.io/[projectID]/my-image
		resourceURL := "https://" + path.Join(string(edge.SrcRegistry.Name),
			string(edge.SrcImageTag.Name)) + "@" + string(edge.Digest)

		projectID, err := parseImageProjectID(&edge)
		if err != nil {
			return nil, fmt.Errorf("ParsingProjectID: %v", err)
		}

		ctx := context.Background()

		req := &grafeaspb.ListOccurrencesRequest{
			Parent: fmt.Sprintf("projects/%s", projectID),
			Filter: fmt.Sprintf("resourceUrl = %q kind = %q",
				resourceURL, "VULNERABILITY"),
		}

		var occurrenceList []*grafeaspb.Occurrence
		it := client.GetGrafeasClient().ListOccurrences(ctx, req)
		for {
			occ, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("occurrence iteration error: %v", err)
			}
			occurrenceList = append(occurrenceList, occ)
		}

		return occurrenceList, nil
	}
}
