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

	pb2_latest "../test_protos/generated/latest/proto2"
	pb3_latest "../test_protos/generated/latest/proto2"

	"github.com/golang/protobuf/proto"
)

func TestBadness(t *testing.T, hashers ProtoHashers) {
	hasher := hashers.DefaultHasher

	badProtos := []proto.Message{
		// Nil messages in repeated fields are bad.
		&pb2_latest.Repetitive{SimpleField: []*pb2_latest.Simple{nil}},
		&pb3_latest.Repetitive{SimpleField: []*pb3_latest.Simple{nil}},

		// Nil messages in maps are bad.
		&pb2_latest.IntMaps{IntToSimple: map[int64]*pb2_latest.Simple{3: nil}},
		&pb3_latest.IntMaps{IntToSimple: map[int64]*pb3_latest.Simple{3: nil}},

		// Custom default values are bad.
		&pb2_latest.BadWithDefaults{},

		&pb2_latest.BadWithDefaults{Text: proto.String("Schlecht!")},

		// Required fields are bad.
		&pb2_latest.BadWithRequirements{},

		&pb2_latest.BadWithRequirements{Text: proto.String("Schlecht!")},

		// Extensions are bad.
		&pb2_latest.BadWithExtensions{},

		// Create proto messages with unknown fields. That's bad.
		forgetAllFields(&pb2_latest.PersonV1{Name: proto.String("Unbekannt")}),
	}
	for _, message := range badProtos {
		_, err := hasher.HashProto(message)
		if err == nil {
			t.Errorf("Attempting to hash %T{ %+v} should have returned an error.", message, message)
		}
	}
}
