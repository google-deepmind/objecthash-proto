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

package tests

import (
	"testing"

	pb2_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto2"
	pb3_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto3"

	"github.com/golang/protobuf/proto"
)

// TestStringFields performs tests on how strings are handled.
func TestStringFields(t *testing.T, hashers ProtoHashers) {
	hasher := hashers.StringPreferringHasher

	testCases := []testCase{
		{
			protos: []proto.Message{
				&pb2_latest.Simple{StringField: proto.String("你好")},
				&pb3_latest.Simple{StringField: "你好"},
			},
			equivalentObject:     map[string]string{"string_field": "你好"},
			equivalentJSONString: "{\"string_field\":\"你好\"}",
			expectedHashString:   "de0086ad683b5f8affffbbcbe57d09e5377aa47cb32f6f0b1bdecd2e54b9137d",
		},

		{
			protos: []proto.Message{
				&pb2_latest.Simple{StringField: proto.String("\u03d3")},
				&pb3_latest.Simple{StringField: "\u03d3"},
			},
			equivalentObject:     map[string]string{"string_field": "\u03d3"},
			equivalentJSONString: "{\"string_field\":\"\u03d3\"}",
			expectedHashString:   "12441188aebffcc3a1e625d825391678d8417c77e645fc992d1ab5b549c659a7",
		},

		// Note that this is the same character as above, but hashes differently
		// unless unicode normalisation is applied.
		{
			protos: []proto.Message{
				&pb2_latest.Simple{StringField: proto.String("\u03d2\u0301")},
				&pb3_latest.Simple{StringField: "\u03d2\u0301"},
			},
			equivalentObject:     map[string]string{"string_field": "\u03d2\u0301"},
			equivalentJSONString: "{\"string_field\":\"\u03d2\u0301\"}",
			expectedHashString:   "1f33a91552e7a527fdf2de0d25f815590f1a3e2dc8340507d20d4ee42462d0a2",
		},

		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{StringField: []string{""}},
				&pb3_latest.Repetitive{StringField: []string{""}},
			},
			equivalentObject:     map[string][]string{"string_field": {""}},
			equivalentJSONString: "{\"string_field\":[\"\"]}",
			expectedHashString:   "63e64f0ed286e0d8f30735e6646ea9ef48174c23ba09a05288b4233c6e6a9419",
		},

		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{StringField: []string{"", "Test", "你好", "\u03d3"}},
				&pb3_latest.Repetitive{StringField: []string{"", "Test", "你好", "\u03d3"}},
			},
			equivalentObject:     map[string][]string{"string_field": {"", "Test", "你好", "\u03d3"}},
			equivalentJSONString: "{\"string_field\":[\"\",\"Test\",\"你好\",\"\u03d3\"]}",
			expectedHashString:   "f76ae15a2685a5ec0e45f9ad7d75e492e6a17d31811480fbaf00af451fb4e98e",
		},
	}

	for _, tc := range testCases {
		if err := tc.check(hasher); err != nil {
			t.Errorf("%s", err)
		}
	}
}
