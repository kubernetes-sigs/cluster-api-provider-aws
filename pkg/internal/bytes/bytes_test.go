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

package bytes

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestSplitBytes(t *testing.T) {
	g := NewWithT(t)

	t.Run("should call 1 time if it fits", func(t *testing.T) {
		maxSize := 100
		input := make([]byte, 100)
		_, err := crand.Read(input)
		g.Expect(err).To(BeNil())

		count := 0
		Split(input, false, maxSize, func(split []byte) {
			g.Expect(split).To(BeEquivalentTo(input))
			count++
		})
		g.Expect(count).To(BeEquivalentTo(1))
	})

	t.Run("should properly encode and split given random input and maxsize", func(t *testing.T) {
		maxSize := 1 + rand.Intn(1024)
		input := make([]byte, rand.Intn(24576))
		_, err := crand.Read(input)
		g.Expect(err).To(BeNil())

		data := []byte{}
		Split(input, true, maxSize, func(split []byte) {
			data = append(data, split...)
		})

		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		l, err := base64.StdEncoding.Decode(decoded, data)
		g.Expect(err).To(BeNil())

		decoded = decoded[:l]

		g.Expect(decoded).To(BeEquivalentTo(input))
	})

	t.Run("should properly split given random input and maxsize", func(t *testing.T) {
		maxSize := 1 + rand.Intn(1024)
		input := make([]byte, rand.Intn(24576))
		_, err := crand.Read(input)
		g.Expect(err).To(BeNil())

		expected := len(input) / maxSize
		if math.Mod(float64(len(input)), float64(maxSize)) > 0 {
			// Add 1 to expected if there is remaining bytes left at the end of all the splits.
			expected++
		}

		data := []byte{}
		count := 0
		Split(input, false, maxSize, func(split []byte) {
			data = append(data, split...)
			count++
		})

		g.Expect(data).To(BeEquivalentTo(input))
		g.Expect(expected).To(BeEquivalentTo(count), fmt.Sprintf("input=%d, maxsize=%d", len(input), maxSize))
	})

	t.Run("should call 0 times if there is no data", func(t *testing.T) {
		maxSize := 100
		input := []byte{}

		count := 0
		Split(input, false, maxSize, func(split []byte) {
			t.Fail()
		})
		g.Expect(count).To(BeEquivalentTo(0))
	})
}
