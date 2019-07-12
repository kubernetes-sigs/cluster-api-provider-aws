# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build the manager binary
FROM golang:1.12.5 as builder

ARG ARCH

# Copy in the go src
WORKDIR $GOPATH/src/sigs.k8s.io/cluster-api
COPY pkg/    pkg/
COPY vendor/ vendor/
COPY cmd/    cmd/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -a -ldflags '-extldflags "-static"' -o ./cmd/example-provider/manager sigs.k8s.io/cluster-api/cmd/example-provider

# Copy the controller-manager into a thin image
FROM gcr.io/distroless/static:latest
WORKDIR /
COPY --from=builder /go/src/sigs.k8s.io/cluster-api/cmd/example-provider/manager .
ENTRYPOINT ["/manager"]
