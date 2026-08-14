package sops

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "last_modified", "2019-04-26T18:43:59Z"),
					resource.TestCheckResourceAttr("data.sops_file.test_basic", "last_modified_unix", "1556304239"),
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
					resource.TestCheckResourceAttr("data.sops_file.test_raw", "last_modified", "2019-01-23T13:39:15Z"),
					resource.TestCheckResourceAttr("data.sops_file.test_raw", "last_modified_unix", "1548250755"),
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

// testExpectPreApply asserts the supplied pre-apply plan expectation and that
// the configuration then fully converges (no remaining diff) on the post-apply
// refresh passes. Adapted from the plan-convergence helper in PR #167.
func testExpectPreApply(p plancheck.PlanCheck) resource.ConfigPlanChecks {
	return resource.ConfigPlanChecks{
		PreApply:             []plancheck.PlanCheck{p},
		PostApplyPreRefresh:  []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
		PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
	}
}

const configTestDataSourceSopsFile_trigger = `
data "sops_file" "trigger_src" {
  source_file = "%s/test-fixtures/%s"
}

resource "terraform_data" "trigger" {
  input = data.sops_file.trigger_src.last_modified
}`

// TestDataSourceSopsFile_lastModifiedTrigger verifies that last_modified behaves
// as a stable version trigger for a downstream (state-persisted) resource: it
// must not churn when the encrypted file is unchanged, and must change when the
// file changes. This is the property that makes it usable as a `wo_version`.
func TestDataSourceSopsFile_lastModifiedTrigger(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// last_modified flows into a managed resource and converges.
			{
				Config:           fmt.Sprintf(configTestDataSourceSopsFile_trigger, wd, "basic.yaml"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectNonEmptyPlan()),
				Check:            resource.TestCheckResourceAttr("terraform_data.trigger", "output", "2019-04-26T18:43:59Z"),
			},
			// Re-reading the same file must not produce a diff (no spurious version bumps).
			{
				Config:           fmt.Sprintf(configTestDataSourceSopsFile_trigger, wd, "basic.yaml"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectEmptyPlan()),
			},
			// A file with a different lastmodified must change the trigger value.
			{
				Config:           fmt.Sprintf(configTestDataSourceSopsFile_trigger, wd, "nested.yaml"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectNonEmptyPlan()),
				Check:            resource.TestCheckResourceAttr("terraform_data.trigger", "output", "2019-01-23T12:37:02Z"),
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
					resource.TestCheckResourceAttr("data.sops_file.test_json", "last_modified", "2019-08-07T23:18:36Z"),
					resource.TestCheckResourceAttr("data.sops_file.test_json", "last_modified_unix", "1565219916"),
				),
			},
		},
	})
}
