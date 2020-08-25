package main

import (
	"github.com/ahawkins/terraform-provider-sops/sops"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sops.Provider,
	})
}
