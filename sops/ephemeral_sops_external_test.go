package sops

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const configTestEphemeralSopsExternal_basic = `
ephemeral "sops_external" "test_basic" {
  source     = file("%s/test-fixtures/basic.yaml")
  input_type = "yaml"
}

provider "echo" {
  data = ephemeral.sops_external.test_basic.data
}

resource "echo" "test_basic" {}
`

func TestEphemeralSopsExternal(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsExternal_basic, wd)
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

const configTestEphemeralSopsExternal_nested = `
ephemeral "sops_external" "test_nested" {
  source     = file("%s/test-fixtures/nested.yaml")
  input_type = "yaml"
}

provider "echo" {
  data = ephemeral.sops_external.test_nested.data
}

resource "echo" "test_nested" {}
`

func TestEphemeralSopsExternal_nested(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsExternal_nested, wd)
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

const configTestEphemeralSopsExternal_raw = `
ephemeral "sops_external" "test_raw" {
  source     = file("%s/test-fixtures/raw.txt")
  input_type = "raw"
}

provider "echo" {
  data = ephemeral.sops_external.test_raw.raw
}

resource "echo" "test_raw" {}
`

func TestEphemeralSopsExternal_raw(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsExternal_raw, wd)
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

const configTestEphemeralSopsExternal_simplelist = `
ephemeral "sops_external" "test_list" {
  source     = file("%s/test-fixtures/simple-list.yaml")
  input_type = "yaml"
}

provider "echo" {
  data = ephemeral.sops_external.test_list.data
}

resource "echo" "test_list" {}
`

func TestEphemeralSopsExternal_simplelist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsExternal_simplelist, wd)
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

const configTestEphemeralSopsExternal_complexlist = `
ephemeral "sops_external" "test_list" {
  source     = file("%s/test-fixtures/complex-list.yaml")
  input_type = "yaml"
}

provider "echo" {
  data = ephemeral.sops_external.test_list.data
}

resource "echo" "test_list" {}
`

func TestEphemeralSopsExternal_complexlist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsExternal_complexlist, wd)
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

const configTestEphemeralSopsExternal_json = `
ephemeral "sops_external" "test_json" {
  source     = file("%s/test-fixtures/basic.json")
  input_type = "json"
}

provider "echo" {
  data = ephemeral.sops_external.test_json.data
}

resource "echo" "test_json" {}
`

func TestEphemeralSopsExternal_json(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestEphemeralSopsExternal_json, wd)
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
