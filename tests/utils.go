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
	pb2_latest "../test_protos/generated/latest/proto2"

	"github.com/golang/protobuf/proto"
)

type ProtoHasher interface {
	HashProto(pb proto.Message) ([]byte, error)
}

type ProtoHashers struct {
	DefaultHasher                 ProtoHasher
	FieldNamesAsKeysHasher        ProtoHasher
	EnumsAsStringsHasher          ProtoHasher
	StringPreferringHasher        ProtoHasher
	CustomMessageIdentifierHasher ProtoHasher
}

func forgetAllFields(originalMessage proto.Message) proto.Message {
	emptyMessage := &pb2_latest.Empty{}

	binaryMessage, err := proto.Marshal(originalMessage)
	if err != nil {
		panic(err)
	}

	err = proto.Unmarshal(binaryMessage, emptyMessage)
	if err != nil {
		panic(err)
	}
	return emptyMessage
}
