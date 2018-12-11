// Copyright 2018 The ObjectHash-Proto Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protohash

import (
	"fmt"
	"reflect"
	"testing"
)

// TestHashTimestampWithBadInputs tests how hashTimestamp handles bad inputs.
func TestHashTimestampWithBadInputs(t *testing.T) {
	hasher := objectHasher{}

	badTimestampValues := []reflect.Value{
		// Not a struct.
		reflect.ValueOf(0.0),

		// Not a valid Timestamp struct (fields have the wrong type).
		reflect.ValueOf(struct {
			Seconds float64
			Nanos   float64
		}{}),
	}

	for i, v := range badTimestampValues {
		t.Run(fmt.Sprintf("TestBadTimestamps-%d", i), func(t *testing.T) {
			_, err := hasher.hashTimestamp(v)
			if err == nil {
				t.Errorf("Attempting to hash %T{ %[1]v } as a timestamp should have returned an error.", v)
			}
		})
	}
}

func TestHashDuration(t *testing.T) {
	durationValue := reflect.ValueOf(struct {
			Seconds int32
			Nanos   int32
		}{})
	hasher := objectHasher{}
	hashedValue, err := hasher.hashDuration(durationValue)
	if err != nil {
		t.Errorf("Hashing duration %v should not error", durationValue)
	}
	if fmt.Sprintf("%x", hashedValue) != "3a82b649344529f03f52c1833f5aecc488a53b31461a1f54c305d149b12b8f53" {
		t.Errorf("Duration %T{ %[1]v } should incorrectly hashed to %v", durationValue, fmt.Sprintf("%x", hashedValue))
	}
}

func TestHashDurationsWithBadInputs(t *testing.T) {
	hasher := objectHasher{}

	badDurationValues := []reflect.Value{
		// Not a struct.
		reflect.ValueOf(0.0),

		// Not a valid Timestamp struct (fields have the wrong type).
		reflect.ValueOf(struct {
			Seconds float64
			Nanos   float64
		}{}),
	}

	for i, v := range badDurationValues {
		t.Run(fmt.Sprintf("TestHashDurations-%d", i), func(t *testing.T) {
			_, err := hasher.hashTimestamp(v)
			if err == nil {
				t.Errorf("Attempting to hash %T{ %[1]v } as a duration should have returned an error.", v)
			}
		})
	}
}