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

package protohash

import (
	"github.com/golang/protobuf/proto"
)

// ProtoHasher is an interface for hashers that are capable of returning an
// ObjectHash for protobufs.
type ProtoHasher interface {
	HashProto(pb proto.Message) ([]byte, error)
}

// NewHasher creates a new ProtoHasher with the options specified in the argument.
func NewHasher(opts ...Option) ProtoHasher {
	hasher := objectHasher{messageIdentifier: mapIdentifier}
	for _, opt := range opts {
		opt.set(&hasher)
	}
	return &hasher
}
