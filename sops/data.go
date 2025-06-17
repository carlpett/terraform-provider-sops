package sops

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type summaryError struct {
	Summary string
	Err     error
}

func (e summaryError) Error() string {
	return fmt.Sprintf("%s: %s", e.Summary, e.Err.Error())
}

func newSummaryError(summary string, err error) summaryError {
	return summaryError{
		Summary: summary,
		Err:     err,
	}
}

func getFileData(sourceFile types.String, inputType types.String) (data map[string]string, raw string, err error) {
	sourceFileValue := sourceFile.ValueString()
	content, err := os.ReadFile(sourceFileValue)
	if err != nil {
		return nil, "", newSummaryError("Error reading file", err)
	}

	var format string
	if !inputType.IsNull() {
		format = inputType.ValueString()
	} else {
		switch ext := path.Ext(sourceFileValue); ext {
		case ".json":
			format = "json"
		case ".yaml", ".yml":
			format = "yaml"
		case ".env":
			format = "dotenv"
		case ".ini":
			format = "ini"
		default:
			return nil, "", newSummaryError("Unknown file type", fmt.Errorf("Don't know how to decode file with extension %s, set input_type as appropriate", ext))
		}
	}

	if err := validateInputType(format); err != nil {
		return nil, "", newSummaryError("Invalid input type", err)
	}

	data, raw, err = readData(content, format)
	if err != nil {
		return nil, "", newSummaryError("Error reading data", err)
	}
	return data, raw, nil
}

func getExternalData(source types.String, inputType types.String) (data map[string]string, raw string, err error) {
	content, err := io.ReadAll(strings.NewReader(source.ValueString()))
	if err != nil {
		return nil, "", newSummaryError("Error reading source", err)
	}

	format := inputType.ValueString()
	if err := validateInputType(format); err != nil {
		return nil, "", newSummaryError("Invalid input type", err)
	}

	data, raw, err = readData(content, format)
	if err != nil {
		return nil, "", newSummaryError("Error reading data", err)
	}

	return data, raw, nil
}
