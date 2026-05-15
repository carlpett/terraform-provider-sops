package sops

import (
	"context"
	"fmt"
	"os"

	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func testGenerateConfigHashResourceInputWO(value string) string {
	return fmt.Sprintf(`
resource "sops_hash" "test" {
  input_wo = %[1]q
}
`, value)
}

func TestAccHashResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testGenerateConfigHashResourceInputWO("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"sops_hash.test",
						tfjsonpath.New("hash"),
						knownvalue.StringExact(hashInput("one")),
					),
				},
			},
			// Update and Read testing
			{
				Config: testGenerateConfigHashResourceInputWO("two"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"sops_hash.test",
						tfjsonpath.New("hash"),
						knownvalue.StringExact(hashInput("two")),
					),
				},
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectNonEmptyPlan()),
			},
			{
				Config: testGenerateConfigHashResourceInputWO("two"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectEmptyPlan()),
			},
			{
				Config: testGenerateConfigHashResourceInputWO("one"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectNonEmptyPlan()),
			},
		},
	})
}

const configTestResourceHashEphemeralVar = `
variable "ephem" {
  type = string
  ephemeral = true
}

resource "sops_hash" "testephem" {
  input_wo = var.ephem
}
`

func TestAccHasherWithEphemeral(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configTestResourceHashEphemeralVar,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"sops_hash.testephem",
						tfjsonpath.New("hash"),
						knownvalue.StringExact(hashInput("testvalue")),
					),
				},
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectNonEmptyPlan()),
				ConfigVariables: map[string]config.Variable{
					"ephem": config.StringVariable("testvalue"),
				},
			},
			{
				Config: configTestResourceHashEphemeralVar,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"sops_hash.testephem",
						tfjsonpath.New("hash"),
						knownvalue.StringExact(hashInput("testvalue")),
					),
				},
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectEmptyPlan()),
				ConfigVariables: map[string]config.Variable{
					"ephem": config.StringVariable("testvalue"),
				},
			},
			{
				Config: configTestResourceHashEphemeralVar,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"sops_hash.testephem",
						tfjsonpath.New("hash"),
						knownvalue.StringExact(hashInput("newvalue")),
					),
				},
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectNonEmptyPlan()),
				ConfigVariables: map[string]config.Variable{
					"ephem": config.StringVariable("newvalue"),
				},
			},
		},
	})
}

// testExpectPreApply returns a set of plan checks that includes the provided
// pre-apply check combined with empty plan checks for both post apply checks.
func testExpectPreApply(p plancheck.PlanCheck) resource.ConfigPlanChecks {
	return resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			p,
		},
		PostApplyPreRefresh: []plancheck.PlanCheck{
			plancheck.ExpectEmptyPlan(),
		},
		PostApplyPostRefresh: []plancheck.PlanCheck{
			plancheck.ExpectEmptyPlan(),
		},
	}
}

// testHashResourcePlanChecker implements PlanCheck and allows custom check functions
// for specific resources in the plan.
type testHashResourcePlanChecker struct {
	randomIntegerChecker func(*tfjson.ResourceChange) error
	sopsHashChecker      func(*tfjson.ResourceChange) error
	otherChecker         func(*tfjson.ResourceChange) error
}

func (c testHashResourcePlanChecker) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	errors := []error{}

	for _, r := range req.Plan.ResourceChanges {
		if r.Type == "random_integer" && c.randomIntegerChecker != nil {
			err := c.randomIntegerChecker(r)
			if err != nil {
				errors = append(errors, fmt.Errorf("random_integer check failed: %w", err))
			}
			continue
		}
		if r.Address == "sops_hash.test" && c.sopsHashChecker != nil {
			err := c.sopsHashChecker(r)
			if err != nil {
				errors = append(errors, fmt.Errorf("sops_hash check failed: %w", err))
			}
			continue
		}
		if c.otherChecker != nil {
			err := c.otherChecker(r)
			if err != nil {
				errors = append(errors, fmt.Errorf("other check failed: %w", err))
			}
		}
	}
	if len(errors) > 0 {
		resp.Error = fmt.Errorf("plan check failed with errors: %v", errors)
	}
}

func testPlanChecker_ExpectRandomChangeAndSopsHashUpdate() testHashResourcePlanChecker {
	return testHashResourcePlanChecker{
		randomIntegerChecker: func(r *tfjson.ResourceChange) error {
			if !r.Change.Actions.Replace() {
				return fmt.Errorf("expected random_integer resource to have Replace action, got: %v", r.Change.Actions)
			}
			return nil
		},
		sopsHashChecker: func(r *tfjson.ResourceChange) error {
			if !r.Change.Actions.Update() {
				return fmt.Errorf("expected sops_hash resource to have Update action, got: %v", r.Change.Actions)
			}
			return nil
		},
		otherChecker: func(r *tfjson.ResourceChange) error {
			return fmt.Errorf("unexpected resource change for resource %s of type %s with actions: %v", r.Address, r.Type, r.Change.Actions)
		},
	}
}

func testPlanChecker_ExpectSopsHashUpdateOnly() testHashResourcePlanChecker {
	return testHashResourcePlanChecker{
		randomIntegerChecker: func(r *tfjson.ResourceChange) error {
			if !r.Change.Actions.NoOp() {
				return fmt.Errorf("unexpected change to random_integer resource %s with actions: %v", r.Address, r.Change.Actions)
			}
			return nil
		},
		sopsHashChecker: func(r *tfjson.ResourceChange) error {
			if !r.Change.Actions.Update() {
				return fmt.Errorf("expected sops_hash resource to have Update action, got: %v", r.Change.Actions)
			}
			return nil
		},
		otherChecker: func(r *tfjson.ResourceChange) error {
			return fmt.Errorf("unexpected resource change for resource %s of type %s with actions: %v", r.Address, r.Type, r.Change.Actions)
		},
	}
}

func testGenerateTransitionConfig(workdir, input, keepers string) string {
	if keepers == "" {
		keepers = "same"
	}
	return fmt.Sprintf(`
resource "random_integer" "fixed_val_zero" { 
	min = 0
	max = 0
	keepers = { always = "%s" }
}
resource "random_integer" "fixed_val_ten" { 
	min = 10
	max = 10
	keepers = { always = "%s" }
}
ephemeral "sops_file" "test_basic" { source_file = "%s/test-fixtures/basic.yaml" }
resource "sops_hash" "test" { input_wo = %s }
output "test" { value = sops_hash.test.hash }
`, keepers, keepers, workdir, input)
}

func TestResourceHashUnknownTransitions(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "registry.terraform.io/hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testGenerateTransitionConfig(wd, "\"0\"", "keeper"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectNonEmptyPlan()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sops_hash.test", "hash", hashInput("0")),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact(hashInput("0")),
					),
				},
			},
			{
				// Changing to a random value should not trigger a plan since the random value is already known (keeper unchanged)
				Config: testGenerateTransitionConfig(wd, "random_integer.fixed_val_zero.result", "keeper"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectEmptyPlan()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sops_hash.test", "hash", hashInput("0")),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact(hashInput("0")),
					),
				},
			},
			{
				// Changing the random keeper should trigger a plan with sops_hash update since the random value is unknown.
				Config: testGenerateTransitionConfig(wd, "random_integer.fixed_val_zero.result", "keeper1"),
				ConfigPlanChecks: testExpectPreApply(testPlanChecker_ExpectRandomChangeAndSopsHashUpdate()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sops_hash.test", "hash", hashInput("0")),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact(hashInput("0")),
					),
				},
			},
			{
				// Same random keeper means empty plan on next apply.
				Config: testGenerateTransitionConfig(wd, "ephemeral.sops_file.test_basic.data.integer", "keeper1"),
				ConfigPlanChecks: testExpectPreApply(plancheck.ExpectEmptyPlan()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sops_hash.test", "hash", hashInput("0")),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact(hashInput("0")),
					),
				},
			},
			{
				// Changing the write only input to a new known value should trigger a plan with sops_hash update.
				Config: testGenerateTransitionConfig(wd, "random_integer.fixed_val_ten.result", "keeper1"),
				ConfigPlanChecks: testExpectPreApply(testPlanChecker_ExpectSopsHashUpdateOnly()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("random_integer.fixed_val_ten", "result", "10"),
					resource.TestCheckResourceAttr("sops_hash.test", "hash", hashInput("10")),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"test",
						knownvalue.StringExact(hashInput("10")),
					),
				},
			},
		},
	})
}
