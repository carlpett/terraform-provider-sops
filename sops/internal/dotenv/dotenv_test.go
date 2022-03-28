package dotenv

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	input := []byte(`# Comment!
password=P@ssw0rd`)
	expectedOutput := map[string]interface{}{
		"password": "P@ssw0rd",
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
