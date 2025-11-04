package sops

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.bool", "true"),
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "data.null_value", "null"),
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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

const configTestDataSourceSopsFile_raw = `
data "sops_file" "test_raw" {
  source_file = "%s/test-fixtures/raw.txt"
  input_type = "raw"
}`

func TestDataSourceSopsFile_raw(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_raw, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_raw", "raw", "Hello raw world!"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFile_simplelist = `
data "sops_file" "test_list" {
  source_file = "%s/test-fixtures/simple-list.yaml"
}`

func TestDataSourceSopsFile_simplelist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_simplelist, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.0", "val1"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.1", "val2"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.2", "null"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFile_complexlist = `
data "sops_file" "test_list" {
  source_file = "%s/test-fixtures/complex-list.yaml"
}`

func TestDataSourceSopsFile_complexlist(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_complexlist, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.0.name", "foo"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.0.index", "0"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.0.value", "null"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.1.name", "bar"),
					resource.TestCheckResourceAttr("data.sops_file.test_list", "data.a_list.1.index", "1"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFile_json = `
data "sops_file" "test_json" {
  source_file = "%s/test-fixtures/basic.json"
}`

func TestDataSourceSopsFile_json(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_json, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_json", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_file.test_json", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_file.test_json", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_file.test_json", "data.bool", "true"),
					resource.TestCheckResourceAttr("data.sops_file.test_json", "data.null", "null"),
				),
			},
		},
	})
}

// age1r6eadvpf2cq967pcgc35ahfkkphqs6ln9frmt6nma52z6tq7zddsszzwrg
const configTestDataSourceSopsFile_envAge0 = `
data "sops_file" "test_env_age0" {
  source_file = "%s/test-fixtures/basic.age0.yaml"
  env = {
    SOPS_AGE_KEY = "AGE-SECRET-KEY-14CSLELDXVWL48V82MSDF8EPVN9KWKKCQ38ZQ80V7Q8ZK7EW26WGQ7W29YP"
  }
}`

func TestDataSourceSopsFile_envAge0(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_envAge0, wd)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_env_age0", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_file.test_env_age0", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_file.test_env_age0", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_file.test_env_age0", "data.bool", "true"),
					resource.TestCheckResourceAttr("data.sops_file.test_env_age0", "data.null_value", "null"),

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
const configTestDataSourceSopsFile_envAge1 = `
data "sops_file" "test_env_age1" {
  source_file = "%s/test-fixtures/basic.age1.yaml"
  env = {
    SOPS_AGE_KEY = "AGE-SECRET-KEY-1NERY6H8EQDUTP2WQML30ME8NX89DATTUJC0SJ5RMDZGDL735EJ0SJ8JGRT"
  }
}`

func TestDataSourceSopsFile_envAge1(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFile_envAge1, wd)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file.test_env_age1", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_file.test_env_age1", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_file.test_env_age1", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_file.test_env_age1", "data.bool", "true"),
					resource.TestCheckResourceAttr("data.sops_file.test_env_age1", "data.null_value", "null"),

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
