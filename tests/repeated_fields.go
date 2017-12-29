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

package tests

import (
	"testing"

	pb2_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto2"
	pb3_latest "github.com/deepmind/objecthash-proto/test_protos/generated/latest/proto3"

	"github.com/golang/protobuf/proto"
)

func TestRepeatedFields(t *testing.T, hashers ProtoHashers) {
	hasher := hashers.StringPreferringHasher

	testCases := []testCase{
		///////////////////
		//  Empty lists. //
		///////////////////

		// Empty repeated fields are ignored when taking a protobuf's objecthash.
		// This is the case for both Proto2 and Proto3.
		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{
					BoolField:       []bool{},
					BytesField:      [][]byte{},
					DoubleField:     []float64{},
					Fixed32Field:    []uint32{},
					Fixed64Field:    []uint64{},
					FloatField:      []float32{},
					Int32Field:      []int32{},
					Int64Field:      []int64{},
					Sfixed32Field:   []int32{},
					Sfixed64Field:   []int64{},
					Sint32Field:     []int32{},
					Sint64Field:     []int64{},
					StringField:     []string{},
					Uint32Field:     []uint32{},
					Uint64Field:     []uint64{},
					SimpleField:     []*pb2_latest.Simple{},
					RepetitiveField: []*pb2_latest.Repetitive{},
					SingletonField:  []*pb2_latest.Singleton{},
				},
				&pb3_latest.Repetitive{
					BoolField:       []bool{},
					BytesField:      [][]byte{},
					DoubleField:     []float64{},
					Fixed32Field:    []uint32{},
					Fixed64Field:    []uint64{},
					FloatField:      []float32{},
					Int32Field:      []int32{},
					Int64Field:      []int64{},
					Sfixed32Field:   []int32{},
					Sfixed64Field:   []int64{},
					Sint32Field:     []int32{},
					Sint64Field:     []int64{},
					StringField:     []string{},
					Uint32Field:     []uint32{},
					Uint64Field:     []uint64{},
					SimpleField:     []*pb3_latest.Simple{},
					RepetitiveField: []*pb3_latest.Repetitive{},
					SingletonField:  []*pb3_latest.Singleton{},
				},
			},
			equivalentJsonString: "{}",
			equivalentObject:     map[string]interface{}{},
			expectedHashString:   "18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4",
		},

		//////////////////////////
		//  Lists with strings. //
		//////////////////////////
		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{StringField: []string{"foo"}},
				&pb3_latest.Repetitive{StringField: []string{"foo"}},
			},
			equivalentJsonString: "{\"string_field\": [\"foo\"]}",
			equivalentObject:     map[string][]string{"string_field": []string{"foo"}},
			expectedHashString:   "54c0b7c6e7c9ff0bb6076a2caeccbc96fad77f49b17b7ec9bc17dfe98a7b343e",
		},

		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{StringField: []string{"foo", "bar"}},
				&pb3_latest.Repetitive{StringField: []string{"foo", "bar"}},
			},
			equivalentJsonString: "{\"string_field\": [\"foo\", \"bar\"]}",
			equivalentObject:     map[string][]string{"string_field": []string{"foo", "bar"}},
			expectedHashString:   "a971a061d199ddf37a365d617f9cd4530efb15e933e0dbaf6602b2908b792056",
		},

		///////////////////////
		//  Lists with ints. //
		///////////////////////

		// JSON treats all numbers as floats, so it is not possible to have an equivalent JSON string.

		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{Int64Field: []int64{-2, -1, 0, 1, 2}},
				&pb3_latest.Repetitive{Int64Field: []int64{-2, -1, 0, 1, 2}},
			},
			equivalentObject:   map[string][]int64{"int64_field": []int64{-2, -1, 0, 1, 2}},
			expectedHashString: "44e78ff73bdf5d0da5141e110b22bab240483ba17c40f83553a0e6bbfa671e22",
		},

		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{Int64Field: []int64{123456789012345, 678901234567890}},
				&pb3_latest.Repetitive{Int64Field: []int64{123456789012345, 678901234567890}},
			},
			equivalentObject:   map[string][]int64{"int64_field": []int64{123456789012345, 678901234567890}},
			expectedHashString: "b0ce1b7dfa71b33a16571fea7f3f27341bf5980b040e9d949a8019f3143ecbc7",
		},

		/////////////////////////
		//  Lists with floats. //
		/////////////////////////
		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{FloatField: []float32{-2, -1, 0, 1, 2}},
				&pb3_latest.Repetitive{FloatField: []float32{-2, -1, 0, 1, 2}},
			},
			equivalentJsonString: "{\"float_field\": [-2, -1, 0, 1, 2]}",
			equivalentObject:     map[string][]float32{"float_field": []float32{-2, -1, 0, 1, 2}},
			expectedHashString:   "68b2552f2f33b5dd38c9be0aeee127170c86d8d2b3ab7daebdc2ea124226593f",
		},

		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{FloatField: []float32{1, 2, 3}},
				&pb3_latest.Repetitive{FloatField: []float32{1, 2, 3}},
			},
			equivalentJsonString: "{\"float_field\": [1, 2, 3]}",
			equivalentObject:     map[string][]float32{"float_field": []float32{1, 2, 3}},
			expectedHashString:   "f26c1502d1f9f7bf672cf669290348f9bfdea0af48261f2822aad01927fe1749",
		},

		{
			protos: []proto.Message{
				&pb2_latest.Repetitive{DoubleField: []float64{123456789012345, 678901234567890}},
				&pb3_latest.Repetitive{DoubleField: []float64{123456789012345, 678901234567890}},
			},
			equivalentJsonString: "{\"double_field\": [123456789012345, 678901234567890]}",
			equivalentObject:     map[string][]float64{"double_field": []float64{123456789012345, 678901234567890}},
			expectedHashString:   "067d25d39b8514b6b905e0eba2d19242bcf4441e2367527dbceac7a9dd0108a0",
		},
	}

	for _, tc := range testCases {
		if err := tc.check(hasher); err != nil {
			t.Errorf("%s", err)
		}
	}
}
