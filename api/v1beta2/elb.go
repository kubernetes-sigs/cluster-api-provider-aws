/*
Copyright 2022 The Kubernetes Authors.

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
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	re = regexp.MustCompile(`^(TCP|SSL|HTTP|HTTPS|TLS|UDP):(\d+)\/?(\w*)$`)
)

func (e *ELBHealthCheckTarget) ValidateELBHealthCheck() []*field.Error {
	if e == nil {
		return nil
	}
	var errs field.ErrorList
	healthCheck := e.String()
	val := re.MatchString(healthCheck)
	if !val {
		errs = append(errs, field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "healthCheck"), healthCheck, "value cannot have characters other than alphabets, numbers, and ':' ."))
	}
	if !validateHealthCheckWithPath(healthCheck) {
		errs = append(errs,
			field.Invalid(field.NewPath("spec", "controlPlaneLoadBalancer", "healthCheck"), healthCheck, "cannot specify paths with protocol other than HTTP and HTTPS"),
		)
	}
	return errs
}

func validateHealthCheckWithPath(healthCheck string) bool {
	protocol := strings.Split(healthCheck, ":")[0]
	if strings.Contains(healthCheck, "/") &&
		!(protocol == ELBProtocolHTTP.String() || protocol == ELBProtocolHTTPS.String()) {
		return false
	}
	return true
}
