package sops

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.bool", "true"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.null_value", "null"),
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0", "val1"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.1", "val2"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.2", "null"),
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0.name", "foo"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0.index", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0.value", "null"),
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.bool", "true"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.null", "null"),
				),
			},
		},
	})
}

// age1r6eadvpf2cq967pcgc35ahfkkphqs6ln9frmt6nma52z6tq7zddsszzwrg
const configTestDataSourceSopsExternal_envAge0 = `
data "sops_external" "test_env_age0" {
  source     = file("%s/test-fixtures/basic.age0.yaml")
  input_type = "yaml"
  env = {
    SOPS_AGE_KEY = "AGE-SECRET-KEY-14CSLELDXVWL48V82MSDF8EPVN9KWKKCQ38ZQ80V7Q8ZK7EW26WGQ7W29YP"
  }
}`

func TestDataSourceSopsExternal_envAge0(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsExternal_envAge0, wd)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_env_age0", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_external.test_env_age0", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_env_age0", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_external.test_env_age0", "data.bool", "true"),
					resource.TestCheckResourceAttr("data.sops_external.test_env_age0", "data.null_value", "null"),

					func(s *terraform.State) error {
						if v := os.Getenv("SOPS_AGE_KEY"); v != "" {
							return fmt.Errorf("expected SOPS_AGE_KEY unset after data read, got %q", v)
						}
						return nil
					},
				),
			},
		},
	})
}

// age120rd0a9p49wtru227l4889q6qzcxnlcwrmdgukk94453uzmyqexshcmmzc
const configTestDataSourceSopsExternal_envAge1 = `
data "sops_external" "test_env_age1" {
  source     = file("%s/test-fixtures/basic.age1.yaml")
  input_type = "yaml"
  env = {
    SOPS_AGE_KEY = "AGE-SECRET-KEY-1NERY6H8EQDUTP2WQML30ME8NX89DATTUJC0SJ5RMDZGDL735EJ0SJ8JGRT"
  }
}`

func TestDataSourceSopsExternal_envAge1(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsExternal_envAge1, wd)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_env_age1", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_external.test_env_age1", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_env_age1", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_external.test_env_age1", "data.bool", "true"),
					resource.TestCheckResourceAttr("data.sops_external.test_env_age1", "data.null_value", "null"),

					func(s *terraform.State) error {
						if v := os.Getenv("SOPS_AGE_KEY"); v != "" {
							return fmt.Errorf("expected SOPS_AGE_KEY unset after data read, got %q", v)
						}
						return nil
					},
				),
			},
		},
	})
}
