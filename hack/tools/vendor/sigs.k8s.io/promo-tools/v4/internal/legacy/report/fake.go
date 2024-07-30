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

package report

import (
	"bytes"
	"log"

	"cloud.google.com/go/errorreporting"
)

// FakeReportingClient is a fake reporting client.
type FakeReportingClient struct {
	reporter     *log.Logger
	reportBuffer bytes.Buffer
}

// GetReportBuffer retrieves the reportBuffer.
func (c *FakeReportingClient) GetReportBuffer() bytes.Buffer {
	return c.reportBuffer
}

// Report simply prints the entry to STDERR. Nothing goes over the network!
func (c *FakeReportingClient) Report(
	e errorreporting.Entry,
) {
	c.reporter.Println(e)
}

// Close is a NOP (there is nothing to close).
func (c *FakeReportingClient) Close() error { return nil }

// NewFakeReportingClient creates a new FakeReportingClient that has the
// reporter initialized to a basic logger to STDERR.
func NewFakeReportingClient() *FakeReportingClient {
	c := FakeReportingClient{}
	c.reporter = log.New(&c.reportBuffer, "FAKE-REPORT", log.LstdFlags)

	return &c
}
