/*
Copyright 2019 The Kubernetes Authors.

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

package cli

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	promoter "sigs.k8s.io/promo-tools/v4/promoter/image"
	options "sigs.k8s.io/promo-tools/v4/promoter/image/options"
)

const (
	// flags.
	PromoterManifestFlag                = "manifest"
	PromoterThinManifestDirFlag         = "thin-manifest-dir"
	PromoterSnapshotFlag                = "snapshot"
	PromoterManifestBasedSnapshotOfFlag = "manifest-based-snapshot-of"
	PromoterOutputFlag                  = "output"
)

type Hook struct {
	lastTime time.Time
	mu       sync.RWMutex
}

func NewHook() *Hook {
	return &Hook{
		lastTime: time.Now(),
		mu:       sync.RWMutex{},
	}
}

func (h *Hook) Fire(e *logrus.Entry) error {
	h.mu.Lock()
	e.Data["diff"] = e.Time.Sub(h.lastTime).Round(time.Millisecond)
	h.lastTime = e.Time
	h.mu.Unlock()
	return nil
}

func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func RunPromoteCmd(opts *options.Options) error {
	cip := promoter.New(opts)

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "15:04:05.000",
	})
	logrus.AddHook(NewHook())

	logrus.Infof("Options to check the Signatures: SignCheckIdentity: %s | SignCheckIdentityRegexp: %s | SignCheckIssuer: %s | SignCheckIssuerRegexp: %s",
		opts.SignCheckIdentity, opts.SignCheckIdentityRegexp, opts.SignCheckIssuer, opts.SignCheckIssuerRegexp,
	)

	// Mode 1: Manifest list verification (removed)

	// Mode 2: Snapshots
	if len(opts.Snapshot) > 0 || len(opts.ManifestBasedSnapshotOf) > 0 {
		return cip.Snapshot(opts)
	}

	// Option summary applies to everything except snapshots
	// TODO: Implement if opts.JSONLogSummary { defer sc.LogJSONSummary() }

	// Mode 3: Security scan
	if opts.SeverityThreshold >= 0 {
		return cip.SecurityScan(opts)
	}

	// Mode 4: Image promotion
	if err := cip.PromoteImages(opts); err != nil {
		return fmt.Errorf("promote images: %w", err)
	}

	return nil
}
