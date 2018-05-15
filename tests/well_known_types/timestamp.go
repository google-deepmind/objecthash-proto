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

package wellknowntypes

import (
	"testing"

	"github.com/golang/protobuf/proto"
	timestamp_pb "github.com/golang/protobuf/ptypes/timestamp"

	oi "github.com/deepmind/objecthash-proto/internal"
	pb2_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto2"
	pb3_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto3"
	ti "github.com/deepmind/objecthash-proto/tests/internal"
)

// TestTimestamps confirms that google.protobuf.timestamp protos are hashed properly.
func TestTimestamps(t *testing.T, hashers oi.ProtoHashers) {
	hasher := hashers.FieldNamesAsKeysHasher

	testCases := []ti.TestCase{
		//////////////////////////////
		//  Empty/Zero Timestamps. //
		//////////////////////////////

		// The semantics of the Timestamp object imply that the distinction between
		// unset and zero happen at the message level, rather than the field level.
		//
		// As a result, an unset timestamp is one where the proto itself is nil,
		// while an explicitly set timestamp with unset fields is considered to be
		// explicitly set to 0.
		//
		// This is unlike normal proto3 messages, where unset/zero fields must be
		// considered to be unset, because they're indistinguishable in the general
		// case.
		{
			Protos: []proto.Message{
				&timestamp_pb.Timestamp{},
				&timestamp_pb.Timestamp{Seconds: 0, Nanos: 0},
			},
			// JSON treats all numbers as floats, so it is not possible to have an equivalent JSON string.
			EquivalentObject:   []int64{0, 0},
			ExpectedHashString: "3a82b649344529f03f52c1833f5aecc488a53b31461a1f54c305d149b12b8f53",
		},

		/////////////////////////
		//  Normal Timestamps. //
		/////////////////////////
		{
			Protos: []proto.Message{
				&timestamp_pb.Timestamp{Seconds: 1525450021, Nanos: 123456789},
			},
			// JSON treats all numbers as floats, so it is not possible to have an equivalent JSON string.
			EquivalentObject:   []int64{1525450021, 123456789},
			ExpectedHashString: "1fd36770664df599ad44e4e4f06b1fad6ef7a4b3f316d79ca11bea668032a199",
		},

		//////////////////////////////////////
		//  Timestamps within other protos. //
		//////////////////////////////////////

		// As mentioned above, a timestamp with unset fields is considered to be a
		// timestamp explicitly set to zero.
		{
			Protos: []proto.Message{
				&pb2_latest.KnownTypes{TimestampField: &timestamp_pb.Timestamp{}},
				&pb2_latest.KnownTypes{TimestampField: &timestamp_pb.Timestamp{Seconds: 0, Nanos: 0}},

				&pb3_latest.KnownTypes{TimestampField: &timestamp_pb.Timestamp{}},
				&pb3_latest.KnownTypes{TimestampField: &timestamp_pb.Timestamp{Seconds: 0, Nanos: 0}},
			},
			// JSON treats all numbers as floats, so it is not possible to have an equivalent JSON string.
			EquivalentObject:   map[string][]int64{"timestamp_field": {0, 0}},
			ExpectedHashString: "8457fe431752dbc5c47301c2546fcf6f0ad8c5317092b443e187d18e312e497e",
		},

		{
			Protos: []proto.Message{
				&pb2_latest.KnownTypes{TimestampField: &timestamp_pb.Timestamp{Seconds: 1525450021, Nanos: 123456789}},
				&pb3_latest.KnownTypes{TimestampField: &timestamp_pb.Timestamp{Seconds: 1525450021, Nanos: 123456789}},
			},
			// JSON treats all numbers as floats, so it is not possible to have an equivalent JSON string.
			EquivalentObject:   map[string][]int64{"timestamp_field": {1525450021, 123456789}},
			ExpectedHashString: "cf99942e3f8d1212f4ce263e206d64e29525b97b91368e71f9595bce83ac6a3e",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
