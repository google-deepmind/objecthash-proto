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

// TestEmptyFields checks that empty proto fields are handled properly.
func TestEmptyFields(t *testing.T, hashers oi.ProtoHashers) {
	hasher := hashers.DefaultHasher

	testCases := []ti.TestCase{
		{
			Protos: []proto.Message{
				&pb2_latest.Empty{},
				&pb3_latest.Empty{},

				// Empty repeated fields are ignored.
				&pb2_latest.Repetitive{StringField: []string{}},
				&pb3_latest.Repetitive{StringField: []string{}},

				// Empty map fields are ignored.
				&pb2_latest.StringMaps{StringToString: map[string]string{}},
				&pb3_latest.StringMaps{StringToString: map[string]string{}},

				// Proto3 scalar fields set to their default values are considered empty.
				&pb3_latest.Simple{BoolField: false},
				&pb3_latest.Simple{BytesField: []byte{}},
				&pb3_latest.Simple{DoubleField: 0},
				&pb3_latest.Simple{DoubleField: 0.0},
				&pb3_latest.Simple{Fixed32Field: 0},
				&pb3_latest.Simple{Fixed64Field: 0},
				&pb3_latest.Simple{FloatField: 0},
				&pb3_latest.Simple{FloatField: 0.0},
				&pb3_latest.Simple{Int32Field: 0},
				&pb3_latest.Simple{Int64Field: 0},
				&pb3_latest.Simple{Sfixed32Field: 0},
				&pb3_latest.Simple{Sfixed64Field: 0},
				&pb3_latest.Simple{Sint32Field: 0},
				&pb3_latest.Simple{Sint64Field: 0},
				&pb3_latest.Simple{StringField: ""},
				&pb3_latest.Simple{Uint32Field: 0},
				&pb3_latest.Simple{Uint64Field: 0},
			},
			EquivalentJSONString: "{}",
			EquivalentObject:     map[string]interface{}{},
			ExpectedHashString:   "18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4",
		},
	}

	for _, tc := range testCases {
		tc.Check(t, hasher)
	}
}
