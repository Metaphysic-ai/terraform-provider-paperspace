package provider

import (
	"context"
	"os"
	"terraform-provider-paperspace/internal/ppclient"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure paperspaceProvider satisfies various provider interfaces.
var _ provider.Provider = &paperspaceProvider{}

// paperspaceProvider defines the provider implementation.
type paperspaceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// paperspaceProviderModel describes the provider data model.
type paperspaceProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

func (p *paperspaceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "paperspace"
	resp.Version = p.version
}

func (p *paperspaceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				// TODO: Add MarkdownDescription
				Optional:  true, // May be set over ENV VAR if true
				Sensitive: true,
			},

			// TODO: Consider 'namespace' here
		},
	}
}

func (p *paperspaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration

	var config paperspaceProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Paperspace API Password",
			"The provider cannot create the Paperspace API client as there is an unknown configuration value for the Paperspace API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PAPERSPACE_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	api_key := os.Getenv("PAPERSPACE_PASSWORD")

	if !config.APIKey.IsNull() {
		api_key = config.APIKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Paperspace API Password",
			"The provider cannot create the Paperspace API client as there is a missing or empty value for the Paperspace API password. "+
				"Set the password value in the configuration or use the PAPERSPACE_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Paperspace client using the configuration values
	// client := http.DefaultClient

	// Create a new Paperspace client using the configuration values
	client, err := ppclient.NewClient(nil, &api_key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Paperspace API Client",
			"An unexpected error occurred when creating the Paperspace API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Paperspace Client Error: "+err.Error(),
		)
		return
	}

	// Make the Paperspace client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *paperspaceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *paperspaceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCustomTemplatesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &paperspaceProvider{
			version: version,
		}
	}
}
