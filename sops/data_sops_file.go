package sops

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"io/ioutil"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFileRead,
		Schema: map[string]*schema.Schema{
			"input_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_file": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"data": &schema.Schema{
				Type:      schema.TypeMap,
				Computed:  true,
				Sensitive: true,
			},
			"raw": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get the environment variables from the provider configuration
	envVars, ok := meta.(map[string]interface{})
	if !ok {
		return diag.Errorf("Unable to get provider configuration")
	}
	// Set the environment variables
	for key, value := range envVars {
		if strValue, ok := value.(string); ok {
			os.Setenv(key, strValue)
		}
	}

	sourceFile := d.Get("source_file").(string)
	content, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return diag.FromErr(err)
	}

	var format string
	if inputType := d.Get("input_type").(string); inputType != "" {
		format = inputType
	} else {
		switch ext := path.Ext(sourceFile); ext {
		case ".json":
			format = "json"
		case ".yaml", ".yml":
			format = "yaml"
		case ".env":
			format = "dotenv"
		case ".ini":
			format = "ini"
		default:
			return diag.Errorf("don't know how to decode file with extension %s, set input_type to json, yaml or raw as appropriate", ext)
		}
	}

	if err := validateInputType(format); err != nil {
		return diag.FromErr(err)
	}

	if err := readData(content, format, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
