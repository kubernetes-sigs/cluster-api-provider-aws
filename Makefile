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

.PHONY: gendeepcopy

all: generate build images

depend:
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure

depend-update:
	dep ensure -update

generate: gendeepcopy

gendeepcopy: depend
	go build -o $$GOPATH/bin/deepcopy-gen sigs.k8s.io/cluster-api-provider-aws/vendor/k8s.io/code-generator/cmd/deepcopy-gen
	$$GOPATH/bin/deepcopy-gen \
	  -i ./cloud/aws/providerconfig,./cloud/aws/providerconfig/v1alpha1 \
	  -O zz_generated.deepcopy \
	  -h boilerplate.go.txt

genmocks: depend
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API" "cloud/aws/services/ec2/mock_ec2iface/mock.go"
	hack/generate-mocks.sh "github.com/aws/aws-sdk-go/service/elb/elbiface ELBAPI" "cloud/aws/services/elb/mock_elbiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 MachineInterface" "cloud/aws/actuators/machine/mock_machineiface/mock.go"
	hack/generate-mocks.sh "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1 ClusterInterface" "cloud/aws/actuators/cluster/mock_clusteriface/mock.go"

build: depend clusterctl-bin cluster-controller machine-controller

clusterctl-bin:
	CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' sigs.k8s.io/cluster-api-provider-aws/clusterctl

images: depend
	$(MAKE) -C cmd/cluster-controller image
	$(MAKE) -C cmd/machine-controller image

dev_push: depend cluster-controller-dev-push machine-controller-dev-push

cluster-controller:
	CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' sigs.k8s.io/cluster-api-provider-aws/cmd/cluster-controller

cluster-controller-dev-push: cluster-controller
	$(MAKE) -C cmd/cluster-controller dev_push

machine-controller:
	CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' sigs.k8s.io/cluster-api-provider-aws/cmd/machine-controller

machine-controller-dev-push: machine-controller
	$(MAKE) -C cmd/machine-controller dev_push

push: depend
	$(MAKE) -C cmd/cluster-controller push
	$(MAKE) -C cmd/machine-controller push

check: depend fmt vet

test: depend
	go test -race -cover ./cmd/... ./cloud/... ./clusterctl/...

fmt:
	hack/verify-gofmt.sh

vet:
	go vet ./...

examples = clusterctl/examples/aws/out/cluster.yaml clusterctl/examples/aws/out/machines.yaml clusterctl/examples/aws/out/provider-components.yaml
templates = clusterctl/examples/aws/cluster.yaml.template clusterctl/examples/aws/machines.yaml.template clusterctl/examples/aws/provider-components.yaml.template
example: $(examples)
$(examples) : envfile $(templates)
	source ./envfile && cd ./clusterctl/examples/aws && ./generate-yaml.sh

envfile: envfile.example
	# create the envfile and exit if the envfile doesn't already exist
	cp -n envfile.example envfile
	echo "\033[0;31mPlease fill out your envfile!\033[0m"
	exit 1

clean:
	rm -rf clusterctl/examples/aws/out
