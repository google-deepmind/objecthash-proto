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

	pb2_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto2"
	pb3_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto3"
)

func TestMaps(t *testing.T, hashers ProtoHashers) {
	hasher := hashers.StringPreferringHasher

	testCases := []testCase{
		////////////////////
		//  Boolean maps. //
		////////////////////
		{
			protos: []proto.Message{
				&pb2_latest.BoolMaps{BoolToString: map[bool]string{true: "NOT FALSE", false: "NOT TRUE"}},
				&pb3_latest.BoolMaps{BoolToString: map[bool]string{true: "NOT FALSE", false: "NOT TRUE"}},
			},
			// No equivalent JSON object because JSON map keys must be strings.
			equivalentObject:   map[string]map[bool]string{"bool_to_string": map[bool]string{true: "NOT FALSE", false: "NOT TRUE"}},
			expectedHashString: "d89d053bf7b37b4784832c72445661db99538fe1d490988575409a9040084f18",
		},

		////////////////////
		//  Integer maps. //
		////////////////////
		{
			protos: []proto.Message{
				&pb2_latest.IntMaps{IntToString: map[int64]string{0: "ZERO"}},
				&pb3_latest.IntMaps{IntToString: map[int64]string{0: "ZERO"}},
			},
			// No equivalent JSON object because JSON map keys must be strings.
			equivalentObject:   map[string]map[int64]string{"int_to_string": map[int64]string{0: "ZERO"}},
			expectedHashString: "53892192fb69cbd93ceb0552ca571b8505887f25d6f12822025341f16983a6af",
		},

		///////////////////
		//  String maps. //
		///////////////////
		{
			protos: []proto.Message{
				&pb2_latest.StringMaps{StringToString: map[string]string{"foo": "bar"}},
				&pb3_latest.StringMaps{StringToString: map[string]string{"foo": "bar"}},
			},
			equivalentJsonString: "{\"string_to_string\": {\"foo\": \"bar\"}}",
			equivalentObject:     map[string]map[string]string{"string_to_string": map[string]string{"foo": "bar"}},
			expectedHashString:   "cadfe560995647c63c20234a6409d2b1b8cf8dcf7d8e420ca33f23ff9ca9abfa",
		},
	}

	for _, tc := range testCases {
		if err := tc.check(hasher); err != nil {
			t.Errorf("%s", err)
		}
	}
}
