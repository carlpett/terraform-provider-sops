package sops

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const configTestDataSourceSopsFile_basic = `
data "sops_file" "test_basic" {
  source_file = "%s/test-fixtures/basic.yaml"
}`

func TestDataSourceSopsFile_basic(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_basic, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.hello", "world"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFile_nested = `
data "sops_file" "test_nested" {
  source_file = "%s/test-fixtures/nested.yaml"
}`

func TestDataSourceSopsFile_nested(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_nested, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_nested", "data.db.user", "foo"),
					resource.TestCheckResourceAttr("data.sops_file.test_nested", "data.db.password", "bar"),
				),
			},
		},
	})
}
