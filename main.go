package main

import (
	"github.com/carlpett/terraform-provider-sops/sops"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sops.Provider,
	})
}
