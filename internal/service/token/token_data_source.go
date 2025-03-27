package token

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	AccountId   types.String      `tfsdk:"account_id"`
	TokenId     types.String      `tfsdk:"token_id"`
	Description types.String      `tfsdk:"description"`
	RoleIds     types.List        `tfsdk:"role_ids"`
	Expires     timetypes.RFC3339 `json:"expires"`
	CreatedAt   timetypes.RFC3339 `json:"created_at"`
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Datafy token data source",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "account id",
				Required:    true,
			},
			"token_id": schema.StringAttribute{
				Description: "account token id",
				Required:    true,
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
			"Unexpected Resource Configure Type",
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

	gat, err := d.client.GetAccountToken(ctx, &datafy.GetAccountTokenRequest{
		AccountId: plan.AccountId.ValueString(),
		TokenId:   plan.TokenId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error read account token",
			"Could not read account token: "+err.Error(),
		)
		return
	}

	plan.Description = types.StringValue(gat.AccountToken.Description)
	plan.RoleIds, _ = types.ListValueFrom(ctx, types.StringType, gat.AccountToken.RoleIds)
	plan.Expires = timetypes.NewRFC3339TimeValue(gat.AccountToken.Expires)
	plan.CreatedAt = timetypes.NewRFC3339TimeValue(gat.AccountToken.CreatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
