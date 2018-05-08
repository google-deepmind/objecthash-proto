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
