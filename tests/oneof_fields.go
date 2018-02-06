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

	"github.com/golang/protobuf/proto"

	pb2_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto2"
	pb3_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto3"
)

// TestOneOfFields checks that oneof fields are handled properly.
func TestOneOfFields(t *testing.T, hashers ProtoHashers) {
	hasher := hashers.DefaultHasher

	testCases := []testCase{
		//////////////////////////
		//  Empty oneof fields. //
		//////////////////////////
		{
			protos: []proto.Message{
				&pb2_latest.Singleton{},
				&pb3_latest.Singleton{},

				&pb2_latest.Empty{},
				&pb3_latest.Empty{},
			},
			equivalentJSONString: "{}",
			equivalentObject:     map[int64]string{},
			expectedHashString:   "18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4",
		},

		/////////////////////////////////////////////
		//  One of the options selected but empty. //
		/////////////////////////////////////////////
		{
			protos: []proto.Message{
				// Only proto2 has empty values.
				&pb2_latest.Simple{BoolField: proto.Bool(false)},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheBool{}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheBool{}},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheBool{TheBool: false}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheBool{TheBool: false}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			equivalentObject:   map[int64]bool{1: false},
			expectedHashString: "8a956cfa8e9b45b738cb8dc8a3dc7126dab3cbd2c07c80fa1ec312a1a31ed709",
		},

		{
			protos: []proto.Message{
				// Only proto2 has empty values.
				&pb2_latest.Simple{StringField: proto.String("")},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheString{}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheString{}},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheString{TheString: ""}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheString{TheString: ""}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			equivalentObject:   map[int64]string{25: ""},
			expectedHashString: "79cff9d2d0ee6c6071c82b58d1a2fcf056b58c4501606862489e5731644c755a",
		},

		{
			protos: []proto.Message{
				// Only proto2 has empty values.
				&pb2_latest.Simple{Int32Field: proto.Int32(0)},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheInt32{}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheInt32{}},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheInt32{TheInt32: 0}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheInt32{TheInt32: 0}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			equivalentObject:   map[int64]int32{13: 0},
			expectedHashString: "bafd42680c987c47a76f72e08ed975877162efdb550d2c564c758dc7d988468f",
		},

		////////////////////////////////////////////////
		//  One of the options selected with content. //
		////////////////////////////////////////////////
		{
			protos: []proto.Message{
				&pb2_latest.Simple{StringField: proto.String("TEST!")},
				&pb3_latest.Simple{StringField: "TEST!"},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheString{TheString: "TEST!"}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheString{TheString: "TEST!"}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			equivalentObject:   map[int64]string{25: "TEST!"},
			expectedHashString: "336cdbca99fd46157bc47bcc456f0ac7f1ef3be7a79acf3535f671434b53944f",
		},

		{
			protos: []proto.Message{
				&pb2_latest.Simple{Int32Field: proto.Int32(99)},
				&pb3_latest.Simple{Int32Field: 99},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheInt32{TheInt32: 99}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheInt32{TheInt32: 99}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			equivalentObject:   map[int64]int32{13: 99},
			expectedHashString: "65517521bc278528d25caf1643da0f094fd88dad50205c9743e3c984a7c53b7d",
		},

		///////////////////////////
		//  Nested oneof fields. //
		///////////////////////////
		{
			protos: []proto.Message{
				&pb2_latest.Simple{SingletonField: &pb2_latest.Singleton{}},
				&pb3_latest.Simple{SingletonField: &pb3_latest.Singleton{}},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheSingleton{TheSingleton: &pb2_latest.Singleton{}}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheSingleton{TheSingleton: &pb3_latest.Singleton{}}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			equivalentObject: map[int64]map[int64]int64{35: {}},
			// equivalentObject:   map[int64]map[int64]map[int64]int64{35: {35: {}}},
			expectedHashString: "4967c72525c764229f9fbf1294764c9aedc0d4f9f4c52e04a19c7f35ca65f517",
		},

		{
			protos: []proto.Message{
				&pb2_latest.Simple{SingletonField: &pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheSingleton{TheSingleton: &pb2_latest.Singleton{}}}},
				&pb3_latest.Simple{SingletonField: &pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheSingleton{TheSingleton: &pb3_latest.Singleton{}}}},

				&pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheSingleton{TheSingleton: &pb2_latest.Singleton{Singleton: &pb2_latest.Singleton_TheSingleton{TheSingleton: &pb2_latest.Singleton{}}}}},
				&pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheSingleton{TheSingleton: &pb3_latest.Singleton{Singleton: &pb3_latest.Singleton_TheSingleton{TheSingleton: &pb3_latest.Singleton{}}}}},
			},
			// No equivalent JSON because JSON maps have to have strings as keys.
			equivalentObject:   map[int64]map[int64]map[int64]int64{35: {35: {}}},
			expectedHashString: "8ea95bbda0f42073a61f46f9f375f48d5a7cb034fce56b44f958470fda5236d0",
		},
	}

	for _, tc := range testCases {
		if err := tc.check(hasher); err != nil {
			t.Errorf("%s", err)
		}

		if err := checkAsSingletonOnTheWire(tc, hasher); err != nil {
			t.Errorf("%s", err)
		}
	}
}

// Checks the provided test case after all its proto messages have been cycled
// to their wire format and unmarshalled back as a Singleton message.
func checkAsSingletonOnTheWire(tc testCase, hasher ProtoHasher) error {
	testCaseAfterAWireTransfer := testCase{
		protos:               tc.protos,
		equivalentJSONString: tc.equivalentJSONString,
		equivalentObject:     tc.equivalentObject,
		expectedHashString:   tc.expectedHashString,
	}

	for i, pb := range tc.protos {
		testCaseAfterAWireTransfer.protos[i] = unmarshalAsSingletonOnTheWire(pb)
	}

	return testCaseAfterAWireTransfer.check(hasher)
}

// Marshals a proto message to its wire format and returns its
// unmarshalled Singleton message.
func unmarshalAsSingletonOnTheWire(original proto.Message) proto.Message {
	binary, err := proto.Marshal(original)
	if err != nil {
		panic(err)
	}

	singleton := &pb3_latest.Singleton{}
	err = proto.Unmarshal(binary, singleton)
	if err != nil {
		panic(err)
	}

	return singleton
}
