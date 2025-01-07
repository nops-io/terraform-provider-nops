package nops

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &computeCopilotIntegrationResource{}
	_ resource.ResourceWithConfigure = &computeCopilotIntegrationResource{}
)

// computeCopilotResource is the resource implementation.
type computeCopilotIntegrationResource struct {
	client *Client
}

type computeCopilotIntegrationModel struct {
	LastUpdated types.String `tfsdk:"last_updated"`
	ClusterArns types.List   `tfsdk:"cluster_arns"`
	RegionName  types.String `tfsdk:"region_name"`
	Version     types.String `tfsdk:"version"`
	AccountID   types.String `tfsdk:"account_id"`
}

// computeCopilotResource is a helper function to simplify the provider implementation.
func computeCopilotResource() resource.Resource {
	return &computeCopilotIntegrationResource{}
}

// Configure adds the provider configured client to the resource.
func (r *computeCopilotIntegrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *computeCopilotIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compute_copilot_integration"
}

// Schema defines the schema for the resource.
func (r *computeCopilotIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Notifies the nOps platform a new cluster has been onboarded to nOps with the required input values." +
			" This resource is mostly used only for secure connection with nOps APIs.",
		Attributes: map[string]schema.Attribute{
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the resource was last updated",
			},
			"cluster_arns": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of EKS cluster arns to be onboarded.",
			},
			"region_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the AWS region where the EKS clusters run.",
			},
			"version": schema.StringAttribute{
				Required:    true,
				Description: "Module version being applied.",
			},
			"account_id": schema.StringAttribute{
				Required:    true,
				Description: "nOps account ID associated with the AWS account where the clusters are hosted.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *computeCopilotIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan computeCopilotIntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster_arns := make([]string, 0, len(plan.ClusterArns.Elements()))
	diags = plan.ClusterArns.ElementsAs(ctx, &cluster_arns, false)
	if diags.HasError() {
		return
	}

	// Notify nOps with new values
	var integration ComputeCopilotOnboarding
	integration.ClusterArns = cluster_arns
	integration.RegionName = plan.RegionName.ValueString()
	integration.Version = plan.Version.ValueString()
	integration.AccountID = plan.AccountID.ValueString()
	err := r.client.NotifyComputeCopilotOnboarding(integration)
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
	tflog.Info(ctx, "Created nOps compute copilot integration resource", map[string]any{"Clusters": plan.ClusterArns, "Region": plan.RegionName})

}

// Read refreshes the Terraform state with the latest data.
func (r *computeCopilotIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state computeCopilotIntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	onboarding, err := r.client.GetComputeCopilotOnboarding(state.RegionName.ValueString(), state.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting remote compute copilot onboarding data",
			err.Error(),
		)
		return
	}

	clusterArnsListValue, valueFromDiags := types.ListValueFrom(ctx, types.StringType, onboarding.ClusterArns)
	if valueFromDiags.HasError() {
		resp.Diagnostics.AddError(
			"Error converting cluster_arns to list",
			"Error converting cluster_arns to list",
		)
		return
	}

	clusterArns := make([]string, 0, len(state.ClusterArns.Elements()))
	diags = state.ClusterArns.ElementsAs(ctx, &clusterArns, false)
	if diags.HasError() {
		return
	}

	// Required as the API might not return the list in the same order as its provided in the module, creating a permadiff.
	sort.Strings(onboarding.ClusterArns)
	sort.Strings(clusterArns)

	if !equal(onboarding.ClusterArns, clusterArns) {
		tflog.Debug(ctx, fmt.Sprintf("%s != %s", strings.Join(onboarding.ClusterArns, ","), strings.Join(clusterArns, ",")))
		state.ClusterArns = clusterArnsListValue
	}
	tflog.Debug(ctx, "Upstream compute copilot integration project data received for clusters "+strings.Join(onboarding.ClusterArns, ",")+" region: "+onboarding.RegionName)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *computeCopilotIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan computeCopilotIntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterArns := make([]string, 0, len(plan.ClusterArns.Elements()))
	diags = plan.ClusterArns.ElementsAs(ctx, &clusterArns, false)
	if diags.HasError() {
		return
	}

	// Notify nOps with new values
	var integration ComputeCopilotOnboarding
	integration.ClusterArns = clusterArns
	integration.RegionName = plan.RegionName.ValueString()
	integration.Version = plan.Version.ValueString()
	integration.AccountID = plan.AccountID.ValueString()
	err := r.client.NotifyComputeCopilotOnboarding(integration)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error notifying nOps",
			"Failed to notify, unexpected error: "+err.Error(),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Created nOps compute copilot integration resource", map[string]any{"Clusters": plan.ClusterArns, "Region": plan.RegionName})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *computeCopilotIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Framework automatically removes resource from state, no action to be taken on that side.
	var state computeCopilotIntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteComputeCopilotOnboarding(state.RegionName.ValueString(), state.AccountID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project",
			err.Error(),
		)
		return
	}
}

// func (r *computeCopilotIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Capability to import existing cc onboarding already integrated in the nOps platform into the TF state without recreation.
// 	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region_name"), req.ID)...)
// }
