package sops

import (
	"reflect"
	"testing"
)

func TestFlattening(t *testing.T) {
	tc := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]string
	}{
		{
			name: "all data types become strings",
			input: map[string]interface{}{
				"a_string":   "foo",
				"an_integer": 12,
				"a_bool":     true,
				"a_float":    1.1,
			},
			expected: map[string]string{
				"a_string":   "foo",
				"an_integer": "12",
				"a_bool":     "true",
				"a_float":    "1.1",
			},
		},
		{
			name: "dicts are unnested",
			input: map[string]interface{}{
				"a_dict": map[string]interface{}{"foo": "bar"},
			},
			expected: map[string]string{
				"a_dict.foo": "bar",
			},
		},
		{
			name: "lists are unpacked with index keys",
			input: map[string]interface{}{
				"a_list": []interface{}{1, 2},
			},
			expected: map[string]string{
				"a_list.0": "1",
				"a_list.1": "2",
			},
		},
		{
			name: "deep nesting",
			/*
				This test corresponds to this yaml structure:
				foo:
				- a: 1
				  b:
				    c:
				    - d: 2
			*/
			input: map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"a": 1,
						"b": map[string]interface{}{
							"c": []interface{}{
								map[string]interface{}{"d": 2},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"foo.0.a":       "1",
				"foo.0.b.c.0.d": "2",
			},
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			output := flatten(c.input)
			if !reflect.DeepEqual(c.expected, output) {
				t.Errorf("Unexpected flattening output, expected %v, got %v", c.expected, output)
			}
		})
	}
}
