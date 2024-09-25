package provider

import (
	"context"
	"fmt"

	"terraform-provider-paperspace/internal/ppclient"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &customTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &customTemplatesDataSource{}
)

// NewCustomTemplatesDataSource is a helper function to simplify the provider implementation.
func NewCustomTemplatesDataSource() datasource.DataSource {
	return &customTemplatesDataSource{}
}

// Allow your data source type to store a reference to the Paperspace client.
type customTemplatesDataSource struct {
	client *ppclient.Client
}

// Metadata returns the data source type name.
func (d *customTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_templates"
}

// Schema defines the schema for the data source.
// The data source uses the Schema method to define the acceptable configuration and state attribute names and types.
func (d *customTemplatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// TODO: Consider removing custom_templates key on top level
			"custom_templates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							// Sensitive: false,

							// TODO: Add descriptions
							// Example: MarkdownDescription: "Example identifier",
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"agent_type": schema.StringAttribute{
							Computed: true,
						},
						"operating_system_label": schema.StringAttribute{
							Computed: true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"default_size_gb": schema.Int32Attribute{
							Computed: true,
						},

						// TODO: Add nested attributes available_machine_types
						// Example:
						// https://github.com/hashicorp/terraform-provider-hashicups/blob/main/03-read-coffees/internal/provider/coffees_data_source.go
						// "available_machine_types":

						"parent_machine_id": schema.StringAttribute{
							Computed: true,
						},
						"dt_created": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

//
// Data model types
//

// customTemplatesDataSourceModel maps the data source schema data.
type customTemplatesDataSourceModel struct {
	CustomTemplates []customTemplatesModel `tfsdk:"custom_templates"`
}

// customTemplatesModel maps customTemplates schema data.
type customTemplatesModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	AgentType            types.String `tfsdk:"agent_type"`
	OperatingSystemLabel types.String `tfsdk:"operating_system_label"`
	Region               types.String `tfsdk:"region"`
	DefaultSizeGb        types.Int32  `tfsdk:"default_size_gb"`
	ParentMachineID      types.String `tfsdk:"parent_machine_id"`
	DtCreated            types.String `tfsdk:"dt_created"`
}

//
// Configure
//

// Configure adds the provider configured client to the data source.
func (d *customTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ppclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ppclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

//
// Read
//

// Read refreshes the Terraform state with the latest data.
func (d *customTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state customTemplatesDataSourceModel

	customTemplates, err := d.client.GetCustomTemplates()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Paperspace CustomTemplates",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, customTemplate := range customTemplates {
		customTemplateState := customTemplatesModel{
			ID:                   types.StringValue(customTemplate.ID),
			Name:                 types.StringValue(customTemplate.Name),
			AgentType:            types.StringValue(customTemplate.AgentType),
			OperatingSystemLabel: types.StringValue(customTemplate.OperatingSystemLabel),
			Region:               types.StringValue(customTemplate.Region),
			DefaultSizeGb:        types.Int32Value(int32(customTemplate.DefaultSizeGb)), // TODO: Consider allowing positive values only
			ParentMachineID:      types.StringValue(customTemplate.ParentMachineID),
			DtCreated:            types.StringValue(customTemplate.DtCreated),
		}

		state.CustomTemplates = append(state.CustomTemplates, customTemplateState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}
