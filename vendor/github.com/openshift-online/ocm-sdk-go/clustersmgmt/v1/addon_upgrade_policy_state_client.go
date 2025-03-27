/*
Copyright (c) 2020 Red Hat, Inc.

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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// AddonUpgradePolicyStateClient is the client of the 'addon_upgrade_policy_state' resource.
//
// Manages a specific upgrade policy state.
type AddonUpgradePolicyStateClient struct {
	transport http.RoundTripper
	path      string
}

// NewAddonUpgradePolicyStateClient creates a new client for the 'addon_upgrade_policy_state'
// resource using the given transport to send the requests and receive the
// responses.
func NewAddonUpgradePolicyStateClient(transport http.RoundTripper, path string) *AddonUpgradePolicyStateClient {
	return &AddonUpgradePolicyStateClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the upgrade policy state.
func (c *AddonUpgradePolicyStateClient) Get() *AddonUpgradePolicyStateGetRequest {
	return &AddonUpgradePolicyStateGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Update the upgrade policy state.
func (c *AddonUpgradePolicyStateClient) Update() *AddonUpgradePolicyStateUpdateRequest {
	return &AddonUpgradePolicyStateUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AddonUpgradePolicyStatePollRequest is the request for the Poll method.
type AddonUpgradePolicyStatePollRequest struct {
	request    *AddonUpgradePolicyStateGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AddonUpgradePolicyStatePollRequest) Parameter(name string, value interface{}) *AddonUpgradePolicyStatePollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AddonUpgradePolicyStatePollRequest) Header(name string, value interface{}) *AddonUpgradePolicyStatePollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AddonUpgradePolicyStatePollRequest) Interval(value time.Duration) *AddonUpgradePolicyStatePollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AddonUpgradePolicyStatePollRequest) Status(value int) *AddonUpgradePolicyStatePollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AddonUpgradePolicyStatePollRequest) Predicate(value func(*AddonUpgradePolicyStateGetResponse) bool) *AddonUpgradePolicyStatePollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AddonUpgradePolicyStateGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AddonUpgradePolicyStatePollRequest) StartContext(ctx context.Context) (response *AddonUpgradePolicyStatePollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AddonUpgradePolicyStatePollResponse{
			response: result.(*AddonUpgradePolicyStateGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AddonUpgradePolicyStatePollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AddonUpgradePolicyStatePollResponse is the response for the Poll method.
type AddonUpgradePolicyStatePollResponse struct {
	response *AddonUpgradePolicyStateGetResponse
}

// Status returns the response status code.
func (r *AddonUpgradePolicyStatePollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AddonUpgradePolicyStatePollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AddonUpgradePolicyStatePollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AddonUpgradePolicyStatePollResponse) Body() *AddonUpgradePolicyState {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddonUpgradePolicyStatePollResponse) GetBody() (value *AddonUpgradePolicyState, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AddonUpgradePolicyStateClient) Poll() *AddonUpgradePolicyStatePollRequest {
	return &AddonUpgradePolicyStatePollRequest{
		request: c.Get(),
	}
}

// AddonUpgradePolicyStateGetRequest is the request for the 'get' method.
type AddonUpgradePolicyStateGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AddonUpgradePolicyStateGetRequest) Parameter(name string, value interface{}) *AddonUpgradePolicyStateGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonUpgradePolicyStateGetRequest) Header(name string, value interface{}) *AddonUpgradePolicyStateGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonUpgradePolicyStateGetRequest) Impersonate(user string) *AddonUpgradePolicyStateGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonUpgradePolicyStateGetRequest) Send() (result *AddonUpgradePolicyStateGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonUpgradePolicyStateGetRequest) SendContext(ctx context.Context) (result *AddonUpgradePolicyStateGetResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "GET",
		URL:    uri,
		Header: header,
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &AddonUpgradePolicyStateGetResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	err = readAddonUpgradePolicyStateGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonUpgradePolicyStateGetResponse is the response for the 'get' method.
type AddonUpgradePolicyStateGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddonUpgradePolicyState
}

// Status returns the response status code.
func (r *AddonUpgradePolicyStateGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonUpgradePolicyStateGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonUpgradePolicyStateGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddonUpgradePolicyStateGetResponse) Body() *AddonUpgradePolicyState {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddonUpgradePolicyStateGetResponse) GetBody() (value *AddonUpgradePolicyState, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AddonUpgradePolicyStateUpdateRequest is the request for the 'update' method.
type AddonUpgradePolicyStateUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AddonUpgradePolicyState
}

// Parameter adds a query parameter.
func (r *AddonUpgradePolicyStateUpdateRequest) Parameter(name string, value interface{}) *AddonUpgradePolicyStateUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonUpgradePolicyStateUpdateRequest) Header(name string, value interface{}) *AddonUpgradePolicyStateUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonUpgradePolicyStateUpdateRequest) Impersonate(user string) *AddonUpgradePolicyStateUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *AddonUpgradePolicyStateUpdateRequest) Body(value *AddonUpgradePolicyState) *AddonUpgradePolicyStateUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonUpgradePolicyStateUpdateRequest) Send() (result *AddonUpgradePolicyStateUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonUpgradePolicyStateUpdateRequest) SendContext(ctx context.Context) (result *AddonUpgradePolicyStateUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAddonUpgradePolicyStateUpdateRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "PATCH",
		URL:    uri,
		Header: header,
		Body:   io.NopCloser(buffer),
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &AddonUpgradePolicyStateUpdateResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	err = readAddonUpgradePolicyStateUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonUpgradePolicyStateUpdateResponse is the response for the 'update' method.
type AddonUpgradePolicyStateUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddonUpgradePolicyState
}

// Status returns the response status code.
func (r *AddonUpgradePolicyStateUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonUpgradePolicyStateUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonUpgradePolicyStateUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddonUpgradePolicyStateUpdateResponse) Body() *AddonUpgradePolicyState {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddonUpgradePolicyStateUpdateResponse) GetBody() (value *AddonUpgradePolicyState, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
