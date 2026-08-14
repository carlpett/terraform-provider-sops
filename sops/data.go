package sops

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

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

func getFileData(sourceFile types.String, inputType types.String) (sopsData, error) {
	sourceFileValue := sourceFile.ValueString()
	content, err := os.ReadFile(sourceFileValue)
	if err != nil {
		return sopsData{}, newSummaryError("Error reading file", err)
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
			return sopsData{}, newSummaryError("Unknown file type", fmt.Errorf("Don't know how to decode file with extension %s, set input_type as appropriate", ext))
		}
	}

	if err := validateInputType(format); err != nil {
		return sopsData{}, newSummaryError("Invalid input type", err)
	}

	result, err := readData(content, format)
	if err != nil {
		return sopsData{}, newSummaryError("Error reading data", err)
	}
	return result, nil
}

func getExternalData(source types.String, inputType types.String) (sopsData, error) {
	content, err := io.ReadAll(strings.NewReader(source.ValueString()))
	if err != nil {
		return sopsData{}, newSummaryError("Error reading source", err)
	}

	format := inputType.ValueString()
	if err := validateInputType(format); err != nil {
		return sopsData{}, newSummaryError("Invalid input type", err)
	}

	result, err := readData(content, format)
	if err != nil {
		return sopsData{}, newSummaryError("Error reading data", err)
	}

	return result, nil
}

// getFileMetadata mirrors getFileData, but reads only the sops metadata of the
// file — it returns the lastmodified timestamp recorded there, and never needs
// access to the data key.
func getFileMetadata(sourceFile types.String, inputType types.String) (time.Time, error) {
	sourceFileValue := sourceFile.ValueString()
	content, err := os.ReadFile(sourceFileValue)
	if err != nil {
		return time.Time{}, newSummaryError("Error reading file", err)
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
			return time.Time{}, newSummaryError("Unknown file type", fmt.Errorf("Don't know how to decode file with extension %s, set input_type as appropriate", ext))
		}
	}

	if err := validateInputType(format); err != nil {
		return time.Time{}, newSummaryError("Invalid input type", err)
	}

	lastModified, err := readLastModified(content, format)
	if err != nil {
		return time.Time{}, newSummaryError("Error reading sops metadata", err)
	}
	return lastModified, nil
}

// getExternalMetadata mirrors getExternalData, but reads only the sops metadata
// of the encrypted string — it returns the lastmodified timestamp recorded
// there, and never needs access to the data key.
func getExternalMetadata(source types.String, inputType types.String) (time.Time, error) {
	content := []byte(source.ValueString())

	format := inputType.ValueString()
	if err := validateInputType(format); err != nil {
		return time.Time{}, newSummaryError("Invalid input type", err)
	}

	lastModified, err := readLastModified(content, format)
	if err != nil {
		return time.Time{}, newSummaryError("Error reading sops metadata", err)
	}
	return lastModified, nil
}
