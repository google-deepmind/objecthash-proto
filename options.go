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
	"fmt"
)

// Option modifies how ObjectHashes for protobufs is calculated.
type Option interface {
	set(*objectHasher)
	fmt.Stringer
}

// EnumsAsStrings returns an Option to specify that enum values should be hashed
// as strings instead of being hashed as integers.
//
// This can be useful for compatibility with non-protobuf formats that
// represent enums as strings, but will have backward-compatiblity consequences
// for the proto message itself.
func EnumsAsStrings() Option { return enumsAsStrings{} }

type enumsAsStrings struct{}

func (x enumsAsStrings) set(oh *objectHasher) {
	oh.enumsAsStrings = true
}

func (x enumsAsStrings) String() string {
	return "EnumsAsStrings"
}

// FieldNamesAsKeys returns an Option to specify that field names should be used
// as their keys instead of using their tag number.
//
// This can be useful for compatibility with non-protobuf formats that
// primarily use strings (rather than integers) as keys, but will have
// backward-compatiblity consequences for the proto message itself.
func FieldNamesAsKeys() Option { return fieldNamesAsKeys{} }

type fieldNamesAsKeys struct{}

func (x fieldNamesAsKeys) set(oh *objectHasher) {
	oh.fieldNamesAsKeys = true
}

func (x fieldNamesAsKeys) String() string {
	return "FieldNamesAsKeys"
}

// MessageIdentifier returns an Option to specify that proto messages should
// use the supplied argument as their type identifier. This will make messages
// have a different hash from maps with equivalent contents.
//
// This can be useful for stricter type checking when comparing hashes, but
// will have consequences for the compatiblity with non-protobuf formats, in
// particular those that do not provide a container type for messages/structs
// distinct from that for maps/dictionaries.
func MessageIdentifier(i string) Option { return messageIdentifier(i) }

type messageIdentifier string

func (x messageIdentifier) set(oh *objectHasher) {
	oh.messageIdentifier = string(x)
}

func (x messageIdentifier) String() string {
	return fmt.Sprintf("MessageIdentifier(%v)", string(x))
}
