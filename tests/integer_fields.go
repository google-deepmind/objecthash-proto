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

// TestIntegerFields performs tests on how integers are handled.
func TestIntegerFields(t *testing.T, hashers oi.ProtoHashers) {
	hasher := hashers.StringPreferringHasher

	testCases := []ti.TestCase{
		///////////////////////////////
		//  Equivalence of Integers. //
		///////////////////////////////
		{
			Protos: []proto.Message{
				&pb2_latest.Fixed32Message{Values: []uint32{0, 1, 2}},
				&pb2_latest.Fixed64Message{Values: []uint64{0, 1, 2}},
				&pb2_latest.Int32Message{Values: []int32{0, 1, 2}},
				&pb2_latest.Int64Message{Values: []int64{0, 1, 2}},
				&pb2_latest.Sfixed32Message{Values: []int32{0, 1, 2}},
				&pb2_latest.Sfixed64Message{Values: []int64{0, 1, 2}},
				&pb2_latest.Sint32Message{Values: []int32{0, 1, 2}},
				&pb2_latest.Sint64Message{Values: []int64{0, 1, 2}},
				&pb2_latest.Uint32Message{Values: []uint32{0, 1, 2}},
				&pb2_latest.Uint64Message{Values: []uint64{0, 1, 2}},

				&pb3_latest.Fixed32Message{Values: []uint32{0, 1, 2}},
				&pb3_latest.Fixed64Message{Values: []uint64{0, 1, 2}},
				&pb3_latest.Int32Message{Values: []int32{0, 1, 2}},
				&pb3_latest.Int64Message{Values: []int64{0, 1, 2}},
				&pb3_latest.Sfixed32Message{Values: []int32{0, 1, 2}},
				&pb3_latest.Sfixed64Message{Values: []int64{0, 1, 2}},
				&pb3_latest.Sint32Message{Values: []int32{0, 1, 2}},
				&pb3_latest.Sint64Message{Values: []int64{0, 1, 2}},
				&pb3_latest.Uint32Message{Values: []uint32{0, 1, 2}},
				&pb3_latest.Uint64Message{Values: []uint64{0, 1, 2}},
			},
			EquivalentObject: map[string][]int32{"values": {0, 1, 2}},
			// No equivalent JSON: JSON does not have an "integer" type. All numbers are floats.
			ExpectedHashString: "42794fb0e73c2b5f427aa76486555d07589359054848396ddf173e9e0b4ab931",
		},

		{
			Protos: []proto.Message{
				&pb2_latest.Int32Message{Values: []int32{-2, -1, 0, 1, 2}},
				&pb2_latest.Int64Message{Values: []int64{-2, -1, 0, 1, 2}},
				&pb2_latest.Sfixed32Message{Values: []int32{-2, -1, 0, 1, 2}},
				&pb2_latest.Sfixed64Message{Values: []int64{-2, -1, 0, 1, 2}},
				&pb2_latest.Sint32Message{Values: []int32{-2, -1, 0, 1, 2}},
				&pb2_latest.Sint64Message{Values: []int64{-2, -1, 0, 1, 2}},

				&pb3_latest.Int32Message{Values: []int32{-2, -1, 0, 1, 2}},
				&pb3_latest.Int64Message{Values: []int64{-2, -1, 0, 1, 2}},
				&pb3_latest.Sfixed32Message{Values: []int32{-2, -1, 0, 1, 2}},
				&pb3_latest.Sfixed64Message{Values: []int64{-2, -1, 0, 1, 2}},
				&pb3_latest.Sint32Message{Values: []int32{-2, -1, 0, 1, 2}},
				&pb3_latest.Sint64Message{Values: []int64{-2, -1, 0, 1, 2}},
			},
			EquivalentObject: map[string][]int32{"values": {-2, -1, 0, 1, 2}},
			// No equivalent JSON: JSON does not have an "integer" type. All numbers are floats.
			ExpectedHashString: "6cb613a53b6086b88dbda40b30e902adb41288b0b1f7a627905beaa764ee49cb",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
