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
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/golang/protobuf/proto"
)

// objectHasher is a configurable object for hashing protocol buffer objects.
// It implements the ProtoHasher interface.
type objectHasher struct {
	// Whether to hash enum values as strings, as opposed to as integer values.
	enumsAsStrings bool

	// Whether to use the proto field name as its key, as opposed to using the
	// tag number as the key.
	fieldNamesAsKeys bool

	// Custom type identifier for hashing proto messages, as opposed to using
	// the map identifier.
	messageIdentifier string
}

// HashProto returns the object hash of a given protocol buffer message.
func (hasher *objectHasher) HashProto(pb proto.Message) (h []byte, err error) {
	// Ensure that we can recover if the proto library panics.
	// See: https://github.com/golang/protobuf/issues/478
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	// Check if the value is nil.
	if pb == nil {
		return hashNil()
	}

	// Expliclity set any custom default values. This is done in order to detect
	// proto2 custom default values and return errors for them as soon as they're
	// introduced to the schema.
	proto.SetDefaults(pb)

	// Make sure the proto itself is actually valid (ie. can be marshalled).
	// If this fails, it probably means there are unset required fields or invalid
	// values.
	if _, err = proto.Marshal(pb); err != nil {
		return nil, err
	}

	val := reflect.ValueOf(pb)

	// Dereference the proto pointer and return its underlying struct.
	v := reflect.Indirect(val)

	return hasher.hashStruct(v)
}

func (hasher *objectHasher) hashRepeatedField(v reflect.Value, sf reflect.StructField, props *proto.Properties) ([]byte, error) {
	b := new(bytes.Buffer)
	for j := 0; j < v.Len(); j++ {
		elem := v.Index(j)
		if elem.Kind() == reflect.Ptr && elem.IsNil() {
			return nil, errors.New("got a nil message in a repeated field, which is invalid")
		}

		h, err := hasher.hashValue(elem, reflect.StructField{}, props)
		if err != nil {
			return nil, err
		}
		b.Write(h[:])
	}
	return hash(listIdentifier, b.Bytes())
}

func (hasher *objectHasher) hashMap(v reflect.Value, sf reflect.StructField, props *proto.Properties) ([]byte, error) {
	mapHashEntries := make([]hashEntry, v.Len())
	n := 0

	keyTag := sf.Tag.Get("protobuf_key")
	keyProps := new(proto.Properties)
	keyProps.Parse(keyTag)

	valTag := sf.Tag.Get("protobuf_val")
	valProps := new(proto.Properties)
	valProps.Parse(valTag)

	keys := v.MapKeys()
	for _, key := range keys {
		val := v.MapIndex(key)

		if val.Kind() == reflect.Ptr && val.IsNil() {
			return nil, errors.New("got a nil message in a map field, which is invalid")
		}

		// Hash the key.
		khash, err := hasher.hashValue(key, reflect.StructField{}, keyProps)
		if err != nil {
			return nil, err
		}
		mapHashEntries[n].khash = khash

		// Hash the value.
		vhash, err := hasher.hashValue(val, reflect.StructField{}, valProps)
		if err != nil {
			return nil, err
		}
		mapHashEntries[n].vhash = vhash

		n++
	}

	sort.Sort(byKHash(mapHashEntries))
	h := new(bytes.Buffer)
	for _, e := range mapHashEntries {
		h.Write(e.khash[:])
		h.Write(e.vhash[:])
	}
	return hash(mapIdentifier, h.Bytes())
}

// hashStruct hashes the struct objects of dereferenced proto messages.
//
// All proto messages are represented as pointers to structs. This method is
// used to calculate a proto message's hash by passing it the reflect.Value of
// the dererferenced message object.
func (hasher *objectHasher) hashStruct(sv reflect.Value) ([]byte, error) {
	name, ok := CheckWellKnownType(sv)
	if ok {
		return hasher.hashWellKnownType(name, sv)
	}

	if isExtendable(sv) {
		return nil, errors.New("extendable messages cannot be hashed reliably")
	}

	st := sv.Type()
	sprops := proto.GetProperties(st)

	structHashEntries := make([]hashEntry, sv.NumField())
	for i := 0; i < sv.NumField(); i++ {
		var entry hashEntry
		var err error

		v := sv.Field(i)
		sf := st.Field(i)

		// Ignore content-independent "XXX_" fields.
		if isContentIndependentField(v, sf) {
			continue
		}

		// Ignore unset fields (and empty proto3 scalar fields).
		unset, err := isUnset(v, sf)
		if err != nil {
			return nil, err
		}
		if unset {
			continue
		}

		if err = failIfUnsupported(v, sf); err != nil {
			return nil, err
		}

		if isAOneOfField(v, sf) {
			entry, err = hasher.hashOneOf(v, sf, sprops.Prop[i])
		} else {
			entry, err = hasher.hashStructField(v, sf, sprops.Prop[i])
		}
		if err != nil {
			return nil, err
		}

		structHashEntries[i] = entry
	}

	sort.Sort(byKHash(structHashEntries))
	h := new(bytes.Buffer)
	for _, e := range structHashEntries {
		h.Write(e.khash[:])
		h.Write(e.vhash[:])
	}

	identifier := mapIdentifier
	if hasher.messageIdentifier != "" {
		identifier = hasher.messageIdentifier
	}
	return hash(identifier, h.Bytes())
}

// hashValue returns the hash of an arbitrary proto field value.
//
// Note that the StructField argument is only used for types that can only
// exist within structs (ie. repeated fields and maps). Therefore, when the
// value does not exist within a struct, it is safe to call this function with
// an empty StructField (ie. reflect.StructField{}).
func (hasher *objectHasher) hashValue(v reflect.Value, sf reflect.StructField, props *proto.Properties) ([]byte, error) {
	switch v.Kind() {
	case reflect.Struct:
		return hasher.hashStruct(v)
	case reflect.Map:
		return hasher.hashMap(v, sf, props)
	case reflect.Slice:
		if props.Repeated {
			return hasher.hashRepeatedField(v, sf, props)
		}

		// If it's not a repeated field, then it must be []byte.
		return hashBytes(v.Bytes())
	case reflect.String:
		return hashUnicode(v.String())
	case reflect.Float32, reflect.Float64:
		return hashFloat(v.Float())
	case reflect.Int32, reflect.Int64:
		// This also includes enums, which are represented as integers.
		if hasher.enumsAsStrings && props.Enum != "" {
			str, err := stringify(v)
			if err != nil {
				return nil, err
			}
			return hashUnicode(str)
		}
		return hashInt64(v.Int())
	case reflect.Uint32, reflect.Uint64:
		return hashUint64(v.Uint())
	case reflect.Bool:
		return hashBool(v.Bool())
	case reflect.Ptr:
		// We know that this is not a null pointer because unset values (incl. null
		// pointer) get skipped and should not get hashed.
		return hasher.hashValue(reflect.Indirect(v), sf, props)
	default:
		return nil, fmt.Errorf("Unsupported type: %T", v)
	}
}

func (hasher *objectHasher) hashStructField(v reflect.Value, sf reflect.StructField, props *proto.Properties) (hashEntry, error) {
	var err error
	var khash []byte
	var vhash []byte

	if props.Required {
		return hashEntry{}, errors.New("required fields are not allowed because they're bad for backwards compatibility")
	}

	if props.HasDefault {
		return hashEntry{}, errors.New("fields with explicit defaults are not allowed because they're bad for backwards compatibility")
	}

	// Hash the tag.
	if hasher.fieldNamesAsKeys {
		khash, err = hashUnicode(props.OrigName)
	} else {
		khash, err = hashInt64(int64(props.Tag))
	}
	if err != nil {
		return hashEntry{}, err
	}

	// Hash the value.
	vhash, err = hasher.hashValue(v, sf, props)
	if err != nil {
		return hashEntry{}, err
	}

	return hashEntry{khash: khash, vhash: vhash}, nil
}

func (hasher *objectHasher) hashOneOf(v reflect.Value, sf reflect.StructField, props *proto.Properties) (hashEntry, error) {
	// A oneof field is an interface which contains a pointer to an inner struct that contains the value.
	fieldPointer := v.Elem()                      // Get the pointer to the inner struct.
	innerStruct := reflect.Indirect(fieldPointer) // Get the inner struct.

	// This check protects innerStruct.Field(0) from panicing.
	if innerStruct.Kind() != reflect.Struct || innerStruct.NumField() != 1 {
		return hashEntry{}, fmt.Errorf("unsupported interface type: %T. Expected it to be a oneof field", v)
	}
	innerValue := innerStruct.Field(0) // Get the inner value.

	// Check if the message is malformed.
	if innerValue.Kind() == reflect.Ptr && innerValue.IsNil() {
		return hashEntry{}, errors.New("got a nil message as a value of a oneof field, which is invalid")
	}

	// Parse the field's properties.
	// Oneof inner structs are defined to have a single field with the "protobuf" tag set.
	innerFd := innerStruct.Type().Field(0)
	innerTag := innerFd.Tag.Get("protobuf")
	innerProps := new(proto.Properties)
	innerProps.Parse(innerTag)

	// The inner field (which is a struct field) should never be considered unset
	// even if the value is a zero value.
	return hasher.hashStructField(innerValue, innerFd, innerProps)
}
