package main

import (
	"github.com/carlpett/terraform-sops/sops"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sops.Provider,
	})
}
