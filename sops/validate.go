package sops

import "fmt"

var validTypes = map[string]bool{
	"json":   true,
	"yaml": true,
	"dotenv": true,
	"ini": true,
	"raw": true,
}

// validateInputType ensures that we can decode the input
func validateInputType(inputType string) error {
	if _, ok := validTypes[inputType]; ok {
		return nil
	}
	return fmt.Errorf("Don't know how to decode file with input type %s, set input_type as appropriate", inputType)
}
