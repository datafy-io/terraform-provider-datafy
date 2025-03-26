package provider

import (
	"context"

	"github.com/datafy-io/terraform-provider-datafy/internal/service/account"
	"github.com/datafy-io/terraform-provider-datafy/internal/service/rolearn"
	"github.com/datafy-io/terraform-provider-datafy/internal/service/token"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure DatafyProvider satisfies provider interfaces.
var _ provider.Provider = &DatafyProvider{}

type DatafyProvider struct {
	version string
}

type DatafyProviderModel struct {
}

func (p *DatafyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "datafy"
	resp.Version = p.version
}

func (p *DatafyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *DatafyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data DatafyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
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
