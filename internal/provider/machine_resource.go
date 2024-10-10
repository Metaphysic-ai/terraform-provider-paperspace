package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"terraform-provider-paperspace/internal/ppclient"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &machineResource{}
	_ resource.ResourceWithConfigure = &machineResource{}
)

// NewMachineResource is a helper function to simplify the provider implementation.
func NewMachineResource() resource.Resource {
	return &machineResource{}
}

// machineResourceModel maps the resource schema data.
// State/Plan structure
type machineResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`         // Required
	MachineType types.String `tfsdk:"machine_type"` // Required
	TemplateID  types.String `tfsdk:"template_id"`  // Required
	DiskSize    types.Int32  `tfsdk:"disk_size"`    // Required
	Region      types.String `tfsdk:"region"`       // Required

	NetworkID types.String `tfsdk:"network_id"`

	// computed only
	CPUs      types.Int32  `tfsdk:"cpus"`
	State     types.String `tfsdk:"state"`
	OS        types.String `tfsdk:"os"`
	AgentType types.String `tfsdk:"agent_type"`
	PublicIP  types.String `tfsdk:"public_ip"`
	PrivateIP types.String `tfsdk:"private_ip"`

	// TODO: Implement remaining
	// // MachineConfig
	// AutoSnapshotEnabled   types.Bool   `tfsdk:"auto_snapshot_enabled"`
	// AutoSnapshotFrequency types.String `tfsdk:"auto_snapshot_frequency"`
	// AutoSnapshotSaveCount types.Int64  `tfsdk:"auto_snapshot_save_count"`
	// AutoShutdownEnabled   types.Bool   `tfsdk:"auto_shutdown_enabled"`
	// AutoShutdownTimeout   types.Int64  `tfsdk:"auto_shutdown_timeout"`
	// AutoShutdownForce     types.Bool   `tfsdk:"auto_shutdown_force"`
	// TakeInitialSnapshot   types.Bool   `tfsdk:"take_initial_snapshot"`
	// RestorePointEnabled   types.Bool   `tfsdk:"restore_point_enabled"`
	// RestorePointFrequency types.String `tfsdk:"restore_point_frequency"`
	PublicIpType types.String `tfsdk:"public_ip_type"`
	// StartupScriptID types.String `tfsdk:"startup_script_id"`
	// EmailPassword   types.Bool   `tfsdk:"email_password"`
	StartOnCreate types.Bool `tfsdk:"start_on_create"`
	// EnableNvlink    types.Bool   `tfsdk:"enable_nvlink"`
	AccessorIds types.List `tfsdk:"accessor_ids"`
}

// machineResource is the resource implementation.
type machineResource struct {
	// Allow resource to store a reference to the client
	client *ppclient.Client
}

// Define the resource type name, which is how the resource is used in Terraform configurations.
func (r *machineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_machine"
}

// Define Schema
// The resource uses the Schema method to define the supported configuration, plan, and state attribute names and types.
// The machine resource will need to save a machine with various attributes to Terraform's state.
func (r *machineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	// elementType: types.StringType, elements: []attr.Value{}

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Machine resource",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":         schema.StringAttribute{Required: true},
			"machine_type": schema.StringAttribute{Required: true},
			"template_id":  schema.StringAttribute{Required: true},
			"disk_size":    schema.Int32Attribute{Required: true},
			"region":       schema.StringAttribute{Required: true},
			"network_id":   schema.StringAttribute{Optional: true, Computed: true},
			"cpus":         schema.Int32Attribute{Computed: true, PlanModifiers: []planmodifier.Int32{int32planmodifier.UseStateForUnknown()}},
			"state":        schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"os":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"agent_type":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"public_ip":    schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"private_ip":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},

			"public_ip_type":  schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("dynamic")}, // TODO maybe do not set here, but add some condition later ?
			"start_on_create": schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(false)},

			// "accessor_ids": schema.ListAttribute{Optional: true, ElementType: types.StringType}, // TODO make default []
			"accessor_ids": schema.ListAttribute{Required: true, ElementType: types.StringType}, // TODO make default []
			// // MachineConfig
			// "auto_snapshot_enabled":    schema.BoolAttribute{Computed: true},
			// "auto_snapshot_frequency":  schema.StringAttribute{Computed: true},
			// "auto_snapshot_save_count": schema.Int64Attribute{Computed: true},
			// "auto_shutdown_enabled":    schema.BoolAttribute{Computed: true},
			// "auto_shutdown_timeout":    schema.Int64Attribute{Computed: true},
			// "auto_shutdown_force":      schema.BoolAttribute{Computed: true},
			// "take_initial_snapshot":    schema.BoolAttribute{Computed: true},
			// "restore_point_enabled":    schema.BoolAttribute{Computed: true},
			// // "restore_point_frequency":  schema.StringAttribute{Computed: true},

			// "startup_script_id": schema.StringAttribute{Computed: true},
			// "email_password":    schema.BoolAttribute{Computed: true},

			// "enable_nvlink":     schema.BoolAttribute{Computed: true},
		},
	}
}

// Create a new resource.
func (r *machineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan machineResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan and create new machine

	var reqData ppclient.MachineConfig

	reqData.Name = plan.Name.ValueString()
	reqData.MachineType = plan.MachineType.ValueString()
	reqData.TemplateID = plan.TemplateID.ValueString()
	reqData.NetworkID = plan.NetworkID.ValueString()
	reqData.DiskSize = int(plan.DiskSize.ValueInt32())
	reqData.Region = plan.Region.ValueString()
	reqData.PublicIpType = plan.PublicIpType.ValueString()
	reqData.StartOnCreate = plan.StartOnCreate.ValueBool()

	var accessors []string
	// TODO: Consider better ways to do so
	for _, item := range plan.AccessorIds.Elements() {
		if strValue, ok := item.(types.String); ok {
			accessors = append(accessors, strValue.ValueString()) // Extract the string value
		}
	}
	reqData.AccessorIds = accessors

	// TODO: Check and implement remaining
	//
	// reqData.AutoSnapshotEnabled = plan.AutoSnapshotEnabled.ValueBool()
	// reqData.AutoSnapshotFrequency = plan.AutoSnapshotFrequency.ValueString()
	// reqData.AutoSnapshotSaveCount = int(plan.AutoSnapshotSaveCount.ValueInt64())
	// reqData.AutoShutdownEnabled = plan.AutoShutdownEnabled.ValueBool()
	// reqData.AutoShutdownTimeout = int(plan.AutoShutdownTimeout.ValueInt64())
	// reqData.AutoShutdownForce = plan.AutoShutdownForce.ValueBool()
	// reqData.TakeInitialSnapshot = plan.TakeInitialSnapshot.ValueBool()
	// reqData.RestorePointEnabled = plan.RestorePointEnabled.ValueBool()
	// reqData.RestorePointFrequency = plan.RestorePointFrequency.ValueString()
	// reqData.StartupScriptID = plan.StartupScriptID.ValueString()
	// reqData.EmailPassword = plan.EmailPassword.ValueBool()
	// reqData.StartOnCreate = plan.StartOnCreate.ValueBool()
	// reqData.EnableNvlink = plan.EnableNvlink.ValueBool()

	jsonData, err := json.MarshalIndent(reqData, "", " ")
	tflog.Info(ctx, "Sent req data: "+string(jsonData))

	// TODO: If start_on_create is true, implement waiting for machine to start
	machine, err := r.client.CreateMachine(reqData, ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating machine",
			"Could not create machine, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, "Created a resource with id "+machine.ID)

	// TODO: Remove indent JSON, it was initially added for better debug only
	machineJson, err := json.MarshalIndent(machine, "", " ")
	if err != nil {
		tflog.Error(ctx, "Could not marshal created machine struct: "+err.Error())
	}
	tflog.Info(ctx, "Created machine data: "+string(machineJson))

	// Map response body to schema and populate Computed attribute values.
	// Save response data into the Terraform state.
	// Only computed attributes must be updated here

	// TODO: Fill plan with all remaining values

	plan.ID = types.StringValue(machine.ID)
	plan.NetworkID = types.StringValue(machine.NetworkID)
	plan.CPUs = types.Int32Value(int32(machine.CPUs))
	plan.State = types.StringValue(machine.State)
	plan.OS = types.StringValue(machine.OS)
	plan.AgentType = types.StringValue(machine.AgentType)
	plan.PrivateIP = types.StringValue(machine.PrivateIP)

	// Nullable field
	if machine.PublicIP != nil {
		plan.PublicIP = types.StringValue(*machine.PublicIP)
	} else {
		plan.PublicIP = types.StringNull()
	}

	// TODO: Consider this, it may be useful
	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *machineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *machineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *machineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state machineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteMachine(state.ID.ValueString(), ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Machine",
			"Could not delete machine, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *machineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ppclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
