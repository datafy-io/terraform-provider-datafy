package autoscaling_rule

import (
	"context"
	"fmt"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
	Id        types.String         `tfsdk:"id"`
	AccountId types.String         `tfsdk:"accountId"`
	Active    types.Bool           `tfsdk:"active"`
	Rule      jsontypes.Normalized `tfsdk:"rule"`
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_autoscaling_rule"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a specific Datafy Autoscaling Rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the Datafy Autoscaling Rule.",
				Required:    true,
			},
			"account_id": schema.StringAttribute{
				Description: "The unique identifier of the Datafy account.",
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

	gaarr, err := d.client.GetAccountAutoscalingRule(ctx, &datafy.GetAccountAutoscalingRuleRequest{
		AccountId: plan.AccountId.ValueString(),
		RuleId:    plan.Id.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error read account autoscaling rule",
			"Could not read account autoscaling rule: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(gaarr.AutoscalingRule.RuleId)
	plan.AccountId = types.StringValue(gaarr.AutoscalingRule.AccountId)
	plan.Active = types.BoolValue(gaarr.AutoscalingRule.Active)
	plan.Rule = jsontypes.NewNormalizedValue(string(gaarr.AutoscalingRule.Rule))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
