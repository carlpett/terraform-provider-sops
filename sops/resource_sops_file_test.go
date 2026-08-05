package sops

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const configTestResourceSopsFile_dataWriteOnly = `
resource "sops_file" "test_data_wo" {
  destination_file = "%[1]s/test-fixtures/resource_data_wo.yaml"
  data_wo = {
    hello = "%[2]s"
    nested = {
      key = "value"
    }
    port    = 5432
    enabled = true
  }
  data_wo_version = "%[3]s"
}

data "sops_file" "test_data_wo_read" {
  source_file = sops_file.test_data_wo.destination_file
}`

func TestResourceSopsFile_dataWriteOnly(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	filename := fmt.Sprintf("%s/test-fixtures/resource_data_wo.yaml", wd)

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_11_0)},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			if _, err := os.Stat(filename); !os.IsNotExist(err) {
				return fmt.Errorf("expected %s to be removed after destroy", filename)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(configTestResourceSopsFile_dataWriteOnly, wd, "world", "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sops_file.test_data_wo", "id", filename),
					resource.TestCheckNoResourceAttr("sops_file.test_data_wo", "data_wo"),
					resource.TestCheckResourceAttr("sops_file.test_data_wo", "data_wo_version", "1"),
					resource.TestCheckResourceAttr("data.sops_file.test_data_wo_read", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_file.test_data_wo_read", "data.nested.key", "value"),
					resource.TestCheckResourceAttr("data.sops_file.test_data_wo_read", "data.port", "5432"),
					resource.TestCheckResourceAttr("data.sops_file.test_data_wo_read", "data.enabled", "true"),
					func(*terraform.State) error {
						raw, err := os.ReadFile(filename)
						if err != nil {
							return err
						}
						if strings.Contains(string(raw), "hello: world") {
							return fmt.Errorf("file on disk is not encrypted: %s", raw)
						}
						if !strings.Contains(string(raw), "sops:") {
							return fmt.Errorf("file on disk does not look like a sops file: %s", raw)
						}
						if !strings.Contains(string(raw), "type:int") {
							return fmt.Errorf("expected port to be encrypted with type:int, got: %s", raw)
						}
						return nil
					},
				),
			},
			{
				Config: fmt.Sprintf(configTestResourceSopsFile_dataWriteOnly, wd, "updated", "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sops_file.test_data_wo", "data_wo_version", "2"),
					resource.TestCheckResourceAttr("data.sops_file.test_data_wo_read", "data.hello", "updated"),
				),
			},
		},
	})
}
