package stores

import (
	"errors"
	"reflect"
	"testing"
)

func TestConvertStructToMap(t *testing.T) {
	type ComplexType struct {
		Val int
	}

	type InterfaceType struct {
		Val interface{}
	}

	type NestedStruct struct {
		A ComplexType
		B ComplexType
	}

	cases := []struct {
		desc        string
		input       interface{}
		expected    map[string]interface{}
		expectedErr error
	}{
		{
			desc: "slice field with primitives",
			input: struct {
				Foo []int
			}{
				Foo: []int{1, 2, 3},
			},
			expected: map[string]interface{}{
				"Foo": []int{1, 2, 3},
			},
			expectedErr: nil,
		},
		{
			desc: "slice field with complex types",
			input: struct {
				Foo []ComplexType
			}{
				Foo: []ComplexType{
					{1}, {2}, {3},
				},
			},
			expected: map[string]interface{}{
				"Foo": []interface{}{
					map[string]interface{}{"Val": 1},
					map[string]interface{}{"Val": 2},
					map[string]interface{}{"Val": 3},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "map field with primitives",
			input: struct {
				Foo map[string]string
			}{
				Foo: map[string]string{
					"Foo": "Bar",
					"Biz": "Baz",
				},
			},
			expected: map[string]interface{}{
				"Foo": map[string]string{
					"Foo": "Bar",
					"Biz": "Baz",
				},
			},
			expectedErr: nil,
		},
		{
			desc: "map field with complex types",
			input: struct {
				Foo map[string]ComplexType
			}{
				Foo: map[string]ComplexType{
					"Biz": {1},
					"Baz": {2},
				},
			},
			expected: map[string]interface{}{
				"Foo": map[string]interface{}{
					"Biz": map[string]interface{}{"Val": 1},
					"Baz": map[string]interface{}{"Val": 2},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "nested structures",
			input: struct {
				Foo NestedStruct
			}{
				Foo: NestedStruct{
					A: ComplexType{Val: 1},
					B: ComplexType{Val: 2},
				},
			},
			expected: map[string]interface{}{
				"Foo": map[string]interface{}{
					"A": map[string]interface{}{"Val": 1},
					"B": map[string]interface{}{"Val": 2},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "slice of interfaces",
			input: struct {
				Foo []interface{}
			}{
				Foo: []interface{}{
					ComplexType{1},
					ComplexType{2},
					ComplexType{3},
				},
			},
			expected: map[string]interface{}{
				"Foo": []interface{}{
					map[string]interface{}{"Val": 1},
					map[string]interface{}{"Val": 2},
					map[string]interface{}{"Val": 3},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			result, err := ConvertStructToMap(tc.input)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("unexpected error got '%v' wanted '%v'", err, tc.expectedErr)
			}

			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("unexpected result got '%v' wanted '%v'", result, tc.expected)
			}
		})
	}
}
