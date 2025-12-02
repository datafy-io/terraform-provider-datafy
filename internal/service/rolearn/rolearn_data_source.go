package rolearn

import (
	"context"
	"fmt"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSourceWithConfigure = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *datafy.Client
}

type DataSourceModel struct {
	AccountId types.String `tfsdk:"account_id"`
	Arn       types.String `tfsdk:"arn"`
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_arn"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a specific Datafy role ARN.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "The unique identifier of the Datafy account.",
				Required:    true,
			},
			"arn": schema.StringAttribute{
				Description: "The Amazon Resource Name (ARN) of the IAM role.",
				Computed:    true,
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		resp.Diagnostics.AddError(
			"Provider Not Configured",
			"The provider has not been configured. Please ensure that you have set up the provider correctly.",
		)
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

	garar, err := d.client.GetAccountRoleArn(ctx, &datafy.GetAccountRoleArnRequest{
		AccountId: plan.AccountId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error read account role arn",
			"Could not read account role arn: "+err.Error(),
		)
		return
	}

	plan.Arn = types.StringValue(garar.AccountRoleArn.RoleArn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
