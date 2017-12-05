package protohash

import (
	"errors"
	"fmt"
	"reflect"
)

// `stringify` returns a string representation of a `reflect.Value` object.
func stringify(v reflect.Value) (string, error) {
	if !v.IsValid() {
		return "", errors.New("Encountered a null pointer.")
	}

	if !v.CanInterface() {
		return "", errors.New("Encountered an unexported struct field.")
	}

	stringerValue, ok := v.Interface().(fmt.Stringer)
	if ok {
		return stringerValue.String(), nil
	} else {
		return "", fmt.Errorf("Failed to represent value '%v' as a string because it does not have a `String()` method", v)
	}
}
