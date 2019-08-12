package sops

import "testing"

func testValidateInputType(inputType string, t *testing.T) {
	err := validateInputType(inputType)
	if err != nil {
		t.Errorf("Failed to validate input type %s", inputType)
	}
}

func TestValidateInputType_yaml(t *testing.T) {
	inputType := "yaml"
	testValidateInputType(inputType, t)
}

func TestValidateInputType_json(t *testing.T) {
	inputType := "json"
	testValidateInputType(inputType, t)
}

func TestValidateInputType_raw(t *testing.T) {
	inputType := "raw"
	testValidateInputType(inputType, t)
}

func TestValidateInputType_bad(t *testing.T) {
	inputType := "tf"
	err := validateInputType(inputType)
	if err == nil {
		t.Errorf("Failed to validate input type %s, expected to be invalid but was valid", inputType)
	}
}
