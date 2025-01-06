package nops

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
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
	ID          types.Int64  `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	ProjectId   types.Int64  `tfsdk:"project_id"`
	Status      types.String `tfsdk:"status"`
	Region      types.String `tfsdk:"region"`
	Bucket      types.String `tfsdk:"bucket"`
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
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Timestamp when the resource was last updated.",
			},
			"project_id": schema.Int64Attribute{
				Required:    true,
				Description: "nOps project ID.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the resource was last updated.",
			},
			"bucket": schema.StringAttribute{
				Computed:    true,
				Description: "AWS bucket name associate with this integration.",
			},
			"region": schema.StringAttribute{
				Computed:    true,
				Description: "AWS region where the bucket resides.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "nOps Container Cost Bucket integration status.",
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

	allContainerBuckets, err := r.client.GetContainerCostBucketSetupStatus()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting remote compute copilot onboarding data",
			err.Error(),
		)
		return
	}

	for _, integration := range *allContainerBuckets {
		if plan.ProjectId.ValueInt64() == integration.Project {
			plan.ID = types.Int64Value(integration.ID)
			plan.Bucket = types.StringValue(integration.Bucket)
			plan.Region = types.StringValue(integration.Region)
			plan.Status = types.StringValue(integration.Status)
		}
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

	containerCostBucketStatus, err := r.client.GetTargetedContainerCostBucketSetupStatus(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting remote compute copilot onboarding data",
			err.Error(),
		)
		return
	}

	state.ID = types.Int64Value(containerCostBucketStatus.ID)
	state.Bucket = types.StringValue(containerCostBucketStatus.Bucket)
	state.Region = types.StringValue(containerCostBucketStatus.Region)
	state.Status = types.StringValue(containerCostBucketStatus.Status)

	tflog.Debug(ctx, "Upstream container cost bucket data received for project "+strconv.Itoa(int(state.ProjectId.ValueInt64())))

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

	allContainerBuckets, err := r.client.GetContainerCostBucketSetupStatus()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting remote compute copilot onboarding data",
			err.Error(),
		)
		return
	}

	for _, integration := range *allContainerBuckets {
		if plan.ProjectId.ValueInt64() == integration.Project {
			plan.ID = types.Int64Value(integration.ID)
			plan.Bucket = types.StringValue(integration.Bucket)
			plan.Region = types.StringValue(integration.Region)
			plan.Status = types.StringValue(integration.Status)

		}
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

	err := r.client.DeleteContainerCostBucket(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting container cost bucket",
			err.Error(),
		)
		return
	}
}

func (r *containerCostBucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Capability to import existing cc onboarding already integrated in the nOps platform into the TF state without recreation.
	val, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing ID for import, please check for a correct project ID", err.Error())
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), val)...)
}
