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
	"testing"

	"github.com/golang/protobuf/proto"

	oi "github.com/deepmind/objecthash-proto/internal"
	pb2_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto2"
	pb3_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto3"
	ti "github.com/deepmind/objecthash-proto/tests/internal"
)

// TestOtherTypes performs tests on types that do not have their own test file.
func TestOtherTypes(t *testing.T, hashers oi.ProtoHashers) {
	hasher := hashers.StringPreferringHasher

	testCases := []ti.TestCase{
		///////////
		//  Nil. //
		///////////
		{
			Protos: []proto.Message{
				nil,
			},
			EquivalentJSONString: "null",
			EquivalentObject:     nil,
			ExpectedHashString:   "1b16b1df538ba12dc3f97edbb85caa7050d46c148134290feba80f8236c83db9",
		},

		/////////////////////
		// Boolean fields. //
		/////////////////////
		{
			Protos: []proto.Message{
				&pb2_latest.Simple{BoolField: proto.Bool(true)},
				&pb3_latest.Simple{BoolField: true},
			},
			EquivalentJSONString: "{\"bool_field\": true}",
			EquivalentObject:     map[string]bool{"bool_field": true},
			ExpectedHashString:   "7b2ac6048e6c8797205505ea486539a5589583be43154da88785a5121e2d6899",
		},

		{
			Protos: []proto.Message{
				&pb2_latest.Simple{BoolField: proto.Bool(false)},
				// proto3 scalar fields set to their default value are considered empty.
			},
			EquivalentJSONString: "{\"bool_field\": false}",
			EquivalentObject:     map[string]bool{"bool_field": false},
			ExpectedHashString:   "1ab5ecdbe4176473024f7efd080593b740d22d076d06ea6edd8762992b484a12",
		},

		///////////////////
		// Bytes fields. //
		///////////////////
		{
			Protos: []proto.Message{
				&pb2_latest.Simple{BytesField: []byte{0, 0, 0}},
				&pb3_latest.Simple{BytesField: []byte{0, 0, 0}},
			},
			// No equivalent JSON: JSON does not have a "bytes" type.
			EquivalentObject:   map[string][]byte{"bytes_field": []byte("\000\000\000")},
			ExpectedHashString: "fdd59e1f3120117943124cb9c39da79ac47ea631343ff9154dffb0e64550789c",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
