// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package instrumentation

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/glog"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
)

var (
	AWSOperation, _ = tag.NewKey("aws.operation")
	AWSService, _   = tag.NewKey("aws.service")
	AWSRegion, _    = tag.NewKey("aws.region")
)

func awsRequestAction(r *http.Request) string {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Info("could not read request body")
		return ""
	}
	q, err := url.ParseQuery(string(b))
	if err != nil {
		glog.Info("could not parse request body")
		return ""
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return q.Get("Action")
}

func AWSInstrumentingHTTPClient() *http.Client {
	// Set up instrumentation

	tracer := func(r *http.Request, t *trace.Span) *httptrace.ClientTrace {

		action := awsRequestAction(r)
		components := strings.Split(r.URL.Hostname(), ".")

		region := components[0]
		service := components[1]

		ctx, _ := tag.New(r.Context(),
			tag.Upsert(AWSOperation, action),
			tag.Upsert(AWSRegion, region),
			tag.Upsert(AWSService, service),
		)

		r = r.WithContext(ctx)
		return &httptrace.ClientTrace{}
	}

	httpClient := &http.Client{
		Transport: &ochttp.Transport{
			NewClientTrace: tracer,
		},
	}

	return httpClient
}

func AWSInstrumentedConfig() *aws.Config {

	awsConfig := aws.NewConfig()

	awsConfig.WithHTTPClient(AWSInstrumentingHTTPClient())

	return awsConfig
}

// MethodName creates a trace name
func MethodName(paths ...string) string {
	arr := []string{"sigs.k8s.io", "cluster-api-provider-aws", "cloud", "aws"}

	arr = append(arr, paths...)
	return strings.Join(arr, "/")
}

// NewCounter register a new counter
func NewCounter(desc string, paths ...string) *stats.Int64Measure {
	paths = append(paths, "count")
	id := MethodName(paths...)
	c := stats.Int64(id, desc, stats.UnitDimensionless)

	if err := view.Register(&view.View{
		Name:        id,
		Description: desc,
		Measure:     c,
		Aggregation: view.Sum(),
	}); err != nil {
		glog.Fatalf("Cannot register view: %v", err)
	}
	return c
}

func RegisterHTTPViews() {
	//view.Register(ochttp.DefaultClientViews...)

	if err := view.Register(
		&view.View{
			Name:        "opencensus.io/http/client/request_count",
			TagKeys:     []tag.Key{AWSOperation},
			Measure:     ochttp.ClientRequestCount,
			Aggregation: view.Count(),
		},
	); err != nil {
		glog.Fatal(err)
	}
}
