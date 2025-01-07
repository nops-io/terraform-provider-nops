package nops

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestComputeCopilotIntegrationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "nops_compute_copilot_integration" "test" {
  cluster_arns = ["arn:aws:eks:us-west-2:844856862745:cluster/nOps-dev2", "arn:aws:eks:us-west-2:844856862745:cluster/nOps-uat", "arn:aws:eks:us-west-2:844856862745:cluster/uat-compute-copilot-testing"]
  region_name = "us-west-2"
	version = "1.0.0"
	account_id = 23986
}
`,

				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// resource.TestCheckResourceAttr("nops_project.test", "items.#", "1"),
					// Verify first order item
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "region_name", "us-west-2"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "cluster_arns.0", "arn:aws:eks:us-west-2:844856862745:cluster/nOps-dev2"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "cluster_arns.1", "arn:aws:eks:us-west-2:844856862745:cluster/nOps-uat"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "version", "1.0.0"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "account_id", "23986"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
			resource "nops_compute_copilot_integration" "test" {
  			cluster_arns = ["arn:aws:eks:us-west-2:844856862745:cluster/nOps-dev2", "arn:aws:eks:us-west-2:844856862745:cluster/nOps-uat", "arn:aws:eks:us-west-2:844856862745:cluster/uat-compute-copilot-testing"]
				region_name = "us-west-2"
				version = "1.0.1"
				account_id = 23986
			}
			`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "region_name", "us-west-2"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "cluster_arns.0", "arn:aws:eks:us-west-2:844856862745:cluster/nOps-dev2"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "cluster_arns.1", "arn:aws:eks:us-west-2:844856862745:cluster/nOps-uat"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "version", "1.0.1"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "account_id", "23986"),
				),
			},
		},
	})
}
