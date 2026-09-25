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

package v1beta2

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func TestNotificationContactsJSONRoundTrip(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	codecs := serializer.NewCodecFactory(scheme)

	tests := []struct {
		name         string
		contacts     []string
		wantInJSON   bool
		wantNil      bool
		wantContacts []string
	}{
		{
			name:       "unset is omitted and decodes as nil",
			contacts:   nil,
			wantInJSON: false,
			wantNil:    true,
		},
		{
			name:         "explicitly empty is preserved",
			contacts:     []string{},
			wantInJSON:   true,
			wantNil:      false,
			wantContacts: []string{},
		},
		{
			name:         "populated is preserved",
			contacts:     []string{"user-a", "user-b"},
			wantInJSON:   true,
			wantNil:      false,
			wantContacts: []string{"user-a", "user-b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			in := &ROSAControlPlane{
				Spec: RosaControlPlaneSpec{NotificationContacts: tt.contacts},
			}
			in.SetGroupVersionKind(GroupVersion.WithKind("ROSAControlPlane"))

			// encoding/json round trip.
			data, err := json.Marshal(in)
			g.Expect(err).NotTo(HaveOccurred())

			var raw struct {
				Spec map[string]json.RawMessage `json:"spec"`
			}
			g.Expect(json.Unmarshal(data, &raw)).To(Succeed())
			_, found := raw.Spec["notificationContacts"]
			g.Expect(found).To(Equal(tt.wantInJSON))

			out := &ROSAControlPlane{}
			g.Expect(json.Unmarshal(data, out)).To(Succeed())
			g.Expect(out.Spec.NotificationContacts == nil).To(Equal(tt.wantNil))
			if !tt.wantNil {
				g.Expect(out.Spec.NotificationContacts).To(Equal(tt.wantContacts))
			}

			// Kubernetes codec round trip, as used by typed clients.
			codec := codecs.LegacyCodec(GroupVersion)
			encoded, err := runtime.Encode(codec, in)
			g.Expect(err).NotTo(HaveOccurred())

			decoded := &ROSAControlPlane{}
			_, _, err = codecs.UniversalDeserializer().Decode(encoded, nil, decoded)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(decoded.Spec.NotificationContacts == nil).To(Equal(tt.wantNil))
			if !tt.wantNil {
				g.Expect(decoded.Spec.NotificationContacts).To(Equal(tt.wantContacts))
			}
		})
	}
}
