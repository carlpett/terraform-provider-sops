package sops

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const configTestDataSourceSopsExternal_basic = `
data "sops_external" "test_basic" {
  source     = file("%s/test-fixtures/basic.yaml")
  input_type = "yaml"
}`

func TestDataSourceSopsExternal(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsExternal_basic, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.bool", "true"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_nested = `
data "sops_external" "test_nested" {
  source     = file("%s/test-fixtures/nested.yaml")
  input_type = "yaml"
}`

func TestDataSourceSopsExternal_nested(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsExternal_nested, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_nested", "data.db.user", "foo"),
					resource.TestCheckResourceAttr("data.sops_external.test_nested", "data.db.password", "bar"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_raw = `
data "sops_external" "test_raw" {
  source     = file("%s/test-fixtures/raw.txt")
  input_type = "raw"
}`

func TestDataSourceSopsExternal_raw(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsExternal_raw, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_raw", "raw", "Hello raw world!"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_simplelist = `
data "sops_external" "test_list" {
  source     = file("%s/test-fixtures/simple-list.yaml")
  input_type = "yaml"
}`

func TestDataSourceSopsExternal_simplelist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsExternal_simplelist, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0", "val1"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.1", "val2"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_complexlist = `
data "sops_external" "test_list" {
  source     = file("%s/test-fixtures/complex-list.yaml")
  input_type = "yaml"
}`

func TestDataSourceSopsExternal_complexlist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsExternal_complexlist, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0.name", "foo"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0.index", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.1.name", "bar"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.1.index", "1"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_json = `
data "sops_external" "test_json" {
  source     = file("%s/test-fixtures/basic.json")
  input_type = "json"
}`

func TestDataSourceSopsExternal_json(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsExternal_json, wd)
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.bool", "true"),
				),
			},
		},
	})
}
