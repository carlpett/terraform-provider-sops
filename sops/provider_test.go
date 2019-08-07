package sops

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-local/local"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccLocalProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccLocalProvider = local.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sops":  testAccProvider,
		"local": testAccLocalProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
