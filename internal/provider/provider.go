package provider

import (
	"context"
	"os"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/datafy-io/terraform-provider-datafy/internal/service/account"
	"github.com/datafy-io/terraform-provider-datafy/internal/service/rolearn"
	"github.com/datafy-io/terraform-provider-datafy/internal/service/token"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure DatafyProvider satisfies provider interfaces.
var _ provider.Provider = &DatafyProvider{}

type DatafyProvider struct {
	version string
}

type DatafyProviderConfig struct {
	Token    types.String `tfsdk:"token"`
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *DatafyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "datafy"
	resp.Version = p.version
}

func (p *DatafyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "Datafy token. Can also be configured using the `DATAFY_TOKEN` environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"endpoint": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *DatafyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config DatafyProviderConfig

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Token",
			"Cannot create Datafy Provider as there is an unknown configuration value.",
		)
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Endpoint",
			"Cannot create Datafy Provider as there is an unknown configuration value.",
		)
	}

	datafyToken := os.Getenv("DATAFY_TOKEN")
	if !config.Token.IsNull() {
		datafyToken = config.Token.ValueString()
	}

	datafyEndpoint := os.Getenv("DATAFY_ENDPOINT")
	if !config.Endpoint.IsNull() {
		datafyEndpoint = config.Endpoint.ValueString()
	}
	if datafyEndpoint == "" {
		datafyEndpoint = "https://api.datafy.io"
	}

	if datafyToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Token",
			"Cannot create Datafy Provider as there is a missing or empty value.",
		)
	}

	client := datafy.NewClient(datafyToken, datafyEndpoint)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *DatafyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		account.NewResource,
		rolearn.NewResource,
		token.NewResource,
	}
}

func (p *DatafyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		account.NewDataSource,
		rolearn.NewDataSource,
		token.NewDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DatafyProvider{
			version: version,
		}
	}
}
