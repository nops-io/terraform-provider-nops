package nops

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestContainerCostBucketResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "nops_container_cost_bucket" "test" {
  project_id = 23986
}
`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nops_container_cost_bucket.test", "bucket", "nops-container-cost-844856862745"),
					resource.TestCheckResourceAttr("nops_container_cost_bucket.test", "region", "us-east-1"),
					resource.TestCheckResourceAttr("nops_container_cost_bucket.test", "project_id", "23986"),
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
			resource "nops_container_cost_bucket" "test" {
				project_id = 23986
			}
			`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nops_container_cost_bucket.test", "bucket", "nops-container-cost-844856862745"),
					resource.TestCheckResourceAttr("nops_container_cost_bucket.test", "region", "us-east-1"),
					resource.TestCheckResourceAttr("nops_container_cost_bucket.test", "project_id", "23986"),
				),
			},
		},
	})
}
