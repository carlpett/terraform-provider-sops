package sops

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const configTestDataSourceSopsExternalMetadata_basic = `
data "sops_external_metadata" "test" {
  source     = file("%s/test-fixtures/basic.yaml")
  input_type = "yaml"
}`

func TestDataSourceSopsExternalMetadata_basic(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsExternalMetadata_basic, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external_metadata.test", "last_modified", "2019-04-26T18:43:59Z"),
					resource.TestCheckResourceAttr("data.sops_external_metadata.test", "last_modified_unix", "1556304239"),
				),
			},
		},
	})
}
