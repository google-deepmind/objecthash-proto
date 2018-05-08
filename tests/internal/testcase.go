// Copyright 2017 The ObjectHash-Proto Authors
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

package internal

import (
	"fmt"
	"testing"

	"github.com/benlaurie/objecthash/go/objecthash"
	oi "github.com/deepmind/objecthash-proto/internal"
	"github.com/golang/protobuf/proto"
)

// TestCase represents typical ObjectHash-Proto test cases.
type TestCase struct {
	// Protos is a list of protobuf messages that should have the same objecthash.
	Protos []proto.Message

	// EquivalentObject is a plain Go object that should have the same objecthash
	// as the messages under the `protos` field.
	EquivalentObject interface{}

	// EquivalentJSONString is a JSON object that should have the same objecthash
	// as the messages under the `protos` field.
	EquivalentJSONString string

	// ExpectedHashString is the expected objecthash for all the objects in the
	// test case.
	ExpectedHashString string
}

// Check tests the ObjectHashes for the protos in a TestCase's Protos field.
//
// It does the following checks:
// - The ObjectHashes of the protos (stringified) are equal to the ExpectedHashString.
// - The ObjectHashes of the protos are equal to the ObjectHash of the EquivalentJSONString, if present.
// - The ObjectHashes of the protos are equal to the ObjectHash of the EquivalentObject, if present.
func (tc TestCase) Check(t *testing.T, hasher oi.ProtoHasher) {
	t.Helper()

	for _, message := range tc.Protos {
		messageHash, err := hasher.HashProto(message)
		if err != nil {
			t.Errorf("Attempting to hash %T{ %[1]v } returned an error: %v", message, err)
		}
		messageHashStr := fmt.Sprintf("%x", messageHash)

		// If the test case has an expected hash string, check it.
		if tc.ExpectedHashString != "" {
			t.Run("Compare to expected hash", func(t *testing.T) {
				if messageHashStr != tc.ExpectedHashString {
					t.Errorf("Got the wrong objecthash for %T{ %[1]v }.\n"+
						"Actual:   %v\nExpected: %v\n", message, messageHashStr, tc.ExpectedHashString)
				}
			})
		}

		// If the test case has an equivalent JSON String, check it.
		if tc.EquivalentJSONString != "" {
			t.Run("Compare to objecthash of the equialent JSON", func(t *testing.T) {
				commonJSONHash, err := objecthash.CommonJSONHash(tc.EquivalentJSONString)
				if err != nil {
					t.Errorf("Attempting to hash %+v returned an error: %v", tc.EquivalentJSONString, err)
				}
				commonJSONHashStr := fmt.Sprintf("%x", commonJSONHash)

				if messageHashStr != commonJSONHashStr {
					t.Errorf("The objecthash for %T{ %[1]v } was expected to be the same as that of %+v.\n"+
						"Actual:   %v\nExpected: %v\n", message, tc.EquivalentJSONString, messageHashStr, commonJSONHashStr)
				}
			})
		}

		// If the test case has an equivalent object, check it.
		if tc.EquivalentObject != nil {
			t.Run("Compare to objecthash of the equialent Go object", func(t *testing.T) {
				EquivalentObjectHash, err := objecthash.ObjectHash(tc.EquivalentObject)
				if err != nil {
					t.Errorf("Attempting to hash %+v returned an error: %v", tc.EquivalentObject, err)
				}
				EquivalentObjectHashStr := fmt.Sprintf("%x", EquivalentObjectHash)

				if messageHashStr != EquivalentObjectHashStr {
					t.Errorf("The objecthash for %T{ %[1]v } was expected to be the same as that of %+v.\n"+
						"Actual:   %v\nExpected: %v\n", message, tc.EquivalentObject, messageHashStr, EquivalentObjectHashStr)
				}
			})
		}
	}
}
