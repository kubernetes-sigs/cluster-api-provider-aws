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

package repeaters

import "sync"

func AtOnce(items uint64, itemFunc func(*sync.WaitGroup)) error {
	var wg sync.WaitGroup
	for i := uint64(1); i < items; i++ {
		wg.Add(1)
		go itemFunc(&wg)
	}
	wg.Wait()
	return nil
}
