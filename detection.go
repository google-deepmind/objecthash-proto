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
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
)

// wktProto is the interface satisfied by the proto library's well-known types.
//
// Notice that the method XXX_WellKnownType requires a pointer receiver, so you
// should take the Addr() of a reflect.value before checking if it satisfies
// this interface.
type wktProto interface {
	proto.Message
	XXX_WellKnownType() string
}

// CheckWellKnownType checks if a value is one of the proto library's well-known
// types. The ok return value reports whether the value is a well-known type or
// not, while the name return value reports the name provided by the proto
// library for this type.
//
// This is done by checking if the proto message has a XXX_WellKnownType method
// defined on it.
func CheckWellKnownType(sv reflect.Value) (name string, ok bool) {
	// The method XXX_WellKnownType requires a pointer receiver.
	wellKnownValue, ok := sv.Addr().Interface().(wktProto)
	if ok {
		return wellKnownValue.XXX_WellKnownType(), true
	}

	return "", false
}

// isExtendable checks if the proto message is extendable.
//
// This is done by checking if it has the ExtensionRangeArray method.
func isExtendable(sv reflect.Value) bool {
	type extendableProto interface {
		proto.Message
		ExtensionRangeArray() []proto.ExtensionRange
	}

	_, ok := sv.Addr().Interface().(extendableProto)
	return ok
}

// isRawMessageField checks if the proto field is a RawMessage.
//
// This is done by checking if it has the Bytes method.
func isRawMessageField(v reflect.Value) bool {
	type rawMessageField interface {
		Bytes() []byte
	}

	_, ok := v.Interface().(rawMessageField)
	return ok
}

// isAOneOfField checks if the proto field is a oneof wrapper field.
//
// This is done by checking if it is an interface whose tag has a
// "protobuf_oneof" entry.
func isAOneOfField(v reflect.Value, sf reflect.StructField) bool {
	return v.Kind() == reflect.Interface && sf.Tag.Get("protobuf_oneof") != ""
}

// isAProto2BytesField checks if the field is a proto2 bytes field.
//
// This is done by checking the field's tag. Byte fields do not have a 'rep'
// tag (for repeated fields).  Additionally, proto3 byte fields have a 'proto3'
// tag, which proto2 byte fields do not have.
func isAProto2BytesField(v reflect.Value, sf reflect.StructField) bool {
	k := v.Kind()
	if k != reflect.Map && k != reflect.Slice {
		return false
	}

	tag := sf.Tag.Get("protobuf")

	// This can potentially break "def" values which come last and can contain
	// commas. This is fine because we don't care about "def" or its value here.
	//
	// Example: "bytes,10,opt,name=F_String,json=FString,def=hello, \"world!\"\n"
	fields := strings.Split(tag, ",")

	// Normal proto tags (incl. bytes) will have at least 2 fields.
	if len(fields) < 2 {
		// We do not return an error because this case can happen for certain kinds
		// of fields, such as 'XXX_unrecognized' for example.
		return false
	}

	for i := 2; i < len(fields); i++ {
		f := fields[i]
		if f == "rep" || f == "proto3" {
			return false
		}
	}

	return true
}

// isContentIndependentField checks if a proto field is not related to its
// contents.
//
// Content independent fields include:
// - "XXX_NoUnkeyedLiteral": This is an empty struct whose presence is used by
//   the proto library to force users to use keyed literals to ensure forward
//   compatible usage of messages.
// - "XXX_sizecache": This is used in serialization. Its value is irrelevant to
//   the actual content of the proto and should not be used in comparisons or
//   for generating hashes.
func isContentIndependentField(v reflect.Value, sf reflect.StructField) bool {
	name := sf.Name
	return name == "XXX_NoUnkeyedLiteral" || name == "XXX_sizecache"
}

// isUnset checks if the proto field has not been set.
//
// This also includes empty proto3 scalar values.
func isUnset(v reflect.Value, sf reflect.StructField) (bool, error) {
	// Default values are considered empty. Otherwise, adding those kinds of
	// fields to a proto's definition would break all older hashes.
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool(), nil
	case reflect.Int32, reflect.Int64:
		return v.Int() == 0, nil
	case reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0, nil
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0, nil
	case reflect.String:
		return v.String() == "", nil
	case reflect.Map, reflect.Slice:
		// Proto2 bytes fields are considered scalar fields, which means that there
		// is a distinction between unset fields and fields set to zero values.
		//
		// Therefore, if we encounter a proto2 bytes field, we should only check if
		// it's nil or not, rather than checking its value.
		if isAProto2BytesField(v, sf) {
			return v.IsNil(), nil
		}

		// If this is not a proto2 bytes field, then empty values are always unset.
		//
		// This applies for: repeated fields, proto3 bytes fields, and special
		// fields (ex. XXX_unrecognized).
		return v.Len() == 0, nil
	case reflect.Interface:
		// This only happens when we have a oneof field.
		//
		// Notice that a non-nil oneof interface implies the presence of the field
		// it represents. Therefore, even if the value of the selected field is zero
		// it is still different from (and not equal to) the empty value.
		//
		// For example, with a proto defined as:
		//
		// message M {
		//   oneof val {
		//     string a = 1;
		//     int32 b = 2;
		//     M2 c = 3;
		//   }
		// }
		//
		// The following are all different values:
		// - <>
		// - < a:"" >
		// - < b:0  >
		// - < c:<> >
		//
		// - The first does not set any oneof value (ie. it is empty).
		// - The second sets the oneof field to "a" and is non-empty even with an
		//   empty string.
		// - The third sets the oneof field to "b" and is non-empty even with its
		//   zero value.
		// - The fourth sets the oneof field to "c" and is non-empty even with an
		//   empty message. In fact, if "c" was set but was nil, the proto will be
		//   invalid. The check for this error is made when hashing the oneof value.
		return v.IsNil(), nil
	case reflect.Ptr:
		// If a pointer is not a null pointer, this means that the value it points
		// to is distinguishable from it being missing. Usually, that value would
		// be another proto message or a proto2 scalar value.
		return v.IsNil(), nil
	case reflect.Struct:
		// This should never happen because protobuf generated code never uses structs
		// as fields, and uses pointers to structs instead.
		// This means that emptiness checks for nested messages would happen in the
		// reflect.Ptr case rather than here.
		return false, fmt.Errorf("got an unexpected struct of type '%+v' for field %+v", v.Type(), sf)
	default:
		return false, fmt.Errorf("got an unexpected type '%+v' for field %+v", v.Type(), sf)
	}
}

// failIfUnsupported returns an error if the provided field cannot be hashed reliably.
//
// Note that unsupported fields are safe to ignore if they've not been set, so
// an isUnset check should be used before this check.
func failIfUnsupported(v reflect.Value, sf reflect.StructField) error {
	// Check "XXX_" fields.
	if name := sf.Name; strings.HasPrefix(name, "XXX_") {
		switch name {
		case "XXX_unrecognized":
			// A non-empty XXX_unrecognized field means that the proto message
			// contains some unrecognized fields.
			return errors.New("unrecognized fields cannot be hashed reliably")
		case "XXX_extensions", "XXX_InternalExtensions":
			return errors.New("extensions cannot be hashed reliably")
		default:
			return fmt.Errorf("found an unknown XXX_ field: '%s'", name)
		}
	}

	if isRawMessageField(v) {
		return errors.New("raw message fields not supported")
	}

	return nil
}
