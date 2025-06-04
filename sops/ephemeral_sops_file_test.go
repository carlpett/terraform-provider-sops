package sops

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const configTestEphemeralSopsFile_basic = `
ephemeral "sops_file" "test_basic" {
  source_file = "%s/test-fixtures/basic.yaml"
}

provider "echo" {
  data = ephemeral.sops_file.test_basic.data
}

resource "echo" "test_basic" {}
`

func TestEphemeralSopsFile_basic(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsFile_basic, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test_basic", "data.hello", "world"),
					resource.TestCheckResourceAttr("echo.test_basic", "data.integer", "0"),
					resource.TestCheckResourceAttr("echo.test_basic", "data.float", "0.2"),
					resource.TestCheckResourceAttr("echo.test_basic", "data.bool", "true"),
					resource.TestCheckResourceAttr("echo.test_basic", "data.null_value", "null"),
				),
			},
		},
	})
}

const configTestEphemeralSopsFile_nested = `
ephemeral "sops_file" "test_nested" {
  source_file = "%s/test-fixtures/nested.yaml"
}

provider "echo" {
  data = ephemeral.sops_file.test_nested.data
}

resource "echo" "test_nested" {}
`

func TestEphemeralSopsFile_nested(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsFile_nested, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test_nested", "data.db.user", "foo"),
					resource.TestCheckResourceAttr("echo.test_nested", "data.db.password", "bar"),
				),
			},
		},
	})
}

const configTestEphemeralSopsFile_raw = `
ephemeral "sops_file" "test_raw" {
  source_file = "%s/test-fixtures/raw.txt"
  input_type = "raw"
}

provider "echo" {
  data = ephemeral.sops_file.test_raw.raw
}

resource "echo" "test_raw" {}
`

func TestEphemeralSopsFile_raw(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsFile_raw, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test_raw", "data", "Hello raw world!"),
				),
			},
		},
	})
}

const configTestEphemeralSopsFile_simplelist = `
ephemeral "sops_file" "test_list" {
  source_file = "%s/test-fixtures/simple-list.yaml"
}

provider "echo" {
  data = ephemeral.sops_file.test_list.data
}

resource "echo" "test_list" {}
`

func TestEphemeralSopsFile_simplelist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsFile_simplelist, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test_list", "data.a_list.0", "val1"),
					resource.TestCheckResourceAttr("echo.test_list", "data.a_list.1", "val2"),
					resource.TestCheckResourceAttr("echo.test_list", "data.a_list.2", "null"),
				),
			},
		},
	})
}

const configTestEphemeralSopsFile_complexlist = `
ephemeral "sops_file" "test_list" {
  source_file = "%s/test-fixtures/complex-list.yaml"
}

provider "echo" {
  data = ephemeral.sops_file.test_list.data
}

resource "echo" "test_list" {}
`

func TestEphemeralSopsFile_complexlist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsFile_complexlist, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test_list", "data.a_list.0.name", "foo"),
					resource.TestCheckResourceAttr("echo.test_list", "data.a_list.0.index", "0"),
					resource.TestCheckResourceAttr("echo.test_list", "data.a_list.0.value", "null"),
					resource.TestCheckResourceAttr("echo.test_list", "data.a_list.1.name", "bar"),
					resource.TestCheckResourceAttr("echo.test_list", "data.a_list.1.index", "1"),
				),
			},
		},
	})
}

const configTestEphemeralSopsFile_json = `
ephemeral "sops_file" "test_json" {
  source_file = "%s/test-fixtures/basic.json"
}

provider "echo" {
  data = ephemeral.sops_file.test_json.data
}

resource "echo" "test_json" {}
`

func TestEphemeralSopsFile_json(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsFile_json, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("echo.test_json", "data.hello", "world"),
					resource.TestCheckResourceAttr("echo.test_json", "data.integer", "0"),
					resource.TestCheckResourceAttr("echo.test_json", "data.float", "0.2"),
					resource.TestCheckResourceAttr("echo.test_json", "data.bool", "true"),
					resource.TestCheckResourceAttr("echo.test_json", "data.null", "null"),
				),
			},
		},
	})
}
