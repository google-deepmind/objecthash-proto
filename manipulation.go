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
)

// stringify returns a string representation of a reflect.Value object.
func stringify(v reflect.Value) (string, error) {
	if !v.IsValid() {
		return "", errors.New("encountered a null pointer")
	}

	if !v.CanInterface() {
		return "", errors.New("encountered an unexported struct field")
	}

	stringerValue, ok := v.Interface().(fmt.Stringer)
	if ok {
		return stringerValue.String(), nil
	}

	return "", fmt.Errorf("failed to represent value '%v' as a string because it does not have a 'String()' method", v)
}
