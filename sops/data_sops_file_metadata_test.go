package sops

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const configTestDataSourceSopsFileMetadata_basic = `
data "sops_file_metadata" "test" {
  source_file = "%s/test-fixtures/basic.yaml"
}`

func TestDataSourceSopsFileMetadata_basic(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFileMetadata_basic, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file_metadata.test", "last_modified", "2019-04-26T18:43:59Z"),
					resource.TestCheckResourceAttr("data.sops_file_metadata.test", "last_modified_unix", "1556304239"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFileMetadata_raw = `
data "sops_file_metadata" "test" {
  source_file = "%s/test-fixtures/raw.txt"
  input_type  = "raw"
}`

func TestDataSourceSopsFileMetadata_raw(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config := fmt.Sprintf(configTestDataSourceSopsFileMetadata_raw, wd)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_file_metadata.test", "last_modified", "2019-01-23T13:39:15Z"),
					resource.TestCheckResourceAttr("data.sops_file_metadata.test", "last_modified_unix", "1548250755"),
				),
			},
		},
	})
}

const configTestDataSourceSopsFileMetadata_trigger = `
data "sops_file_metadata" "trigger_src" {
  source_file = "%s/test-fixtures/%s"
}

resource "terraform_data" "trigger" {
  input = data.sops_file_metadata.trigger_src.last_modified_unix
}`

// TestDataSourceSopsFileMetadata_lastModifiedTrigger demonstrates the intended
// use of the metadata data source: driving a state-persisted version (here a
// terraform_data input, standing in for a write-only `wo_version`) from the
// integer last_modified_unix, without decrypting the file. The trigger must be
// stable when the file is unchanged and change when the file changes.
func TestDataSourceSopsFileMetadata_lastModifiedTrigger(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:           fmt.Sprintf(configTestDataSourceSopsFileMetadata_trigger, wd, "basic.yaml"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectNonEmptyPlan()),
				Check:            resource.TestCheckResourceAttr("terraform_data.trigger", "output", "1556304239"),
			},
			{
				Config:           fmt.Sprintf(configTestDataSourceSopsFileMetadata_trigger, wd, "basic.yaml"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectEmptyPlan()),
			},
			{
				Config:           fmt.Sprintf(configTestDataSourceSopsFileMetadata_trigger, wd, "nested.yaml"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectNonEmptyPlan()),
				Check:            resource.TestCheckResourceAttr("terraform_data.trigger", "output", "1548247022"),
			},
		},
	})
}
