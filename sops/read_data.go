package sops

import (
	"encoding/json"
	"fmt"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/decrypt"
	"gopkg.in/yaml.v3"

	"github.com/carlpett/terraform-provider-sops/sops/internal/dotenv"
	"github.com/carlpett/terraform-provider-sops/sops/internal/ini"
)

func readData(content []byte, format string) (map[string]string, string, error) {
	cleartext, err := decrypt.Data(content, format)
	if userErr, ok := err.(sops.UserError); ok {
		err = fmt.Errorf(userErr.UserError())
	}
	if err != nil {
		return nil, "", fmt.Errorf("Error decrypting sops file: %w", err)
	}

	var data map[string]interface{}
	switch format {
	case "json":
		err = json.Unmarshal(cleartext, &data)
	case "yaml":
		err = yaml.Unmarshal(cleartext, &data)
	case "dotenv":
		err = dotenv.Unmarshal(cleartext, &data)
	case "ini":
		err = ini.Unmarshal(cleartext, &data)
	}
	if err != nil {
		return nil, "", fmt.Errorf("Error parsing decrypted data: %w", err)
	}

	return flatten(data), string(cleartext), nil
}
