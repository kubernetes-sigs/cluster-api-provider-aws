// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elb

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/awserrors"
)

var _ error = &ELBError{}

// ELBError is an error exposed to users of this library.
type ELBError struct {
	err error

	Code int
}

// Error implements the Error interface.
func (e *ELBError) Error() string {
	return e.err.Error()
}

// NewNotFound returns a new error which indicates that the resource of the kind and the name was not found.
func NewNotFound(err error) error {
	return &ELBError{
		err:  err,
		Code: http.StatusNotFound,
	}
}

// NewConflict returns a new error which indicates that the request cannot be processed due to a conflict.
func NewConflict(err error) error {
	return &ELBError{
		err:  err,
		Code: http.StatusConflict,
	}
}

// IsNotFound returns true if the error was created by NewNotFound.
func IsNotFound(err error) bool {
	if ReasonForError(err) == http.StatusNotFound {
		return true
	}
	if code, ok := awserrors.Code(errors.Cause(err)); ok {
		switch code {
		case elb.ErrCodeAccessPointNotFoundException:
			return true
		}
	}
	return false
}

// IsConflict returns true if the error was created by NewConflict.
func IsConflict(err error) bool {
	return ReasonForError(err) == http.StatusConflict
}

// IsSDKError returns true if the error is of type awserr.Error.
func IsSDKError(err error) (ok bool) {
	_, ok = errors.Cause(err).(awserr.Error)
	return
}

// ReasonForError returns the HTTP status for a particular error.
func ReasonForError(err error) int {
	switch t := errors.Cause(err).(type) {
	case *ELBError:
		return t.Code
	}
	return -1
}
