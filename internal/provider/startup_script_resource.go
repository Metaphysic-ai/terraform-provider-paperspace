package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"terraform-provider-paperspace/internal/psclient"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &startupScriptResource{}
	_ resource.ResourceWithConfigure = &startupScriptResource{}
)

// NewStartupScriptResource is a helper function to simplify the provider implementation.
func NewStartupScriptResource() resource.Resource {
	return &startupScriptResource{}
}

// Maps the resource schema data.
// State/Plan structure.
type startupScriptResourceModel struct {
	Name      types.String `tfsdk:"name"`   // required
	Script    types.String `tfsdk:"script"` // required
	IsRunOnce types.Bool   `tfsdk:"is_run_once"`

	// Computed only
	ID                 types.String `tfsdk:"id"`
	Description        types.String `tfsdk:"description"`
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	AssignedMachineIDs types.List   `tfsdk:"assigned_machine_ids"`
	DtCreated          types.String `tfsdk:"dt_created"`
}

// startupScriptResource is the resource implementation.
type startupScriptResource struct {
	// Allow resource to store a reference to the client
	client *psclient.Client
}

// Define the resource type name, which is how the resource is used in Terraform configurations.
func (r *startupScriptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_startup_script"
}

// Define Schema
// The resource uses the Schema method to define the supported configuration, plan, and state attribute names and types.
func (r *startupScriptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	// elementType: types.StringType, elements: []attr.Value{}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Startup script resource",
		Attributes: map[string]schema.Attribute{
			// TODO: Implement update and remove RequiresReplace plan modifiers

			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the startup script.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// TODO: Consider sensitive or encoded
			"script": schema.StringAttribute{
				MarkdownDescription: "The script to run on startup.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_run_once": schema.BoolAttribute{
				MarkdownDescription: "Whether the script should only run once on first boot or on every boot.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},

			// Computed only
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the startup script.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the startup script.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"is_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the startup script is enabled.",
				Computed:            true, // Cannot be set during creation, so computed only to simplify
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"assigned_machine_ids": schema.ListAttribute{
				MarkdownDescription: "The IDs of the machines the startup script is assigned to.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"dt_created": schema.StringAttribute{
				MarkdownDescription: "The date the startup script was created.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

// TODO: Add timeouts https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts
// Create a new resource.
func (r *startupScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan startupScriptResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan and create new startup script

	reqData := psclient.StartupScriptCreateConfig{
		Name:      plan.Name.ValueString(),   // required
		Script:    plan.Script.ValueString(), // required
		IsRunOnce: plan.IsRunOnce.ValueBool(),
	}

	// TODO: Save script into state in base64

	jsonData, _ := json.MarshalIndent(reqData, "", " ")
	tflog.Info(ctx, "Sending create req data: "+string(jsonData))

	startupScript, err := r.client.CreateStartupScript(reqData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating startup script",
			"Could not create startup script, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, "Created a startup script resource with id "+startupScript.ID)

	// Map response body to schema and populate Computed attribute values.
	// Save response data into the Terraform state.
	// Only computed attributes must be updated here

	plan.ID = types.StringValue(startupScript.ID)
	// Handling List values
	assignedMachineIDs, diags := types.ListValueFrom(ctx, types.StringType, startupScript.AssignedMachineIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.AssignedMachineIDs = assignedMachineIDs
	fillStateWithStartupScriptData(&plan, startupScript)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *startupScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state startupScriptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed data from Paperspace
	startupScript, err := r.client.GetStartupScript(state.ID.ValueString())
	if err != nil {
		// Avoid error due to 404 status. Resource could be deleted outside provider, so handle this as expected case.
		if strings.Contains(err.Error(), "status: 404") {
			tflog.Warn(ctx, fmt.Sprintf("Startup script %s not found, removing from state", state.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Paperspace startup script",
			"Could not read Paperspace startup script ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// TODO: Remove log, or remove indent
	startupScriptJson, err := json.MarshalIndent(startupScript, "", " ")
	if err != nil {
		tflog.Error(ctx, "Could not marshal fetched startup script struct: "+err.Error())
	}
	tflog.Info(ctx, "Fetched startup script data: "+string(startupScriptJson))

	// ID not needed here
	// TODO: Consider moving list handling to function
	// Handling List values
	assignedMachineIDs, diags := types.ListValueFrom(ctx, types.StringType, startupScript.AssignedMachineIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.AssignedMachineIDs = assignedMachineIDs
	fillStateWithStartupScriptData(&state, startupScript)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Updates the resource and sets the updated Terraform state on success.
func (r *startupScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *startupScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state startupScriptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	startupScriptID := state.ID.ValueString()
	err := r.client.DeleteStartupScript(startupScriptID)
	if err != nil {
		if strings.Contains(err.Error(), "status: 404") {
			tflog.Info(ctx, fmt.Sprintf("Startup script %s not found, assuming already deleted", startupScriptID))
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting startup script",
			"Could not delete startup script, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *startupScriptResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*psclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *psclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func fillStateWithStartupScriptData(state *startupScriptResourceModel, startupScript *psclient.StartupScript) {
	state.Name = types.StringValue(startupScript.Name)
	state.IsRunOnce = types.BoolValue(startupScript.IsRunOnce)
	state.Description = types.StringPointerValue(startupScript.Description)
	state.IsEnabled = types.BoolValue(startupScript.IsEnabled)
	state.DtCreated = types.StringValue(startupScript.DtCreated)
}
