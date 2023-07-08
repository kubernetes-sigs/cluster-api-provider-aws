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

package reqcounter

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	tw "sigs.k8s.io/promo-tools/v3/internal/legacy/timewrapper"
)

// RequestCounter records the number of HTTP requests to GCR.
type RequestCounter struct {
	Mutex     sync.Mutex    // Lock to prevent race-conditions with concurrent processes.
	Requests  uint64        // Number of HTTP requests since recording started.
	Since     time.Time     // When the current request counter began recording requests.
	Interval  time.Duration // The duration of time between each log.
	Threshold uint64        // When to warn of a high request count during a logging cycle. Setting a
	// non-zero threshold allows the request counter to reset each interval. If left uninitialized,
	// the request counter will be persistent and never warn or reset.
}

// increment adds 1 to the request counter, signifying another call to GCR.
func (rc *RequestCounter) Increment() {
	rc.Mutex.Lock()
	rc.Requests++
	rc.Mutex.Unlock()
}

// Flush records the number of HTTP requests found and resets the request counter.
func (rc *RequestCounter) Flush() {
	// Hold onto the lock when reading & writing the request counter.
	rc.Mutex.Lock()
	defer rc.Mutex.Unlock()

	rc.log()

	// Only allow request counters wi
	if rc.Threshold > 0 {
		// Reset the request counter.
		rc.reset()
	}
}

// log the number of HTTP requests found. If the number of requests exceeds the
// threshold, log an additional warning message.
func (rc *RequestCounter) log() {
	msg := fmt.Sprintf("From %s to %s [%d min] there have been %d requests to GCR.", rc.Since.Format(TimestampFormat), Clock.Now().Format(TimestampFormat), rc.Interval/time.Minute, rc.Requests)
	Debug(msg)
	if rc.Threshold > 0 && rc.Requests > rc.Threshold {
		msg = fmt.Sprintf("The threshold of %d requests has been surpassed.", rc.Threshold)
		Warn(msg)
	}
}

// reset clears the request counter and stamps the current time of reset.
func (rc *RequestCounter) reset() {
	rc.Requests = 0
	rc.Since = Clock.Now()
}

// watch indefinitely performs repeated sleep/log cycles.
func (rc *RequestCounter) watch() {
	// TODO: @tylerferrara create a way to cleanly terminate this goroutine.
	go func() {
		for {
			rc.Cycle()
		}
	}()
}

// Cycle sleeps for the request counter's interval and flushes itself.
func (rc *RequestCounter) Cycle() {
	Clock.Sleep(rc.Interval)
	rc.Flush()
}

// RequestCounters holds multiple request counters.
type RequestCounters []*RequestCounter

// NetworkMonitor is the primary means of monitoring network traffic between CIP and GCR.
type NetworkMonitor struct {
	RequestCounters RequestCounters
}

// increment adds 1 to each request counter, signifying a new request has been made to GCR.
func (nm *NetworkMonitor) increment() {
	for _, rc := range nm.RequestCounters {
		rc.Increment()
	}
}

// Log begins logging each request counter at their specified intervals.
func (nm *NetworkMonitor) Log() {
	for _, rc := range nm.RequestCounters {
		rc.watch()
	}
}

const (
	// QuotaWindowShort specifies the length of time to wait before logging in order to estimate the first
	// GCR Quota of 50,000 HTTP requests per 10 min.
	// NOTE: These metrics are only a rough approximation of the actual GCR quotas. The specific 10min measurement
	// is ambiguous, as the start and end time are not specified in the docs. Therefore, it's impossible for our
	// requests counters to perfectly line up with the actual GCR quota.
	// Source: https://cloud.google.com/container-registry/quotas
	QuotaWindowShort time.Duration = time.Minute * 10
	// QuotaWindowLong specifies the length of time to wait before logging in order to estimate the second
	// GCR Quota of 1,000,000 HTTP requests per day.
	QuotaWindowLong time.Duration = time.Hour * 24
	// TimestampFormat specifies the syntax for logging time stamps of request counters.
	TimestampFormat string = "2006-01-02 15:04:05"
)

var (
	// EnableCounting will only become true if the Init function is called. This allows
	// requests to be counted and logged.
	EnableCounting bool
	// NetMonitor holds all request counters for recording HTTP requests to GCR.
	NetMonitor *NetworkMonitor
	// Debug is defined to simplify testing of logrus.Debug calls.
	Debug func(args ...interface{}) = logrus.Debug
	// Warn is defined to simplify testing of logrus.Warn calls.
	Warn func(args ...interface{}) = logrus.Warn
	// Clock is defined to allow mocking of time functions.
	Clock tw.Time = tw.RealTime{}
)

// Init allows request counting to begin.
func Init() {
	EnableCounting = true

	// Create a request counter for logging traffic every 10mins. This aims to mimic the actual
	// GCR quota, but acts as a rough estimation of this quota, indicating when throttling may occur.
	requestCounters := RequestCounters{
		{
			Mutex:     sync.Mutex{},
			Requests:  0,
			Since:     Clock.Now(),
			Interval:  QuotaWindowShort,
			Threshold: 50000,
		},
		{
			Mutex:     sync.Mutex{},
			Requests:  0,
			Since:     Clock.Now(),
			Interval:  QuotaWindowLong,
			Threshold: 1000000,
		},
		{
			Mutex:    sync.Mutex{},
			Requests: 0,
			Since:    Clock.Now(),
			Interval: QuotaWindowShort,
		},
	}

	// Create a new network monitor.
	NetMonitor = &NetworkMonitor{
		RequestCounters: requestCounters,
	}

	// Begin logging network traffic.
	NetMonitor.Log()
}

// Increment increases the all request counters by 1, signifying an HTTP
// request to GCR has been made.
func Increment() {
	if EnableCounting {
		NetMonitor.increment()
	}
}
