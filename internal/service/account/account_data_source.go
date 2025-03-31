package account

import (
	"context"
	"fmt"
	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *datafy.Client
}

type DataSourceModel struct {
	Name            types.String `tfsdk:"name"`
	Id              types.String `tfsdk:"id"`
	ParentAccountId types.String `tfsdk:"parent_account_id"`
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Datafy account data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "account id",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "account name",
				Computed:    true,
			},
			"parent_account_id": schema.StringAttribute{
				Description: "parent account id",
				Computed:    true,
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*datafy.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *datafy.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gcr, err := d.client.GetAccount(ctx, &datafy.GetAccountRequest{
		AccountId: plan.Id.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error read account",
			"Could not read account: "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(gcr.Account.AccountName)
	plan.ParentAccountId = types.StringValue(gcr.Account.ParentAccountId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
