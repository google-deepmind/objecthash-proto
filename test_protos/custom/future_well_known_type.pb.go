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

package custom

// FutureWellKnownType is a manually created mock proto that can be used in
// tests as an unrecognized (or new) well known type.
//
// Note that this is not registered with the proto library (using init) to keep
// things simple.
type FutureWellKnownType struct {
}

func (m *FutureWellKnownType) Reset()                    { *m = FutureWellKnownType{} }
func (m *FutureWellKnownType) String() string            { return "FutureWellKnownType" }
func (*FutureWellKnownType) ProtoMessage()               {}
func (*FutureWellKnownType) Descriptor() ([]byte, []int) { return []byte{}, []int{0} }
func (*FutureWellKnownType) XXX_WellKnownType() string   { return "FutureWellKnownType" }

// The following line is used to prevent linters from running on this file:
// Code generated manually. DO NOT EDIT.
