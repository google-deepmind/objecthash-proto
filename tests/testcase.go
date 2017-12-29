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

package tests

import (
	"fmt"

	"github.com/benlaurie/objecthash/go/objecthash"
	"github.com/golang/protobuf/proto"
)

type testCase struct {
	protos               []proto.Message
	equivalentObject     interface{}
	equivalentJsonString string
	expectedHashString   string
}

func (tc testCase) check(hasher ProtoHasher) error {
	for _, message := range tc.protos {
		messageHash, err := hasher.HashProto(message)
		if err != nil {
			return fmt.Errorf("Attempting to hash %T{ %[1]v } returned an error: %v", message, err)
		}
		messageHashStr := fmt.Sprintf("%x", messageHash)

		// If the test case has an expected hash string, check it.
		if tc.expectedHashString != "" {
			if messageHashStr != tc.expectedHashString {
				return fmt.Errorf("Got the wrong objecthash for %T{ %[1]v }.\n"+
					"Actual:   %v\nExpected: %v\n", message, messageHashStr, tc.expectedHashString)
			}
		}

		// If the test case has an equivalent JSON String, check it.
		if tc.equivalentJsonString != "" {
			commonJsonHash, err := objecthash.CommonJSONHash(tc.equivalentJsonString)
			if err != nil {
				return fmt.Errorf("Attempting to hash %+v returned an error: %v", tc.equivalentJsonString, err)
			}
			commonJsonHashStr := fmt.Sprintf("%x", commonJsonHash)

			if messageHashStr != commonJsonHashStr {
				return fmt.Errorf("The objecthash for %T{ %[1]v } was expected to be the same as that of %+v.\n"+
					"Actual:   %v\nExpected: %v\n", message, tc.equivalentJsonString, messageHashStr, commonJsonHashStr)
			}
		}

		// If the test case has an equivalent object, check it.
		if tc.equivalentObject != nil {
			equivalentObjectHash, err := objecthash.ObjectHash(tc.equivalentObject)
			if err != nil {
				return fmt.Errorf("Attempting to hash %+v returned an error: %v", tc.equivalentObject, err)
			}
			equivalentObjectHashStr := fmt.Sprintf("%x", equivalentObjectHash)

			if messageHashStr != equivalentObjectHashStr {
				return fmt.Errorf("The objecthash for %T{ %[1]v } was expected to be the same as that of %+v.\n"+
					"Actual:   %v\nExpected: %v\n", message, tc.equivalentObject, messageHashStr, equivalentObjectHashStr)
			}
		}
	}

	return nil
}
