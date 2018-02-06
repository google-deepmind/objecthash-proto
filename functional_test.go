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

import "testing"

import "github.com/deepmind/objecthash-proto/tests"

func TestFunctional(t *testing.T) {
	protoHashers := tests.ProtoHashers{
		DefaultHasher:                 NewHasher(),
		FieldNamesAsKeysHasher:        NewHasher(FieldNamesAsKeys()),
		EnumsAsStringsHasher:          NewHasher(EnumsAsStrings()),
		StringPreferringHasher:        NewHasher(FieldNamesAsKeys(), EnumsAsStrings()),
		CustomMessageIdentifierHasher: NewHasher(MessageIdentifier(`m`)),
	}

	t.Run("TestBadness", func(t *testing.T) { tests.TestBadness(t, protoHashers) })
	t.Run("TestEmptyFields", func(t *testing.T) { tests.TestEmptyFields(t, protoHashers) })
	t.Run("TestFloatFields", func(t *testing.T) { tests.TestFloatFields(t, protoHashers) })
	t.Run("TestIntegerFields", func(t *testing.T) { tests.TestIntegerFields(t, protoHashers) })
	t.Run("TestMaps", func(t *testing.T) { tests.TestMaps(t, protoHashers) })
	t.Run("TestOneOfFields", func(t *testing.T) { tests.TestOneOfFields(t, protoHashers) })
	t.Run("TestOtherTypes", func(t *testing.T) { tests.TestOtherTypes(t, protoHashers) })
	t.Run("TestRepeatedFields", func(t *testing.T) { tests.TestRepeatedFields(t, protoHashers) })
	t.Run("TestStringFields", func(t *testing.T) { tests.TestStringFields(t, protoHashers) })
}
