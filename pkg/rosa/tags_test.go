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

package rosa

import (
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	sdk "github.com/openshift-online/ocm-sdk-go"
	ocmlogging "github.com/openshift-online/ocm-sdk-go/logging"
	ocmsdk "github.com/openshift-online/ocm-sdk-go/testing"
	"github.com/openshift/rosa/pkg/ocm"
	"k8s.io/utils/ptr"
)

func TestSupportsDayTwoTags(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{version: "4.20.0", want: true},
		{version: "4.20.0-ec.2", want: true},
		{version: "4.21.3", want: true},
		{version: "5.0.0", want: true},
		{version: "openshift-v4.20.1", want: true},
		{version: "openshift-v4.20.1-candidate", want: true},
		{version: "4.19.12", want: false},
		{version: "4.19.99-rc.1", want: false},
		{version: "", want: false},
		{version: "not-a-version", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(SupportsDayTwoTags(tt.version)).To(Equal(tt.want))
		})
	}
}

func TestTagChanges(t *testing.T) {
	tests := []struct {
		name    string
		current map[string]string
		desired map[string]string
		want    map[string]*string
	}{
		{
			name:    "equal sets need no changes",
			current: map[string]string{"team": "a"},
			desired: map[string]string{"team": "a"},
			want:    nil,
		},
		{
			name:    "nil and empty sets are equal",
			current: map[string]string{},
			desired: nil,
			want:    nil,
		},
		{
			name:    "adds, updates and removes",
			current: map[string]string{"team": "a", "old": "x", "keep": "k"},
			desired: map[string]string{"team": "b", "new": "y", "keep": "k"},
			want:    map[string]*string{"team": ptr.To("b"), "new": ptr.To("y"), "old": nil},
		},
		{
			name:    "clearing removes every key",
			current: map[string]string{"team": "a", "env": "dev"},
			desired: nil,
			want:    map[string]*string{"team": nil, "env": nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(TagChanges(tt.current, tt.desired)).To(Equal(tt.want))
		})
	}
}

func TestNodePoolTagChanges(t *testing.T) {
	tests := []struct {
		name        string
		clusterTags map[string]string
		effective   map[string]string
		desired     map[string]string
		want        map[string]*string
	}{
		{
			name:        "inherited cluster tags aren't written as overrides",
			clusterTags: map[string]string{"env": "prod"},
			effective:   map[string]string{"env": "prod", "team": "a"},
			desired:     map[string]string{"team": "a"},
			want:        nil,
		},
		{
			name:        "desired tag equal to an inherited tag needs no update",
			clusterTags: map[string]string{"env": "prod"},
			effective:   map[string]string{"env": "prod"},
			desired:     map[string]string{"env": "prod"},
			want:        nil,
		},
		{
			name:        "adds a pool override",
			clusterTags: map[string]string{"env": "prod"},
			effective:   map[string]string{"env": "prod"},
			desired:     map[string]string{"team": "a"},
			want:        map[string]*string{"team": ptr.To("a"), "env": nil},
		},
		{
			name:        "overrides an inherited value",
			clusterTags: map[string]string{"env": "prod"},
			effective:   map[string]string{"env": "prod"},
			desired:     map[string]string{"env": "dev"},
			want:        map[string]*string{"env": ptr.To("dev")},
		},
		{
			name:        "removing an override restores inheritance",
			clusterTags: map[string]string{"env": "prod"},
			effective:   map[string]string{"env": "dev", "team": "a"},
			desired:     map[string]string{"team": "a"},
			want:        map[string]*string{"env": nil},
		},
		{
			name:        "re-asserts an override after the cluster tag changes",
			clusterTags: map[string]string{"env": "staging"},
			effective:   map[string]string{"env": "staging"},
			desired:     map[string]string{"env": "prod"},
			want:        map[string]*string{"env": ptr.To("prod")},
		},
		{
			name:        "clearing removes pool overrides",
			clusterTags: nil,
			effective:   map[string]string{"team": "a"},
			desired:     nil,
			want:        map[string]*string{"team": nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(NodePoolTagChanges(tt.clusterTags, tt.effective, tt.desired)).To(Equal(tt.want))
		})
	}
}

func TestTagsDrifted(t *testing.T) {
	g := NewWithT(t)
	actual := map[string]string{"red-hat-managed": "true", "team": "a"}

	g.Expect(TagsDrifted(actual, nil)).To(BeFalse())
	g.Expect(TagsDrifted(actual, map[string]string{"team": "a"})).To(BeFalse())
	g.Expect(TagsDrifted(actual, map[string]string{"team": "b"})).To(BeTrue())
	g.Expect(TagsDrifted(actual, map[string]string{"env": "prod"})).To(BeTrue())
}

func TestTagPatchBody(t *testing.T) {
	g := NewWithT(t)

	body, err := tagPatchBody([]byte(`{"kind":"NodePool","id":"np","aws_node_pool":{"instance_type":"m5.xlarge"}}`),
		nodePoolAWSField, map[string]*string{"team": ptr.To("a"), "old": nil})
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(body).To(MatchJSON(
		`{"kind":"NodePool","id":"np","aws_node_pool":{"instance_type":"m5.xlarge","tags":{"team":"a","old":null}}}`))

	body, err = tagPatchBody([]byte(`{"kind":"Cluster"}`), clusterAWSField, map[string]*string{"team": nil})
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(body).To(MatchJSON(`{"kind":"Cluster","aws":{"tags":{"team":null}}}`))
}

func TestDecodeUserTags(t *testing.T) {
	g := NewWithT(t)

	tags, err := decodeUserTags([]byte(`{"aws":{"tags":{"team":"a","gone":null}}}`), clusterAWSField)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(tags).To(Equal(map[string]string{"team": "a"}))

	tags, err = decodeUserTags([]byte(`{"id":"np"}`), nodePoolAWSField)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(tags).To(BeEmpty())
}

func newTagsTestClient(t *testing.T) (*ocmclient, *ghttp.Server) {
	t.Helper()
	// The OCM SDK test helpers and ghttp verifiers assert through the global Gomega instance.
	RegisterTestingT(t)
	g := NewWithT(t)

	apiServer := ocmsdk.MakeTCPServer()
	t.Cleanup(apiServer.Close)

	logger, err := ocmlogging.NewGoLoggerBuilder().Debug(false).Build()
	g.Expect(err).ToNot(HaveOccurred())
	connection, err := sdk.NewConnectionBuilder().
		Logger(logger).
		Tokens(ocmsdk.MakeTokenString("Bearer", 15*time.Minute)).
		URL(apiServer.URL()).
		Build()
	g.Expect(err).ToNot(HaveOccurred())
	t.Cleanup(func() { _ = connection.Close() })

	return &ocmclient{ocmClient: ocm.NewClientWithConnection(connection), connection: connection}, apiServer
}

func TestOCMClientTags(t *testing.T) {
	t.Run("GetClusterUserTags requests customer-only tags", func(t *testing.T) {
		g := NewWithT(t)
		c, server := newTagsTestClient(t)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters/c1", "fetchUserTagsOnly=true"),
			ocmsdk.RespondWithJSON(http.StatusOK, `{"kind":"Cluster","id":"c1","aws":{"tags":{"team":"a"}}}`),
		))

		tags, err := c.GetClusterUserTags("c1")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(tags).To(Equal(map[string]string{"team": "a"}))
	})

	t.Run("GetNodePoolUserTags reads aws_node_pool tags", func(t *testing.T) {
		g := NewWithT(t)
		c, server := newTagsTestClient(t)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters/c1/node_pools/np1", "fetchUserTagsOnly=true"),
			ocmsdk.RespondWithJSON(http.StatusOK, `{"kind":"NodePool","id":"np1","aws_node_pool":{"tags":{"env":"prod"}}}`),
		))

		tags, err := c.GetNodePoolUserTags("c1", "np1")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(tags).To(Equal(map[string]string{"env": "prod"}))
	})

	t.Run("UpdateClusterTags sends nullable tag changes", func(t *testing.T) {
		g := NewWithT(t)
		c, server := newTagsTestClient(t)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPatch, "/api/clusters_mgmt/v1/clusters/c1"),
			ghttp.VerifyJSON(`{"kind":"Cluster","aws":{"tags":{"team":"b","old":null}}}`),
			ocmsdk.RespondWithJSON(http.StatusOK, `{"kind":"Cluster","id":"c1"}`),
		))

		err := c.UpdateClusterTags("c1", map[string]*string{"team": ptr.To("b"), "old": nil})
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("UpdateNodePoolTags sends nullable tag changes", func(t *testing.T) {
		g := NewWithT(t)
		c, server := newTagsTestClient(t)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodPatch, "/api/clusters_mgmt/v1/clusters/c1/node_pools/np1"),
			ghttp.VerifyJSON(`{"kind":"NodePool","id":"np1","aws_node_pool":{"tags":{"env":null}}}`),
			ocmsdk.RespondWithJSON(http.StatusOK, `{"kind":"NodePool","id":"np1"}`),
		))

		err := c.UpdateNodePoolTags("c1", "np1", map[string]*string{"env": nil})
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("OCM errors keep the API reason", func(t *testing.T) {
		g := NewWithT(t)
		c, server := newTagsTestClient(t)
		server.AppendHandlers(ocmsdk.RespondWithJSON(http.StatusBadRequest,
			`{"kind":"Error","id":"400","code":"CLUSTERS-MGMT-400","reason":"Tag key 'aws:x' is reserved"}`))

		err := c.UpdateClusterTags("c1", map[string]*string{"aws:x": ptr.To("y")})
		g.Expect(err).To(MatchError(ContainSubstring("Tag key 'aws:x' is reserved")))
	})

	t.Run("clients without a raw connection fail", func(t *testing.T) {
		g := NewWithT(t)
		c := &ocmclient{ocmClient: &ocm.Client{}}

		_, err := c.GetClusterUserTags("c1")
		g.Expect(err).To(MatchError(errNoRawConnection))
		g.Expect(c.UpdateNodePoolTags("c1", "np1", nil)).To(MatchError(errNoRawConnection))
	})
}
