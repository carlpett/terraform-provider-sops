package sops

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"sops_file":     dataSourceFile(),
			"sops_external": dataSourceExternal(),
		},
	}
}
