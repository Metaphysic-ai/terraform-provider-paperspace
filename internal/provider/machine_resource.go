package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"terraform-provider-paperspace/internal/psclient"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

// Maps the resource schema data.
// State/Plan structure.
type machineResourceModel struct {
	// MachineCreateConfig
	Name                   types.String `tfsdk:"name"`         // required
	MachineType            types.String `tfsdk:"machine_type"` // required
	TemplateID             types.String `tfsdk:"template_id"`  // required
	DiskSize               types.Int64  `tfsdk:"disk_size"`    // required
	Region                 types.String `tfsdk:"region"`       // required
	NetworkID              types.String `tfsdk:"network_id"`
	AutoSnapshotEnabled    types.Bool   `tfsdk:"auto_snapshot_enabled"`
	AutoSnapshotFrequency  types.String `tfsdk:"auto_snapshot_frequency"`
	AutoSnapshotSaveCount  types.Int64  `tfsdk:"auto_snapshot_save_count"`
	AutoShutdownEnabled    types.Bool   `tfsdk:"auto_shutdown_enabled"`
	AutoShutdownTimeout    types.Int64  `tfsdk:"auto_shutdown_timeout"`
	AutoShutdownForce      types.Bool   `tfsdk:"auto_shutdown_force"`
	RestorePointEnabled    types.Bool   `tfsdk:"restore_point_enabled"`
	RestorePointFrequency  types.String `tfsdk:"restore_point_frequency"`
	RestorePointSnapshotID types.String `tfsdk:"restore_point_snapshot_id"`
	PublicIPType           types.String `tfsdk:"public_ip_type"`
	EnableNvlink           types.Bool   `tfsdk:"enable_nvlink"`
	TakeInitialSnapshot    types.Bool   `tfsdk:"take_initial_snapshot"`
	StartupScriptID        types.String `tfsdk:"startup_script_id"`
	EmailPassword          types.Bool   `tfsdk:"email_password"`
	AccessorIDs            types.List   `tfsdk:"accessor_ids"`

	// Computed only
	ID           types.String  `tfsdk:"id"`
	RegionFull   types.String  `tfsdk:"region_full"`
	CPUs         types.Int64   `tfsdk:"cpus"`
	State        types.String  `tfsdk:"state"`
	OS           types.String  `tfsdk:"os"`
	AgentType    types.String  `tfsdk:"agent_type"`
	PublicIP     types.String  `tfsdk:"public_ip"`
	PrivateIP    types.String  `tfsdk:"private_ip"`
	RAM          types.String  `tfsdk:"ram"`
	StorageTotal types.String  `tfsdk:"storage_total"`
	StorageUsed  types.String  `tfsdk:"storage_used"`
	UsageRate    types.Float64 `tfsdk:"usage_rate"`
	StorageRate  types.Float64 `tfsdk:"storage_rate"`
	DtCreated    types.String  `tfsdk:"dt_created"`
	DtModified   types.String  `tfsdk:"dt_modified"`

	//// Note: These fields are omitted
	// DtDeleted    types.String  `tfsdk:"dtDeleted"`
	// Reservation  *Reservation  `tfsdk:"reservation"`
	// Accelerators []Accelerator `tfsdk:"accelerators"`
}

// machineResource is the resource implementation.
type machineResource struct {
	// Allow resource to store a reference to the client
	client *psclient.Client
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
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the machine.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the new machine.",
				Required:            true,
			},
			"machine_type": schema.StringAttribute{
				MarkdownDescription: "The machine type. Updates to this field will trigger a stop/start of the machine.",
				Required:            true,
			},
			"template_id": schema.StringAttribute{
				MarkdownDescription: "The template ID.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disk_size": schema.Int64Attribute{
				MarkdownDescription: "The disk size in gigabytes. Updates to this field will trigger a stop/start of the machine.",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(50, 100, 250, 500, 1000, 2000),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region to create the machine in.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The network ID. You can migrate machines between private networks and from the default network to a private network." +
					" It is not possible to migrate a machine back to the default network." +
					" If this is required, please file a support ticket.",
				Optional: true,
				Computed: true,
				// stringplanmodifier.UseStateForUnknown() is not used here intentionally.
				// If private network ID is not set explicitly in Terraform configuration,
				// the ID of default one MUST not be passed during machine update.
				// Otherwise API returns 404 error.
				// So, cannot use ID from state for unknown, because it contains ID of default network.
			},
			"region_full": schema.StringAttribute{
				MarkdownDescription: "Full machine region name.",
				Computed:            true,
			},
			"cpus": schema.Int64Attribute{
				MarkdownDescription: "Number of CPUs.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "Desired state of the machine. Possible values: `off`, `ready`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("off"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"off", "ready"}...),
				},
			},
			"os": schema.StringAttribute{
				MarkdownDescription: "Operating system of the machine.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"agent_type": schema.StringAttribute{
				MarkdownDescription: "Agent type of the machine.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"public_ip_type": schema.StringAttribute{
				MarkdownDescription: "The public IP type. Possible values: `static`, `dynamic`, `none`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("dynamic"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"static", "dynamic", "none"}...),
				},
			},
			"public_ip": schema.StringAttribute{
				MarkdownDescription: "Public IP address of the machine.",
				Computed:            true,
			},
			"private_ip": schema.StringAttribute{
				MarkdownDescription: "Private IP address of the machine.",
				Computed:            true,
			},
			// Auto Snapshot
			"auto_snapshot_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable auto snapshots.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"auto_snapshot_frequency": schema.StringAttribute{
				MarkdownDescription: "The auto snapshot frequency. Possible values: `hourly`, `daily`, `weekly`, `monthly`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"hourly", "daily", "weekly", "monthly"}...),
				},
			},
			"auto_snapshot_save_count": schema.Int64Attribute{
				MarkdownDescription: "The number of auto snapshots to save. Must be between 1 and 9 if `auto_snapshot_enabled` is true.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 9),
				},
			},
			// Auto Shutdown
			"auto_shutdown_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable auto shutdown.",
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"auto_shutdown_timeout": schema.Int64Attribute{
				MarkdownDescription: "The auto shutdown timeout in hours. Must be set if `auto_shutdown_enabled` is true. " +
					"May be troubles with updating the value, seems like Paperspace API issue." +
					"Disable auto shutdown and then enable with different option to update.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.OneOf(1, 8, 24, 168),
				},
			},
			"auto_shutdown_force": schema.BoolAttribute{
				MarkdownDescription: "Whether to force shutdown the machine. " +
					"May be troubles with updating the value, seems like Paperspace API issue." +
					"Disable auto shutdown and then enable with different option to update.",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			// Restore Point (only computed, user input is not implemented yet)
			"restore_point_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether to use initial snapshot as a restore point.",
				Computed:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"restore_point_frequency": schema.StringAttribute{
				MarkdownDescription: "The restore point frequency. Possible values: `shutdown`.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"shutdown"}...),
				},
			},
			"restore_point_snapshot_id": schema.StringAttribute{
				MarkdownDescription: "The restore point snapshot ID.",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			// Attributes which apply only on creation
			"enable_nvlink": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable NVLink.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(), // Forces resource replacement if changed
				},
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"take_initial_snapshot": schema.BoolAttribute{
				MarkdownDescription: "Whether to take an initial snapshot. Applies only on resource creation.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"startup_script_id": schema.StringAttribute{
				MarkdownDescription: "The startup script ID. Forces resource replacement if changed.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email_password": schema.BoolAttribute{
				MarkdownDescription: "Whether to email the password. Applies only on resource creation.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			// Computed only
			"ram": schema.StringAttribute{
				MarkdownDescription: "RAM amount of the machine.",
				Computed:            true,
			},
			"storage_total": schema.StringAttribute{
				MarkdownDescription: "Storage total of the machine.",
				Computed:            true,
			},
			"storage_used": schema.StringAttribute{
				MarkdownDescription: "Storage used of the machine.",
				Computed:            true,
			},
			"usage_rate": schema.Float64Attribute{
				MarkdownDescription: "Usage rate of the machine.",
				Computed:            true,
			},
			"storage_rate": schema.Float64Attribute{
				MarkdownDescription: "Storage rate of the machine.",
				Computed:            true,
			},
			"dt_created": schema.StringAttribute{
				MarkdownDescription: "Created date timestamp of the machine.",
				Computed:            true,
			},
			"dt_modified": schema.StringAttribute{
				MarkdownDescription: "Modified date timestamp of the machine.",
				Computed:            true,
			},
			// TODO: Implement change on resource update, so machine accessors are being updated
			"accessor_ids": schema.ListAttribute{
				MarkdownDescription: "The IDs of users to grant access to the machine. Applies only on resource creation.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// TODO: Add timeouts https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts
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

	accessors := []string{}
	for _, item := range plan.AccessorIDs.Elements() {
		if strValue, ok := item.(types.String); ok {
			accessors = append(accessors, strValue.ValueString()) // Extract the string value
		}
	}

	reqData := psclient.MachineCreateConfig{
		Name:                  plan.Name.ValueString(),        // required
		MachineType:           plan.MachineType.ValueString(), // required
		TemplateID:            plan.TemplateID.ValueString(),  // required
		DiskSize:              plan.DiskSize.ValueInt64(),     // required
		Region:                plan.Region.ValueString(),      // required
		NetworkID:             plan.NetworkID.ValueString(),
		PublicIPType:          plan.PublicIPType.ValueString(),
		StartOnCreate:         plan.State.ValueString() == "ready",
		AutoSnapshotEnabled:   getValueBoolPointer(plan.AutoSnapshotEnabled),
		AutoSnapshotFrequency: plan.AutoSnapshotFrequency.ValueString(),
		AutoSnapshotSaveCount: getValueInt64Pointer(plan.AutoSnapshotSaveCount),
		AutoShutdownEnabled:   getValueBoolPointer(plan.AutoShutdownEnabled),
		AutoShutdownTimeout:   getValueInt64Pointer(plan.AutoShutdownTimeout),
		AutoShutdownForce:     getValueBoolPointer(plan.AutoShutdownForce),
		EnableNvlink:          getValueBoolPointer(plan.EnableNvlink),
		TakeInitialSnapshot:   getValueBoolPointer(plan.TakeInitialSnapshot),
		StartupScriptID:       plan.StartupScriptID.ValueString(),
		EmailPassword:         getValueBoolPointer(plan.EmailPassword),
		AccessorIDs:           accessors,
	}

	jsonData, _ := json.MarshalIndent(reqData, "", " ")
	tflog.Info(ctx, "Sending create req data: "+string(jsonData))

	machine, err := r.client.CreateMachine(reqData)
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

	plan.ID = types.StringValue(machine.ID)
	fillStateWithMachineData(&plan, machine)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *machineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state machineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed machine value from Paperspace
	machine, err := r.client.GetMachine(state.ID.ValueString())
	if err != nil {
		// Avoid error due to 404 status. Machine may be deleted outside provider, so handle this as expected case.
		if strings.Contains(err.Error(), "status: 404") {
			tflog.Warn(ctx, fmt.Sprintf("Machine %s not found, removing from state", state.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Paperspace Machine",
			"Could not read Paperspace machine ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// TODO: Remove log, or remove indent
	machineJson, err := json.MarshalIndent(machine, "", " ")
	if err != nil {
		tflog.Error(ctx, "Could not marshal fetched machine struct: "+err.Error())
	}
	tflog.Info(ctx, "Fetched machine data: "+string(machineJson))

	// ID not needed here
	fillStateWithMachineData(&state, machine)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Updates the resource and sets the updated Terraform state on success.
func (r *machineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Fetch the entire plan and prior state
	var plan, state machineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	machineID := plan.ID.ValueString()
	machineStateCurrent := state.State.ValueString()
	machineStateTarget := plan.State.ValueString()

	if (machineStateCurrent != "off") && (machineStateCurrent != "ready") {
		resp.Diagnostics.AddError(
			"Error updating machine",
			fmt.Sprintf("Could not update machine, it must be 'off' or 'ready', but it is '%s'", machineStateCurrent),
		)
		return
	}

	// If machine type or disk size is changed, machine must be stopped before such update
	if !plan.MachineType.Equal(state.MachineType) || !plan.DiskSize.Equal(state.DiskSize) {
		// Make sure machine is off
		tflog.Info(ctx, "Stopping machine before update, ID: "+machineID)
		err := r.client.ManageMachineState(machineID, psclient.MachineStateOff)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error stopping Paperspace machine",
				"Could not stop Paperspace machine ID "+machineID+": "+err.Error(),
			)
			return
		}
	}

	// Generate API request body from plan

	reqData := psclient.MachineUpdateConfig{
		Name:                  plan.Name.ValueString(),
		MachineType:           plan.MachineType.ValueString(),
		NetworkID:             plan.NetworkID.ValueString(),
		DiskSize:              plan.DiskSize.ValueInt64(),
		PublicIPType:          plan.PublicIPType.ValueString(),
		AutoSnapshotEnabled:   getValueBoolPointer(plan.AutoSnapshotEnabled),
		AutoSnapshotFrequency: plan.AutoSnapshotFrequency.ValueString(),
		AutoSnapshotSaveCount: getValueInt64Pointer(plan.AutoSnapshotSaveCount),
		AutoShutdownEnabled:   getValueBoolPointer(plan.AutoShutdownEnabled),
		AutoShutdownTimeout:   getValueInt64Pointer(plan.AutoShutdownTimeout),
		AutoShutdownForce:     getValueBoolPointer(plan.AutoShutdownForce),
	}

	jsonData, _ := json.MarshalIndent(reqData, "", " ")
	tflog.Info(ctx, "Sending update req data: "+string(jsonData))

	err := r.client.UpdateMachine(machineID, reqData)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating machine",
			"Could not update machine, unexpected error: "+err.Error(),
		)
		return
	}

	// Start/stop the machine based on target state
	tflog.Info(ctx, fmt.Sprintf("Ensuring machine '%s' is '%s'", machineID, machineStateTarget))
	err = r.client.ManageMachineState(machineID, machineStateTarget)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error managing Paperspace machine state",
			"Could update state of Paperspace machine ID "+machineID+": "+err.Error(),
		)
		return
	}

	// Fetch updated machine
	updatedMachine, err := r.client.GetMachine(machineID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading updated Paperspace machine",
			"Could not read Paperspace machine ID "+machineID+": "+err.Error(),
		)
		return
	}

	// TODO: Remove indent JSON, it was initially added for better debug only
	updatedMachineJson, err := json.MarshalIndent(updatedMachine, "", " ")
	if err != nil {
		tflog.Error(ctx, "Could not marshal updated machine struct: "+err.Error())
	}
	tflog.Info(ctx, "Updated machine data: "+string(updatedMachineJson))

	fillStateWithMachineData(&plan, updatedMachine)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	machineID := state.ID.ValueString()
	err := r.client.DeleteMachine(machineID)
	if err != nil {
		if strings.Contains(err.Error(), "status: 404") {
			tflog.Info(ctx, fmt.Sprintf("Machine %s not found, assuming already deleted", machineID))
			return
		}

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

func fillStateWithMachineData(state *machineResourceModel, machine *psclient.Machine) {
	state.Name = types.StringValue(machine.Name)
	state.State = types.StringValue(machine.State)
	state.OS = types.StringValue(machine.OS)
	state.MachineType = types.StringValue(machine.MachineType)
	state.AgentType = types.StringValue(machine.AgentType)
	state.CPUs = types.Int64Value(machine.CPUs)
	state.RAM = types.StringValue(machine.RAM)
	state.StorageTotal = types.StringValue(machine.StorageTotal)
	state.StorageUsed = types.StringValue(machine.StorageUsed)
	state.RegionFull = types.StringValue(machine.RegionFull)
	state.PrivateIP = types.StringValue(machine.PrivateIP)
	state.NetworkID = types.StringValue(machine.NetworkID)
	state.PublicIP = types.StringPointerValue(machine.PublicIP) // Nullable field
	state.PublicIPType = types.StringValue(machine.PublicIPType)
	state.AutoShutdownEnabled = types.BoolValue(machine.AutoShutdownEnabled)
	state.AutoShutdownTimeout = types.Int64PointerValue(machine.AutoShutdownTimeout) // Nullable field
	state.AutoShutdownForce = types.BoolValue(machine.AutoShutdownForce)
	state.AutoSnapshotEnabled = types.BoolValue(machine.AutoSnapshotEnabled)
	state.AutoSnapshotFrequency = types.StringPointerValue(machine.AutoSnapshotFrequency) // Nullable field
	state.AutoSnapshotSaveCount = types.Int64PointerValue(machine.AutoSnapshotSaveCount)  // Nullable field
	state.RestorePointEnabled = types.BoolValue(machine.RestorePointEnabled)
	state.RestorePointFrequency = types.StringPointerValue(machine.RestorePointFrequency)   // Nullable field
	state.RestorePointSnapshotID = types.StringPointerValue(machine.RestorePointSnapshotID) // Nullable field
	state.UsageRate = types.Float64Value(machine.UsageRate)
	state.StorageRate = types.Float64Value(machine.StorageRate)
	state.DtCreated = types.StringValue(machine.DtCreated)
	state.DtModified = types.StringValue(machine.DtModified)
}

// Private

// Returns nil for unknown and ValueBoolPointer for known.
func getValueBoolPointer(attr basetypes.BoolValue) *bool {
	if attr.IsUnknown() {
		return nil
	}
	return attr.ValueBoolPointer()
}

// Returns nil for unknown and ValueInt64Pointer for known.
func getValueInt64Pointer(attr basetypes.Int64Value) *int64 {
	if attr.IsUnknown() {
		return nil
	}
	return attr.ValueInt64Pointer()
}
