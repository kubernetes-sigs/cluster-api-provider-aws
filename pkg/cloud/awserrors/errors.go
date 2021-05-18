/*
Copyright 2018 The Kubernetes Authors.

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

package awserrors

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	AuthFailure                = "AuthFailure"
	InUseIPAddress             = "InvalidIPAddress.InUse"
	GroupNotFound              = "InvalidGroup.NotFound"
	PermissionNotFound         = "InvalidPermission.NotFound"
	VPCNotFound                = "InvalidVpcID.NotFound"
	SubnetNotFound             = "InvalidSubnetID.NotFound"
	InternetGatewayNotFound    = "InvalidInternetGatewayID.NotFound"
	NATGatewayNotFound         = "InvalidNatGatewayID.NotFound"
	GatewayNotFound            = "InvalidGatewayID.NotFound"
	EIPNotFound                = "InvalidElasticIpID.NotFound"
	RouteTableNotFound         = "InvalidRouteTableID.NotFound"
	LoadBalancerNotFound       = "LoadBalancerNotFound"
	ResourceNotFound           = "InvalidResourceID.NotFound"
	InvalidSubnet              = "InvalidSubnet"
	AssociationIDNotFound      = "InvalidAssociationID.NotFound"
	InvalidInstanceID          = "InvalidInstanceID.NotFound"
	LaunchTemplateNameNotFound = "InvalidLaunchTemplateName.NotFoundException"
	ResourceExists             = "ResourceExistsException"
	NoCredentialProviders      = "NoCredentialProviders"
)

var _ error = &EC2Error{}

// Code returns the AWS error code as a string
func Code(err error) (string, bool) {
	if awserr, ok := err.(awserr.Error); ok {
		return awserr.Code(), true
	}
	return "", false
}

// Message returns the AWS error message as a string
func Message(err error) string {
	if awserr, ok := err.(awserr.Error); ok {
		return awserr.Message()
	}
	return ""
}

// EC2Error is an error exposed to users of this library.
type EC2Error struct {
	msg string

	Code int
}

// Error implements the Error interface.
func (e *EC2Error) Error() string {
	return e.msg
}

// NewNotFound returns an error which indicates that the resource of the kind and the name was not found.
func NewNotFound(msg string) error {
	return &EC2Error{
		msg:  msg,
		Code: http.StatusNotFound,
	}
}

// NewConflict returns an error which indicates that the request cannot be processed due to a conflict.
func NewConflict(msg string) error {
	return &EC2Error{
		msg:  msg,
		Code: http.StatusConflict,
	}
}

func IsResourceExists(err error) bool {
	if code, ok := Code(err); ok {
		return code == ResourceExists
	}
	return false
}

// NewFailedDependency returns an error which indicates that a dependency failure status
func NewFailedDependency(msg string) error {
	return &EC2Error{
		msg:  msg,
		Code: http.StatusFailedDependency,
	}
}

// IsFailedDependency checks if the error is pf http.StatusFailedDependency
func IsFailedDependency(err error) bool {
	return ReasonForError(err) == http.StatusFailedDependency
}

// IsNotFound returns true if the error was created by NewNotFound.
func IsNotFound(err error) bool {
	if ReasonForError(err) == http.StatusNotFound {
		return true
	}
	return IsInvalidNotFoundError(err)
}

// IsConflict returns true if the error was created by NewConflict.
func IsConflict(err error) bool {
	return ReasonForError(err) == http.StatusConflict
}

// IsSDKError returns true if the error is of type awserr.Error.
func IsSDKError(err error) (ok bool) {
	_, ok = err.(awserr.Error)
	return
}

// IsInvalidNotFoundError tests for common aws not found errors
func IsInvalidNotFoundError(err error) bool {
	if code, ok := Code(err); ok {
		switch code {
		case VPCNotFound:
			return true
		case InvalidInstanceID:
			return true
		case ssm.ErrCodeParameterNotFound:
			return true
		case LaunchTemplateNameNotFound:
			return true
		}
	}

	return false
}

// ReasonForError returns the HTTP status for a particular error.
func ReasonForError(err error) int {
	if t, ok := err.(*EC2Error); ok {
		return t.Code
	}

	return -1
}

// IsIgnorableSecurityGroupError checks for errors in SG that can be ignored and then return nil.
func IsIgnorableSecurityGroupError(err error) error {
	if code, ok := Code(err); ok {
		switch code {
		case GroupNotFound, PermissionNotFound:
			return nil
		default:
			return err
		}
	}
	return nil
}
