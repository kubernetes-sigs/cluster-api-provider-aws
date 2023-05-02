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

package rate

import (
	"context"
	"testing"
)

func TestLimiter_ResetTokens(t *testing.T) {
	lim := NewLimiter(1, 1)
	ctx := context.Background()
	lim.Wait(ctx)
	if lim.tokens != 0.0 {
		t.Errorf("Expected tokens to be 0 after Wait, got %v", lim.tokens)
	}
	lim.tokens = 1.1
	lim.ResetTokens()
	if lim.tokens != 0.0 {
		t.Errorf("Expected tokens to be 0 after ResetTokens, got %v", lim.tokens)
	}
}
