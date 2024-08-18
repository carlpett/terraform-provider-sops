package sops

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceExternal() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExternalRead,
		Schema: map[string]*schema.Schema{
			"input_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
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

func dataSourceExternalRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	source := d.Get("source").(string)
	content, err := ioutil.ReadAll(strings.NewReader(source))
	if err != nil {
		return diag.FromErr(err)
	}

	format := d.Get("input_type").(string)
	if err := validateInputType(format); err != nil {
		return diag.FromErr(err)
	}

	if err := readData(content, format, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
