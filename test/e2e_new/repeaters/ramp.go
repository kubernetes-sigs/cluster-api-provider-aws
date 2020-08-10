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

func Ramp(items uint64, duration time.Duration, itemFunc func(*sync.WaitGroup)) error {
	if duration < 0 {
		return errors.New("Duration must be >= 0")
	}
	if items == 0 {
		return nil
	}
	if duration == 0 {
		return AtOnce(items, itemFunc)
	}

	interval := duration.Nanoseconds() / int64(items)

	var wg sync.WaitGroup
	for i := uint64(1); i < items; i++ {
		wg.Add(1)
		go itemFunc(&wg)
		time.Sleep(time.Duration(interval))
	}
	wg.Wait()
	return nil
}

func RampConcurrent(rate1, rate2 uint64, forDuration time.Duration, itemFunc func(*sync.WaitGroup)) error {
	if forDuration < 0 {
		return errors.New("Duration must be >= 0")
	}
	if rate1 == 0 {
		return nil
	}
	if forDuration == 0 {
		return AtOnce(rate2, itemFunc)
	}
	deltaRate := (float64(rate2) - float64(rate1)) / forDuration.Seconds()
	rate := float64(rate1)
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for i := uint64(1); i < uint64(rate); i++ {
					wg.Add(1)
					go itemFunc(&wg)
				}
				rate = math.Round(rate + deltaRate)
			}
		}
	}()
	time.Sleep(forDuration)

	wg.Wait()
	return nil
}
