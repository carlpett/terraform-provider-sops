package ini

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	input := []byte(`; Comment!
rootKey = foo
[Some section]
example_key = example_value`)
	expectedOutput := map[string]interface{}{
		"rootKey": "foo",
		"Some section": map[string]interface{}{
			"example_key": "example_value",
		},
	}
	var data map[string]interface{}
	err := Unmarshal(input, &data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedOutput, data) {
		t.Errorf("Unexpected output, expected %v, got %v", expectedOutput, data)
	}
}
