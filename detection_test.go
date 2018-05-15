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
	t.Run("Zero values should return false", func(t *testing.T) {
		var emptyV reflect.Value
		var emptySF reflect.StructField
		if isAProto2BytesField(emptyV, emptySF) {
			t.Error("isAProto2BytesField incorrectly returned true for zero values.")
		}
	})

	t.Run("String values should return false", func(t *testing.T) {
		v := reflect.ValueOf(struct{ s string }{"Hello"})
		tp := v.Type()

		if isAProto2BytesField(v.Field(0), tp.Field(0)) {
			t.Error("isAProto2BytesField incorrectly returned true for a string.")
		}
	})

	t.Run("Integer values should return false", func(t *testing.T) {
		v := reflect.ValueOf(struct{ i int64 }{42})
		tp := v.Type()

		if isAProto2BytesField(v.Field(0), tp.Field(0)) {
			t.Error("isAProto2BytesField incorrectly returned true for an int.")
		}
	})
}

func TestIsUnsetWithBadArguments(t *testing.T) {
	t.Run("Should return errors for struct fields", func(t *testing.T) {
		v := reflect.ValueOf(struct {
			s struct{}
		}{})
		tp := v.Type()

		_, err := isUnset(v.Field(0), tp.Field(0))

		if err == nil {
			t.Error("isUnset should have returned an error when running on a struct field.")
		}

		// Make sure that the returned error is meaningful.
		expectedError := "got an unexpected struct of type 'struct {}' for field {Name:s PkgPath:github.com/deepmind/objecthash-proto Type:struct {} Tag: Offset:0 Index:[0] Anonymous:false}"
		if err.Error() != expectedError {
			t.Errorf("Expected the error returned by isUnset to be '%s'. Instead got '%s'.", expectedError, err)
		}
	})

	t.Run("Should return errors for fields with unsupported types", func(t *testing.T) {
		v := reflect.ValueOf(struct {
			c chan interface{}
		}{})
		tp := v.Type()

		_, err := isUnset(v.Field(0), tp.Field(0))

		if err == nil {
			t.Error("isUnset should have returned an error when running on a field with an unexpected type.")
		}

		// Make sure that the returned error is meaningful.
		expectedError := "got an unexpected type 'chan interface {}' for field {Name:c PkgPath:github.com/deepmind/objecthash-proto Type:chan interface {} Tag: Offset:0 Index:[0] Anonymous:false}"
		if err.Error() != expectedError {
			t.Errorf("Expected the error returned by isUnset to be '%s'. Instead got '%s'.", expectedError, err)
		}
	})
}
