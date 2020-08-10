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

import (
	"errors"
	"math"
	"sync"
	"time"
)

func Heaviside(items uint64, duration time.Duration, itemFunc func(*sync.WaitGroup)) error {
	if duration < 0 {
		return errors.New("Duration must be >= 0")
	}
	if items == 0 {
		return nil
	}
	if duration == 0 {
		return AtOnce(items, itemFunc)
	}

	i := heavisideRepeater{
		iteration: 0,
		items:     items,
		iterator:  itemFunc,
	}
	i.t0 = int64(math.Abs(i.heavisideInv(1)))
	i.k = duration.Milliseconds() / (i.t0 * 2)
	return nil
}

type heavisideRepeater struct {
	iteration uint64
	items     uint64
	k         int64
	t0        int64
	iterator  func(*sync.WaitGroup)
}

func (i heavisideRepeater) hasNext() bool {
	return i.iteration < i.items
}

func (i heavisideRepeater) heavisideInv(u uint64) float64 {
	x := float64(u) / float64(i.items+2)
	return math.Erfinv(2*x - 1)
}

func (i heavisideRepeater) next() (time.Duration, error) {
	if !i.hasNext() {
		return 0, errors.New("No new element")
	}
	i.iteration++
	t := int64(math.Abs(i.heavisideInv(i.iteration)))
	newDuration := i.k * (t + i.t0) * int64(time.Millisecond)
	return time.Duration(newDuration), nil
}

func heavisideInv(u uint64, items uint64) float64 {
	x := float64(u) / float64(items+2)
	return math.Erfinv(2*x - 1)
}
