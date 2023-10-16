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

package logclient

import (
	"context"
	"log"

	"cloud.google.com/go/logging"
)

// GcpLogClient is a GCP log client.
type GcpLogClient struct {
	logClient *logging.Client
	loggers   [3]*log.Logger
}

// Close simply calls Close() to the underlying logging client (from which the
// child loggers are derived).
func (c *GcpLogClient) Close() error {
	return c.logClient.Close()
}

// GetInfoLogger exposes the internal Info logger.
func (c *GcpLogClient) GetInfoLogger() *log.Logger {
	return c.loggers[IndexLogInfo]
}

// GetErrorLogger exposes the internal Error logger.
func (c *GcpLogClient) GetErrorLogger() *log.Logger {
	return c.loggers[IndexLogError]
}

// GetAlertLogger exposes the internal Alert logger.
func (c *GcpLogClient) GetAlertLogger() *log.Logger {
	return c.loggers[IndexLogAlert]
}

// NewGcpLogClient returns a new LoggingFacility that logs to GCP resources. As
// such, it requires the GCP projectID as well as the logName to log to.
func NewGcpLogClient(
	projectID, logName string,
) (*GcpLogClient, error) {
	c := GcpLogClient{}
	ctx := context.Background()

	// This creates a logging client that performs better logging than the
	// default behavior on GCP Stackdriver. For instance, logs sent with this
	// client are not split up over newlines, and also the severity levels are
	// actually understood by Stackdriver.
	logClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	c.logClient = logClient

	c.loggers[IndexLogInfo] = logClient.
		Logger(logName).StandardLogger(logging.Info)
	c.loggers[IndexLogError] = logClient.
		Logger(logName).StandardLogger(logging.Error)
	c.loggers[IndexLogAlert] = logClient.
		Logger(logName).StandardLogger(logging.Alert)

	return &c, nil
}
