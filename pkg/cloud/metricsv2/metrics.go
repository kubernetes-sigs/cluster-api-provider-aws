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

// Package metricsv2 provides a way to capture request metrics.
package metricsv2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/v2/version"
)

const (
	metricAWSSubsystem       = "aws"
	metricRequestCountKey    = "api_requests_total_v2"
	metricRequestDurationKey = "api_request_duration_seconds_v2"
	metricAPICallRetries     = "api_call_retries_v2"
	metricServiceLabel       = "service"
	metricRegionLabel        = "region"
	metricOperationLabel     = "operation"
	metricControllerLabel    = "controller"
	metricStatusCodeLabel    = "status_code"
	metricErrorCodeLabel     = "error_code"
)

var (
	awsRequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: metricAWSSubsystem,
		Name:      metricRequestCountKey,
		Help:      "Total number of AWS requests",
	}, []string{metricControllerLabel, metricServiceLabel, metricRegionLabel, metricOperationLabel, metricStatusCodeLabel, metricErrorCodeLabel})
	awsRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricAWSSubsystem,
		Name:      metricRequestDurationKey,
		Help:      "Latency of HTTP requests to AWS",
	}, []string{metricControllerLabel, metricServiceLabel, metricRegionLabel, metricOperationLabel})
	awsCallRetries = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricAWSSubsystem,
		Name:      metricAPICallRetries,
		Help:      "Number of retries made against an AWS API",
		Buckets:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}, []string{metricControllerLabel, metricServiceLabel, metricRegionLabel, metricOperationLabel})
	getRawResponse = func(metadata middleware.Metadata) *http.Response {
		switch res := awsmiddleware.GetRawResponse(metadata).(type) {
		case *http.Response:
			return res
		default:
			return nil
		}
	}
)

func init() {
	metrics.Registry.MustRegister(awsRequestCount)
	metrics.Registry.MustRegister(awsRequestDurationSeconds)
	metrics.Registry.MustRegister(awsCallRetries)
}

type requestContextKey struct{}

// RequestData holds information related to request metrics.
type RequestData struct {
	RequestStartTime time.Time
	RequestEndTime   time.Time
	StatusCode       int
	ErrorCode        string
	RequestCount     int
	Service          string
	OperationName    string
	Region           string
	UserAgent        string
	Controller       string
	Target           runtime.Object
	Attempts         int
}

// WithMiddlewares adds instrumentation middleware stacks to AWS GO SDK V2 service clients.
// Inspired by https://github.com/jonathan-innis/aws-sdk-go-prometheus/v2.
func WithMiddlewares(controller string, target runtime.Object) func(stack *middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		if err := stack.Initialize.Add(getMetricCollectionMiddleware(controller, target), middleware.Before); err != nil {
			return err
		}
		if err := stack.Build.Add(getAddToUserAgentMiddleware(), middleware.Before); err != nil {
			return err
		}
		if err := stack.Finalize.Add(getRequestMetricContextMiddleware(), middleware.Before); err != nil {
			return err
		}
		if err := stack.Finalize.Insert(getAttemptContextMiddleware(), "Retry", middleware.After); err != nil {
			return err
		}
		return stack.Deserialize.Add(getRecordAWSPermissionsIssueMiddleware(target), middleware.After)
	}
}

func getMetricCollectionMiddleware(controller string, target runtime.Object) middleware.InitializeMiddleware {
	return middleware.InitializeMiddlewareFunc("capa/MetricCollectionMiddleware", func(ctx context.Context, input middleware.InitializeInput, handler middleware.InitializeHandler) (middleware.InitializeOutput, middleware.Metadata, error) {
		ctx = initRequestContext(ctx, controller, target)
		request := getContext(ctx)

		request.RequestStartTime = time.Now().UTC()
		out, metadata, err := handler.HandleInitialize(ctx, input)
		request.RequestEndTime = time.Now().UTC()

		request.CaptureRequestMetrics()

		return out, metadata, err
	})
}

func getRequestMetricContextMiddleware() middleware.FinalizeMiddleware {
	return middleware.FinalizeMiddlewareFunc("capa/RequestMetricContextMiddleware", func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
		request := getContext(ctx)
		request.Service = awsmiddleware.GetServiceID(ctx)
		request.OperationName = awsmiddleware.GetOperationName(ctx)
		request.Region = awsmiddleware.GetRegion(ctx)

		return handler.HandleFinalize(ctx, input)
	})
}

// For capturing retry count and status codes.
func getAttemptContextMiddleware() middleware.FinalizeMiddleware {
	return middleware.FinalizeMiddlewareFunc("capa/AttemptMetricContextMiddleware", func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
		request := getContext(ctx)
		request.Attempts++
		out, metadata, err := handler.HandleFinalize(ctx, input)
		response := getRawResponse(metadata)

		if response.Body != nil {
			defer response.Body.Close()
		}

		// This will record only last attempts status code.
		// Can be further extended to capture status codes of all attempts
		if response != nil {
			request.StatusCode = response.StatusCode
		} else {
			request.StatusCode = -1
		}

		return out, metadata, err
	})
}
func getRecordAWSPermissionsIssueMiddleware(target runtime.Object) middleware.DeserializeMiddleware {
	return middleware.DeserializeMiddlewareFunc("capa/RecordAWSPermissionsIssueMiddleware", func(ctx context.Context, input middleware.DeserializeInput, handler middleware.DeserializeHandler) (middleware.DeserializeOutput, middleware.Metadata, error) {
		r, ok := input.Request.(*smithyhttp.ResponseError)
		if !ok {
			return middleware.DeserializeOutput{}, middleware.Metadata{}, fmt.Errorf("unknown transport type %T", input.Request)
		}

		var ae smithy.APIError
		if errors.As(r.Err, &ae) {
			switch ae.ErrorCode() {
			case "AuthFailure", "UnauthorizedOperation", "NoCredentialProviders":
				record.Warnf(target, ae.ErrorCode(), "Operation %s failed with a credentials or permission issue", awsmiddleware.GetOperationName(ctx))
			}
		}
		return handler.HandleDeserialize(ctx, input)
	})
}

func getAddToUserAgentMiddleware() middleware.BuildMiddleware {
	return middleware.BuildMiddlewareFunc("capa/AddUserAgentMiddleware", func(ctx context.Context, input middleware.BuildInput, handler middleware.BuildHandler) (middleware.BuildOutput, middleware.Metadata, error) {
		request := getContext(ctx)
		r, ok := input.Request.(*smithyhttp.Request)
		if !ok {
			return middleware.BuildOutput{}, middleware.Metadata{}, fmt.Errorf("unknown transport type %T", input.Request)
		}

		if curUA := r.Header.Get("User-Agent"); curUA != "" {
			request.UserAgent = curUA + " " + request.UserAgent
		}
		r.Header.Set("User-Agent", request.UserAgent)

		return handler.HandleBuild(ctx, input)
	})
}

func initRequestContext(ctx context.Context, controller string, target runtime.Object) context.Context {
	if middleware.GetStackValue(ctx, requestContextKey{}) == nil {
		ctx = middleware.WithStackValue(ctx, requestContextKey{}, &RequestData{
			Controller: controller,
			Target:     target,
			UserAgent:  fmt.Sprintf("aws.cluster.x-k8s.io/%s", version.Get().String()),
		})
	}
	return ctx
}

func getContext(ctx context.Context) *RequestData {
	rctx := middleware.GetStackValue(ctx, requestContextKey{})
	if rctx == nil {
		return nil
	}
	return rctx.(*RequestData)
}

// CaptureRequestMetrics will monitor and capture request metrics.
func (r *RequestData) CaptureRequestMetrics() {
	requestDuration := r.RequestStartTime.Sub(r.RequestEndTime)
	retryCount := r.Attempts - 1

	awsRequestCount.WithLabelValues(r.Controller, r.Service, r.Region, r.OperationName, strconv.Itoa(r.StatusCode), r.ErrorCode).Inc()
	awsRequestDurationSeconds.WithLabelValues(r.Controller, r.Service, r.Region, r.OperationName).Observe(requestDuration.Seconds())
	awsCallRetries.WithLabelValues(r.Controller, r.Service, r.Region, r.OperationName).Observe(float64(retryCount))
}
