package sops

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/decrypt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"

	"github.com/carlpett/terraform-provider-sops/sops/internal/dotenv"
	"github.com/carlpett/terraform-provider-sops/sops/internal/ini"
)

var decryptMutex sync.Mutex

func readData(content []byte, format string, env types.Map) (map[string]string, string, error) {

	if !env.IsNull() && !env.IsUnknown() {
		rawMap := env.Elements()
		envMap := make(map[string]string, len(rawMap))
		for k, v := range rawMap {
			if vStr, ok := v.(types.String); ok && !vStr.IsNull() && !vStr.IsUnknown() {
				envMap[k] = vStr.ValueString()
			}
		}

		decryptMutex.Lock()
		defer decryptMutex.Unlock()

		for k, v := range envMap {
			if err := os.Setenv(k, v); err != nil {
				return nil, "", fmt.Errorf("failed to set environment variable %q: %w", k, err)
			}
		}

		defer func() {
			for k := range envMap {
				_ = os.Unsetenv(k)
			}
		}()
	}

	cleartext, err := decrypt.Data(content, format)
	if userErr, ok := err.(sops.UserError); ok {
		err = userErr
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
