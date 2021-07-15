// +build e2e

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

package shared

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	. "github.com/onsi/gomega"

	"github.com/gofrs/flock"
	"sigs.k8s.io/yaml"
)

type TestResource struct {
	EC2       int `json:"ec2"`
	VPC       int `json:"vpc"`
	EIP       int `json:"eip"`
	IGW       int `json:"igw"`
	NGW       int `json:"ngw"`
	ClassicLB int `json:"classiclb"`
}

func WriteResourceQuotesToFile(logPath string, serviceQuotas map[string]*ServiceQuota) {
	if _, err := os.Stat(logPath); err == nil {
		// If resource-quotas file exists, remove it. Should not fail on error, another ginkgo node might have already deleted it.
		os.Remove(logPath)
	}

	resources := TestResource{
		EC2:       serviceQuotas["ec2"].Value,
		VPC:       serviceQuotas["vpc"].Value,
		EIP:       serviceQuotas["eip"].Value,
		IGW:       serviceQuotas["igw"].Value,
		NGW:       serviceQuotas["ngw"].Value,
		ClassicLB: serviceQuotas["classiclb"].Value,
	}
	data, err := yaml.Marshal(resources)
	Expect(err).NotTo(HaveOccurred())

	err = ioutil.WriteFile(logPath, data, 0644)
	Expect(err).NotTo(HaveOccurred())
}

func (r *TestResource) String() string {
	return fmt.Sprintf("{ec2:%v, vpc:%v, eip:%v, ngw:%v, igw:%v, classiclb:%v}", r.EC2, r.VPC, r.EIP, r.NGW, r.IGW, r.ClassicLB)
}

func (r *TestResource) WriteRequestedResources(e2eCtx *E2EContext, testName string) {
	requestedResourceFilePath := path.Join(e2eCtx.Settings.ArtifactFolder, "requested-resources.yaml")
	if _, err := os.Stat(ResourceQuotaFilePath); err != nil {
		// If requested-resources file does not exist, create it
		f, err := os.Create(requestedResourceFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(f.Close()).NotTo(HaveOccurred())
	}

	fileLock := flock.New(requestedResourceFilePath)
	defer func() error {
		if err := fileLock.Unlock(); err != nil {
			return err
		}
		return nil
	}()

	err := fileLock.Lock()
	Expect(err).NotTo(HaveOccurred())

	requestedResources, err := ioutil.ReadFile(requestedResourceFilePath)
	Expect(err).NotTo(HaveOccurred())

	resources := struct {
		TestResourceMap map[string]TestResource `json:"requested-resources"`
	}{}
	err = yaml.Unmarshal(requestedResources, &resources)
	Expect(err).NotTo(HaveOccurred())

	if resources.TestResourceMap == nil {
		resources.TestResourceMap = make(map[string]TestResource, 0)
	}
	resources.TestResourceMap[testName] = *r
	str, err := yaml.Marshal(resources)
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(requestedResourceFilePath, str, 0644)).To(Succeed())
}

func (r *TestResource) doesSatisfy(request *TestResource) bool {
	if request.EC2 != 0 && r.EC2 < request.EC2 {
		return false
	}
	if request.IGW != 0 && r.IGW < request.IGW {
		return false
	}
	if request.NGW != 0 && r.NGW < request.NGW {
		return false
	}
	if request.ClassicLB != 0 && r.ClassicLB < request.ClassicLB {
		return false
	}
	if request.VPC != 0 && r.VPC < request.VPC {
		return false
	}
	if request.EIP != 0 && r.EIP < request.EIP {
		return false
	}

	return true
}

func (r *TestResource) acquire(request *TestResource) {
	r.EC2 -= request.EC2
	r.VPC -= request.VPC
	r.EIP -= request.EIP
	r.NGW -= request.NGW
	r.IGW -= request.IGW
	r.ClassicLB -= request.ClassicLB
}

func (r *TestResource) release(request *TestResource) {
	r.EC2 += request.EC2
	r.VPC += request.VPC
	r.EIP += request.EIP
	r.NGW += request.NGW
	r.IGW += request.IGW
	r.ClassicLB += request.ClassicLB
}

func AcquireResources(request *TestResource, nodeNum int, fileLock *flock.Flock) error {
	timeoutAfter := time.Now().Add(time.Hour * 6)
	defer func() error {
		if err := fileLock.Unlock(); err != nil {
			return err
		}
		return nil
	}()

	Byf("Node %d acquiring resources: %s", nodeNum, request.String())
	for range time.Tick(time.Second) {
		if time.Now().After(timeoutAfter) {
			Byf("Timeout reached for node %d", nodeNum)
			break
		}
		err := fileLock.Lock()
		if err != nil {
			continue
		}
		resourceText, err := ioutil.ReadFile(ResourceQuotaFilePath)
		if err != nil {
			return err
		}

		resources := &TestResource{}
		if err = yaml.Unmarshal(resourceText, resources); err != nil {
			return err
		}

		if resources.doesSatisfy(request) {
			resources.acquire(request)
			data, err := yaml.Marshal(resources)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(ResourceQuotaFilePath, data, 0644); err != nil {
				return err
			}
			Byf("Node %d acquired resources: %s", nodeNum, request.String())
			return nil
		}
		e2eDebugBy("Insufficient resources, retrying")
		if err := fileLock.Unlock(); err != nil {
			return err
		}
	}
	return errors.New("giving up on acquiring resource due to timeout")
}

func e2eDebugBy(msg string) {
	if os.Getenv("E2E_DEBUG") != "" {
		Byf(msg)
	}
}

func ReleaseResources(request *TestResource, nodeNum int, fileLock *flock.Flock) error {
	timeoutInSec := 20

	defer fileLock.Unlock()
	var tryCount = 0
	for range time.Tick(1 * time.Second) {
		tryCount++
		if tryCount > timeoutInSec {
			break
		}
		if err := fileLock.Lock(); err != nil {
			continue
		}
		resourceText, err := ioutil.ReadFile(ResourceQuotaFilePath)
		if err != nil {
			return err
		}
		resources := &TestResource{}
		if err := yaml.Unmarshal(resourceText, resources); err != nil {
			return err
		}
		resources.release(request)
		data, err := yaml.Marshal(resources)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(ResourceQuotaFilePath, data, 0644); err != nil {
			return err
		}
		Byf("Node %d released resources: %s", nodeNum, request.String())
		return nil
	}
	return errors.New("giving up on releasing resource due to timeout")
}
