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

package protohash

import (
	"reflect"
	"testing"
)

func TestIsAProto2BytesFieldWithBadArguments(t *testing.T) {
	var emptyV reflect.Value
	var emptySF reflect.StructField
	if isAProto2BytesField(emptyV, emptySF) {
		t.Error("isAProto2BytesField incorrectly returned true for zero values.")
	}

	v := reflect.ValueOf(struct {
		s string
		i int64
	}{"Hello", 42})
	tp := v.Type()

	if isAProto2BytesField(v.Field(0), tp.Field(0)) {
		t.Error("isAProto2BytesField incorrectly returned true for a string.")
	}

	if isAProto2BytesField(v.Field(1), tp.Field(1)) {
		t.Error("isAProto2BytesField incorrectly returned true for an int.")
	}
}
