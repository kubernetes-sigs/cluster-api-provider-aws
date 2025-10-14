/*
Copyright (c) 2021 Red Hat, Inc.

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

// This file types and functions useful for testing Prometheus metrics.

package testing

import (
	"io"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	. "github.com/onsi/gomega"       // nolint
	. "github.com/onsi/gomega/ghttp" // nolint
)

// MetricsServer is an HTTP server configured to return Prometheus metrics. Don't create objects of
// this type directly, use the MakeMetricsServer function instead.
type MetricsServer struct {
	server   *Server
	registry *prometheus.Registry
}

// NewMetricsServer creates a metrics server.
func NewMetricsServer() *MetricsServer {
	// Create the registry:
	registry := prometheus.NewPedanticRegistry()

	// Create the server:
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	server := NewServer()
	server.AppendHandlers(handler.ServeHTTP)

	// Create and populate the object:
	return &MetricsServer{
		server:   server,
		registry: registry,
	}
}

// Metrics returns an array of strings containing the metrics available in this server. Each item in
// this array is a line in the Prometheus exposition format. This is intended to be used together
// with the MatchLine matcher.
func (s *MetricsServer) Metrics() []string {
	response, err := http.Get(s.server.URL() + "/metrics")
	Expect(err).ToNot(HaveOccurred())
	defer func() {
		err = response.Body.Close()
		Expect(err).ToNot(HaveOccurred())
	}()
	data, err := io.ReadAll(response.Body)
	Expect(err).ToNot(HaveOccurred())
	return strings.Split(string(data), "\n")
}

// Registry returns the registry that should be used to register metrics for this server.
func (s *MetricsServer) Registry() prometheus.Registerer {
	return s.registry
}

// Close stops the server and releases the resources it uses.
func (s *MetricsServer) Close() {
	s.server.Close()
}

// MatchLine succeeds if actual is an slice of strings that contains at least one items that matches
// the passed regular expression.
func MatchLine(regexp string, args ...interface{}) OmegaMatcher {
	return ContainElement(MatchRegexp(regexp, args...))
}
