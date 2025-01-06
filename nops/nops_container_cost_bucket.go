package nops

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &containerCostBucketResource{}
	_ resource.ResourceWithConfigure = &containerCostBucketResource{}
)

// computeCopilotResource is the resource implementation.
type containerCostBucketResource struct {
	client *Client
}

type containerCostBucketModel struct {
	LastUpdated types.String `tfsdk:"last_updated"`
	ProjectId   types.Int64  `tfsdk:"project_id"`
}

// computeCopilotResource is a helper function to simplify the provider implementation.
func containerCostResource() resource.Resource {
	return &containerCostBucketResource{}
}

// Configure adds the provider configured client to the resource.
func (r *containerCostBucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *containerCostBucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_cost_bucket"
}

// Schema defines the schema for the resource.
func (r *containerCostBucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Notifies the nOps platform a new container cost bucket was created for the backend to fetch metadata from it." +
			" This resource is mostly used only for secure connection with nOps APIs.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Required:    true,
				Description: "nOps project ID.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the resource was last updated.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *containerCostBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan containerCostBucketModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var containerCostBucket ContainerCostBucketSetup
	containerCostBucket.Project = plan.ProjectId.ValueInt64()

	err := r.client.NotifyContainerCostBucketSetup(containerCostBucket)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error notifying nOps",
			"Failed to notify, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Created nOps container cost bucket resource", map[string]any{"project_id": plan.ProjectId})

}

// Read refreshes the Terraform state with the latest data.
func (r *containerCostBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state containerCostBucketModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *containerCostBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan containerCostBucketModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Created nOps container cost bucket resource", map[string]any{"project_id": plan.ProjectId})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *containerCostBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Framework automatically removes resource from state, no action to be taken on that side.
	var state containerCostBucketModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
