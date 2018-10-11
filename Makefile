# Copyright 2018 The Kubernetes Authors.
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

GOFLAGS += -ldflags '-extldflags "-static"'
GOREBUILD :=
GOPATH := $(shell go env GOPATH)

.PHONY: gendeepcopy

all: generate build

ifndef FASTBUILD
# If FASTBUILD isn't defined, fully rebuild Go binaries, and always
# run dep ensure
GOREBUILD += -a
.PHONY: vendor
endif

vendor:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure

depend-update:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure -update

generate: gendeepcopy

gendeepcopy: vendor
	go build -o $(GOPATH)/bin/deepcopy-gen sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/code-generator/cmd/deepcopy-gen
	$(GOPATH)/bin/deepcopy-gen \
	  -i ./pkg/cloud/aws/providerconfig,./pkg/cloud/aws/providerconfig/v1alpha1 \
	  -O zz_generated.deepcopy \
	  -h boilerplate.go.txt

genmocks: vendor
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API" "cloud/aws/services/ec2/mock_ec2iface/mock.go"
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/elb/elbiface ELBAPI" "cloud/aws/services/elb/mock_elbiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 MachineInterface" "cloud/aws/actuators/machine/mock_machineiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 ClusterInterface" "cloud/aws/actuators/cluster/mock_clusteriface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services EC2Interface" "cloud/aws/services/mocks/ec2.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services ELBInterface" "cloud/aws/services/mocks/elb.go"

build: clusterctl-bin

clusterctl-bin: vendor
	CGO_ENABLED=0 go install $(GOFLAGS) $(GOREBUILD) sigs.k8s.io/cluster-api-provider-aws/clusterctl

clusterawsadm-bin: vendor
	CGO_ENABLED=0 go install $(GOFLAGS) $(GOREBUILD) sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm

cluster-api-dev-helper-bin: vendor
	CGO_ENABLED=0 go install $(GOFLAGS) sigs.k8s.io/cluster-api-provider-aws/hack/cluster-api-dev-helper

images: vendor
	# TODO: manager image here

dev_push:
	# TODO: manager image here

check: fmt vet

test: vendor
	go test -race -cover ./cmd/... ./pkg/... ./clusterctl/...

fmt: vendor
	hack/verify-gofmt.sh

vet: vendor
	go vet ./...

lint:
	golint || go get -u golang.org/x/lint/golint
	golint -set_exit_status ./cmd/... ./pkg/... ./clusterctl/...

examples = clusterctl/examples/aws/out/cluster.yaml clusterctl/examples/aws/out/machines.yaml clusterctl/examples/aws/out/provider-components.yaml
templates = clusterctl/examples/aws/cluster.yaml.template clusterctl/examples/aws/machines.yaml.template clusterctl/examples/aws/provider-components.yaml.template
example: $(examples)
$(examples) : envfile $(templates)
	source ./envfile && cd ./clusterctl/examples/aws && ./generate-yaml.sh

envfile:
	# create the envfile and exit if the envfile doesn't already exist
	cp -n envfile.example envfile
	echo "\033[0;31mPlease fill out your envfile!\033[0m"
	exit 1

clean:
	rm -rf clusterctl/examples/aws/out
