package nops

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestComputeCopilotIntegrationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "nops_compute_copilot_integration" "test" {
  cluster_arns = ["arn:aws:eks:us-west-2:844856862745:cluster/nOps-dev2"]
  region_name = "us-west-2"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// resource.TestCheckResourceAttr("nops_project.test", "items.#", "1"),
					// Verify first order item
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "region_name", "us-west-2"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "cluster_arns.0", "arn:aws:eks:us-west-2:844856862745:cluster/nOps-dev2"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "nops_compute_copilot_integration" "test" {
  cluster_arns = ["arn:aws:eks:us-west-2:844856862745:cluster/nOps-uat"]
  region_name = "us-west-2"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "region_name", "us-west-2"),
					resource.TestCheckResourceAttr("nops_compute_copilot_integration.test", "cluster_arns.0", "arn:aws:eks:us-west-2:844856862745:cluster/nOps-uat"),
				),
			},
		},
	})
}
