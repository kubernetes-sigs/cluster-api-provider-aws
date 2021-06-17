/*
Copyright 2021 The Kubernetes Authors.

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

package util

import (
	"fmt"
	"os"
)

// ErrEnvironmentVariableNotFound is an error string for environment variable not found error.
type ErrEnvironmentVariableNotFound string

// Error defines the error interface for ErrEnvironmentVariableNotFound.
func (e ErrEnvironmentVariableNotFound) Error() string {
	return fmt.Sprintf("environment variable %q not found", string(e))
}

// GetEnv will lookup and return an environment variable.
func GetEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", ErrEnvironmentVariableNotFound(key)
	}
	return val, nil
}
