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

// An `Option` modifies how ObjectHashes for protobufs is calculated.
type Option interface {
	set(*objectHasher)
	fmt.Stringer
}

// EnumsAsStrings returns an Option to specify that enum values should be hashed
// as strings instead of being hashed as integers.
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
func FieldNamesAsKeys() Option { return fieldNamesAsKeys{} }

type fieldNamesAsKeys struct{}

func (x fieldNamesAsKeys) set(oh *objectHasher) {
	oh.fieldNamesAsKeys = true
}

func (x fieldNamesAsKeys) String() string {
	return "FieldNamesAsKeys"
}

// MessageIdentifier returns an Option to specify that proto messages
// should use `i` as their type identifier. This is useful to make sure that
// messages have a different hash from maps.
func MessageIdentifier(i string) Option { return messageIdentifier(i) }

type messageIdentifier string

func (x messageIdentifier) set(oh *objectHasher) {
	oh.messageIdentifier = string(x)
}

func (x messageIdentifier) String() string {
	return fmt.Sprintf("MessageIdentifier(%v)", string(x))
}
