// Copyright 2018 the Service Broker Project Authors.
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

package broker

import (
	"errors"
	"reflect"
	"testing"

	"github.com/cloudfoundry/cloud-service-broker/pkg/validation"
)

func TestBrokerVariable_ToSchema(t *testing.T) {
	cases := map[string]struct {
		BrokerVar BrokerVariable
		Expected  map[string]any
	}{
		"blank": {
			BrokerVariable{}, map[string]any{},
		},
		"enums get copied": {
			BrokerVariable{Enum: map[any]string{"a": "description", "b": "description"}},
			map[string]any{
				"enum": []any{"a", "b"},
			},
		},
		"details are copied": {
			BrokerVariable{Details: "more information"},
			map[string]any{
				"description": "more information",
			},
		},
		"type is copied": {
			BrokerVariable{Type: JSONTypeString},
			map[string]any{
				"type": JSONTypeString,
			},
		},
		"default is copied": {
			BrokerVariable{Default: "some-value"},
			map[string]any{
				"default": "some-value",
			},
		},
		"template defaults": {
			BrokerVariable{Default: "${33}", Details: "Some value."},
			map[string]any{
				"description": `Some value. If you do not specify this field, it will be generated by the template "${33}"`,
			},
		},
		"full test": {
			BrokerVariable{
				FieldName: "full_test_field_name",
				Default:   "some-value",
				Type:      JSONTypeString,
				Details:   "more information",
				Enum:      map[any]string{"b": "description", "a": "description"},
				Constraints: map[string]any{
					"examples": []string{"SAMPLEA", "SAMPLEB"},
				},
				TFAttribute: "test.fake.field",
			},
			map[string]any{
				"title":        "Full Test Field Name",
				"default":      "some-value",
				"type":         JSONTypeString,
				"description":  "more information",
				"enum":         []any{"a", "b"},
				"examples":     []string{"SAMPLEA", "SAMPLEB"},
				"tf_attribute": "test.fake.field",
			},
		},
		"prohibit update is copied": {
			BrokerVariable{ProhibitUpdate: true},
			map[string]any{
				"prohibitUpdate": true,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.BrokerVar.ToSchema()
			if !reflect.DeepEqual(actual, tc.Expected) {
				t.Errorf("Expected ToSchema to be: %v, got: %v", tc.Expected, actual)
			}
		})
	}
}

func TestBrokerVariable_ValidateVariables(t *testing.T) {
	cases := map[string]struct {
		Parameters map[string]any
		Variables  []BrokerVariable
		Expected   error
	}{
		"nil params": {
			Parameters: nil,
			Variables:  nil,
			Expected:   errors.New("1 error(s) occurred: (root): Invalid type. Expected: object, given: null"),
		},
		"nil vars check": {
			Parameters: map[string]any{},
			Variables:  nil,
			Expected:   nil,
		},
		"integer": {
			Parameters: map[string]any{
				"test": 12,
			},
			Variables: []BrokerVariable{
				{
					Required:  true,
					FieldName: "test",
					Type:      JSONTypeInteger,
				},
			},
			Expected: nil,
		},
		"unexpected type": {
			Parameters: map[string]any{
				"test": "didn't see that coming",
			},
			Variables: []BrokerVariable{
				{
					Required:  true,
					FieldName: "test",
					Type:      JSONTypeInteger,
				},
			},
			Expected: errors.New("1 error(s) occurred: test: Invalid type. Expected: integer, given: string"),
		},
		"test constraints": {
			Parameters: map[string]any{
				"test": 0,
			},
			Variables: []BrokerVariable{
				{
					Required:  true,
					FieldName: "test",
					Type:      JSONTypeInteger,
					Constraints: validation.NewConstraintBuilder().
						Minimum(10).
						Build(),
				},
			},
			Expected: errors.New("1 error(s) occurred: test: Must be greater than or equal to 10"),
		},
		"test enum": {
			Parameters: map[string]any{
				"test": "not this one",
			},
			Variables: []BrokerVariable{
				{
					Required:  true,
					FieldName: "test",
					Type:      JSONTypeString,
					Enum: map[any]string{
						"one":      "it's either this one",
						"theother": "or this one",
					},
				},
			},
			Expected: errors.New(`1 error(s) occurred: test: test must be one of the following: "one", "theother"`),
		},
		"test missing": {
			Parameters: map[string]any{},
			Variables: []BrokerVariable{
				{
					Required:  true,
					FieldName: "test",
					Type:      JSONTypeString,
					Enum: map[any]string{
						"one":      "it's either this one",
						"theother": "or this one",
					},
				},
			},
			Expected: errors.New("1 error(s) occurred: (root): test is required"),
		},
		"test incorrect schema": {
			Parameters: map[string]any{},
			Variables: []BrokerVariable{
				{
					Type: "garbage",
				},
			},
			Expected: errors.New("has a primitive type that is NOT VALID -- given: /garbage/ Expected valid values are:[array boolean integer number null object string]"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := ValidateVariables(tc.Parameters, tc.Variables)
			if tc.Expected == nil {
				if actual != nil {
					t.Fatalf("Expected ValidateVariables not to raise an error but got %v", actual)
				}
			} else {
				if actual == nil {
					t.Fatalf("Expected ValidateVariables to be: %q, got: %v", tc.Expected.Error(), actual)
				}

				if actual.Error() != tc.Expected.Error() {
					t.Errorf("Expected ValidateVariables error to be: %q, got: %q", tc.Expected.Error(), actual.Error())
				}
			}
		})
	}
}

func TestBrokerVariable_Validate(t *testing.T) {
	cases := map[string]struct {
		Parameters map[string]any
		Variable   BrokerVariable
		Expected   error
	}{
		"valid fields": {
			Variable: BrokerVariable{
				FieldName:   "test",
				Details:     "test variable",
				Type:        JSONTypeInteger,
				TFAttribute: "type.name.attribute",
			},
			Expected: nil,
		},
		"blank fields": {
			Variable: BrokerVariable{
				FieldName: "",
				Details:   "",
			},
			Expected: errors.New("missing field(s): details, field_name"),
		},
		"invalid JSON type": {
			Variable: BrokerVariable{
				FieldName: "test",
				Details:   "test variable",
				Type:      "map",
			},
			Expected: errors.New("field must match '^(|object|boolean|array|number|string|integer)$': type"),
		},
		"invalid tf_attribute": {
			Variable: BrokerVariable{
				FieldName:   "test",
				Details:     "test variable",
				Type:        JSONTypeInteger,
				TFAttribute: "thisisnot.validtfattribute",
			},
			Expected: errors.New("field must match '^([-a-zA-Z0-9_-]*\\.[-a-zA-Z0-9_-]*){2}': tf_attribute"),
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual := tc.Variable.Validate()
			if tc.Expected == nil {
				if actual != nil {
					t.Fatalf("Expected ValidateVariables not to raise an error but got %v", actual)
				}
			} else {
				if actual == nil {
					t.Fatalf("Expected ValidateVariables to be: %q, got: %v", tc.Expected.Error(), actual)
				}

				if actual.Error() != tc.Expected.Error() {
					t.Errorf("Expected ValidateVariables error to be: %q, got: %q", tc.Expected.Error(), actual.Error())
				}
			}
		})
	}
}

func TestBrokerVariable_ApplyDefaults(t *testing.T) {
	cases := map[string]struct {
		Parameters map[string]any
		Variables  []BrokerVariable
		Expected   map[string]any
	}{
		"nil check": {
			Parameters: nil,
			Variables:  nil,
			Expected:   nil,
		},
		"simple": {
			Parameters: map[string]any{},
			Variables: []BrokerVariable{
				{
					FieldName: "test",
					Type:      JSONTypeInteger,
					Default:   123,
				},
			},
			Expected: map[string]any{
				"test": 123,
			},
		},
		"do not replace": {
			Parameters: map[string]any{
				"test": 123,
			},
			Variables: []BrokerVariable{
				{
					FieldName: "test",
					Type:      JSONTypeInteger,
					Default:   456,
				},
			},
			Expected: map[string]any{
				"test": 123,
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			ApplyDefaults(tc.Parameters, tc.Variables)

			if !reflect.DeepEqual(tc.Parameters, tc.Expected) {
				t.Errorf("Expected ValidateVariables to be: %v, got: %v", tc.Expected, tc.Parameters)
			}
		})
	}
}

func TestFieldNameToLabel(t *testing.T) {
	cases := map[string]struct {
		Field    string
		Expected string
	}{
		"snake_case":       {Field: "my_field", Expected: "My Field"},
		"kebab-case":       {Field: "kebab-case", Expected: "Kebab Case"},
		"dot.notation":     {Field: "dot.notation", Expected: "Dot Notation"},
		"uri":              {Field: "my_uri", Expected: "My URI"},
		"url":              {Field: "my_url", Expected: "My URL"},
		"id":               {Field: "my_id", Expected: "My ID"},
		"gb":               {Field: "size_gb", Expected: "Size GB"},
		"jdbc":             {Field: "jdbc", Expected: "JDBC"},
		"double_separator": {Field: "my__id", Expected: "My ID"},
		"single_char":      {Field: "a b c", Expected: "A B C"},
		"blank":            {Field: "", Expected: ""},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			field := fieldNameToLabel(tc.Field)
			if field != tc.Expected {
				t.Errorf("Expected: %q got %q", tc.Expected, field)
			}
		})
	}
}
