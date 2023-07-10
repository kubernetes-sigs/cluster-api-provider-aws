/*
Copyright 2022 The Kubernetes Authors.

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

package imagepromoter

import (
	"fmt"

	reg "sigs.k8s.io/promo-tools/v3/internal/legacy/dockerregistry"
	options "sigs.k8s.io/promo-tools/v3/promoter/image/options"
)

// ScanEdges runs the vulnerability scans on the new images
// detected by the promoter.
func (di *DefaultPromoterImplementation) ScanEdges(
	opts *options.Options, sc *reg.SyncContext,
	promotionEdges map[reg.PromotionEdge]interface{},
) error {
	if err := sc.RunChecks(
		[]reg.PreCheck{
			reg.MKImageVulnCheck(
				sc,
				promotionEdges,
				opts.SeverityThreshold,
				nil,
			),
		},
	); err != nil {
		return fmt.Errorf("checking image vulnerabilities: %w", err)
	}
	di.PrintSection("END (VULNSCAN)", opts.Confirm)
	return nil
}
