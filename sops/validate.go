package sops

import "fmt"

// validateInputType ensures that we can decode the input
func validateInputType(inputType string) error {
	switch inputType {
	case "json":
		return nil
	case "yaml", "yml":
		return nil
	case "raw":
		return nil
	default:
		return fmt.Errorf("Don't know how to decode file with input type %s, set input_type to json, yaml or raw as appropriate", inputType)
	}
}
