package sops

import "fmt"

// validateInputType ensures that we can decode the input
func validateInputType(inputType string) error {
	switch inputType {
	case "json":
		return nil
	case "yaml":
		return nil
	case "dotenv":
		return nil
	case "ini":
		return nil
	case "raw":
		return nil
	default:
		return fmt.Errorf("Don't know how to decode file with input type %s, set input_type to json, yaml, ini, dotenv or raw as appropriate", inputType)
	}
}
