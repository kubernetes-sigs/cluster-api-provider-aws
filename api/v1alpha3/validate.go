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

package v1alpha3

import (
	"fmt"
	"net"
	"regexp"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	sshKeyValidNameRegex = regexp.MustCompile(`^[[:graph:]]+([[:print:]]*[[:graph:]]+)*$`)
)

// Validate will validate the bastion fields
func (b *Bastion) Validate() []*field.Error {
	var errs field.ErrorList

	if b.DisableIngressRules && len(b.AllowedCIDRBlocks) > 0 {
		errs = append(errs,
			field.Forbidden(field.NewPath("spec", "bastion", "allowedCIDRBlocks"), "cannot be set if spec.bastion.disableIngressRules is true"),
		)
		return errs
	}

	for i, cidr := range b.AllowedCIDRBlocks {
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			errs = append(errs,
				field.Invalid(field.NewPath("spec", "bastion", fmt.Sprintf("allowedCIDRBlocks[%d]", i)), cidr, "must be a valid CIDR block"),
			)
		}
	}
	return errs
}

func validateSSHKeyName(sshKeyName *string) field.ErrorList {
	var allErrs field.ErrorList
	switch {
	case sshKeyName == nil:
	// nil is accepted
	case sshKeyName != nil && *sshKeyName == "":
	// empty string is accepted
	case sshKeyName != nil && !sshKeyValidNameRegex.Match([]byte(*sshKeyName)):
		allErrs = append(allErrs, field.Invalid(field.NewPath("sshKeyName"), sshKeyName, "Name is invalid. Must be specified in ASCII and must not start or end in whitespace"))
	}
	return allErrs
}

// Validate validates S3Bucket fields.
func (b *S3Bucket) Validate() []*field.Error {
	var errs field.ErrorList

	if !b.Create {
		return errs
	}

	if b.ControlPlaneIAMInstanceProfile == "" {
		errs = append(errs,
			field.Required(field.NewPath("spec", "s3bucket", "controlPlaneIAMInstanceProfiles"), "can't be empty"))
	}

	if len(b.NodesIAMInstanceProfiles) == 0 {
		errs = append(errs,
			field.Required(field.NewPath("spec", "s3bucket", "nodesIAMInstanceProfiles"), "can't be empty"))
	}

	for i, iamInstanceProfile := range b.NodesIAMInstanceProfiles {
		if iamInstanceProfile == "" {
			errs = append(errs,
				field.Required(field.NewPath("spec", "s3bucket", fmt.Sprintf("nodesIAMInstanceProfiles[%d]", i)), "can't be empty"))
		}
	}

	if b.Name != "" {
		errs = append(errs, validateS3BucketName(b.Name)...)
	}

	return errs
}

// Validation rules taken from https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html.
func validateS3BucketName(name string) []*field.Error {
	var errs field.ErrorList

	path := field.NewPath("spec", "s3bucket", "name")

	if len(name) < 3 || len(name) > 63 {
		errs = append(errs, field.Invalid(path, name, "must be between 3 and 63 characters long"))
	}

	s3BucketNotAllowedCharacters := regexp.MustCompile("[^a-z1-9.-]")
	if s3BucketNotAllowedCharacters.MatchString(name) {
		errs = append(errs, field.Invalid(path, name, "consist only of lowercase letters, numbers, dots (.), and hyphens (-)"))
	}

	startsWithNumberOrLetter := regexp.MustCompile("^[a-z1-9]")
	endsWithNumberOrLetter := regexp.MustCompile("[a-z1-9]$")

	if !startsWithNumberOrLetter.MatchString(name) || !endsWithNumberOrLetter.MatchString(name) {
		errs = append(errs, field.Invalid(path, name, "must begin and end with a letter or number"))
	}

	if net.ParseIP(name) != nil {
		errs = append(errs, field.Invalid(path, name, "must not be formatted as an IP address (for example, 192.168.5.4)"))
	}

	return errs
}
