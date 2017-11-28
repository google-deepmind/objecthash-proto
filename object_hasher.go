package protohash

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/golang/protobuf/proto"
)

// ObjectHasher is a configurable object for hashing protocol buffer objects.
type ObjectHasher struct {
	// Whether to hash enum values as strings, as opposed to as integer values.
	EnumsAsStrings bool

	// Whether to use the proto field name its key, as opposed to using the tag
	// number as the key.
	FieldNamesAsKeys bool

	// Whether to hash proto messages as maps, as opposed to using a separate
	// hash identifier for them. Enabling this will make the ObjectHash of a
	// proto message equivalent to the ObjectHash of an equivalent map object.
	TreatMessagesAsMaps bool
}

// `HashProto` returns the object hash of a given protocol buffer message.
func (hasher *ObjectHasher) HashProto(pb proto.Message) ([]byte, error) {
	val := reflect.ValueOf(pb)
	// The Go generated code defines protos as pointers to structs.
	// This means that the `IsNil` check here is safe to make and wont panic.
	if pb == nil || val.IsNil() {
		return hashNil()
	}

	// Dereference the proto pointer and return its underlying struct.
	v := reflect.Indirect(val)

	return hasher.hashStruct(v)
}

func (hasher *ObjectHasher) hashRepeatedField(v reflect.Value, sf reflect.StructField, props *proto.Properties) ([]byte, error) {
	b := new(bytes.Buffer)
	for j := 0; j < v.Len(); j++ {

		elem := v.Index(j)
		if elem.Kind() == reflect.Ptr && elem.IsNil() {
			return nil, errors.New("Got a nil message in a repeated field, which is invalid.")
		}

		h, err := hasher.hashValue(elem, reflect.StructField{}, props)
		if err != nil {
			return nil, err
		}
		b.Write(h[:])
	}
	return hash(listIdentifier, b.Bytes())
}

func (hasher *ObjectHasher) hashMap(v reflect.Value, sf reflect.StructField, props *proto.Properties) ([]byte, error) {
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
			return nil, errors.New("Got a nil message in a map field, which is invalid.")
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

func (hasher *ObjectHasher) hashStruct(sv reflect.Value) ([]byte, error) {
	if isAny(sv) {
		return nil, errors.New("google.protobuf.Any messages cannot be hashed reliably.")
	}

	if isExtendable(sv) {
		return nil, errors.New("Extendable messages cannot be hashed reliably.")
	}

	st := sv.Type()
	sprops := proto.GetProperties(st)

	structHashEntries := make([]hashEntry, sv.NumField())
	for i := 0; i < sv.NumField(); i++ {
		var entry hashEntry
		var err error

		v := sv.Field(i)
		sf := st.Field(i)

		// Ignore unused/empty fields.
		empty, err := isEmpty(v)
		if err != nil {
			return nil, err
		}
		if empty {
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

	identifier := protoMessageIdentifier
	if hasher.TreatMessagesAsMaps {
		identifier = mapIdentifier
	}
	return hash(identifier, h.Bytes())
}

// `hashValue` returns the hash of an arbitrary proto field value.
//
// Note that the StructField argument is only used for types that can only
// exist within structs (ie. repeated fields and maps). Therefore, when the
// value does not exist within a struct, it is safe to call this function with
// an empty StructField (ie. `reflect.StructField{}`).
func (hasher *ObjectHasher) hashValue(v reflect.Value, sf reflect.StructField, props *proto.Properties) ([]byte, error) {
	switch v.Kind() {
	case reflect.Struct:
		return hasher.hashStruct(v)
	case reflect.Map:
		return hasher.hashMap(v, sf, props)
	case reflect.Slice:
		if props.Repeated {
			return hasher.hashRepeatedField(v, sf, props)
		} else {
			// If it's not a repeated field, then it must be []byte.
			return hashBytes(v.Bytes())
		}
	case reflect.String:
		return hashUnicode(v.String())
	case reflect.Float32, reflect.Float64:
		return hashFloat(v.Float())
	case reflect.Int32, reflect.Int64:
		// This also includes enums, which are represented as integers.
		if hasher.EnumsAsStrings && props.Enum != "" {
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
		// We know that this is not a null pointer because empty values (incl. null
		// pointer) get skipped and should not get hashed.
		return hasher.hashValue(reflect.Indirect(v), sf, props)
	default:
		return nil, fmt.Errorf("Unsupported type: %T", v)
	}
}

func (hasher *ObjectHasher) hashStructField(v reflect.Value, sf reflect.StructField, props *proto.Properties) (hashEntry, error) {
	var err error
	var khash []byte
	var vhash []byte

	// Hash the tag.
	if hasher.FieldNamesAsKeys {
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

func (hasher *ObjectHasher) hashOneOf(v reflect.Value, sf reflect.StructField, props *proto.Properties) (hashEntry, error) {
	// A oneof field is an interface which contains a pointer to an inner struct that contains the value.
	fieldPointer := v.Elem()                      // Get the pointer to the inner struct.
	innerStruct := reflect.Indirect(fieldPointer) // Get the inner struct.

	// This check protects `innerStruct.Field(0)` from panicing.
	if innerStruct.Kind() != reflect.Struct || innerStruct.NumField() != 1 {
		return hashEntry{}, fmt.Errorf("Unsupported interface type: %T. Expected it to be a oneof field.", v)
	}
	innerValue := innerStruct.Field(0) // Get the inner value.

	// Check if the message is malformed.
	if innerValue.Kind() == reflect.Ptr && innerValue.IsNil() {
		return hashEntry{}, errors.New("Got a nil message as a value of a oneof field, which is invalid.")
	}

	// Parse the field's properties.
	// Oneof inner structs are defined to have a single field with the "protobuf" tag set.
	innerFd := innerStruct.Type().Field(0)
	innerTag := innerFd.Tag.Get("protobuf")
	innerProps := new(proto.Properties)
	innerProps.Parse(innerTag)

	// The inner field (which is a struct field) cannot be considered empty even
	// if the value is a zero value.
	return hasher.hashStructField(innerValue, innerFd, innerProps)
}
