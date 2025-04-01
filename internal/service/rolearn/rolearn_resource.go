package rolearn

import (
	"context"
	"fmt"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *datafy.Client
}

type ResourceModel struct {
	AccountId types.String `tfsdk:"account_id"`
	Arn       types.String `tfsdk:"arn"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_arn"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Datafy role ARN, which represents an AWS IAM role associated with a Datafy account.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "The unique identifier of the Datafy account.",
				Required:    true,
			},
			"arn": schema.StringAttribute{
				Description: "The Amazon Resource Name (ARN) of the IAM role.",
				Required:    true,
			},
		},
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CreateAccountRoleArn(ctx, &datafy.CreateAccountRoleArnRequest{
		AccountId: plan.AccountId.ValueString(),
		Arn:       plan.Arn.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating account role arn",
			"Could not create account role arn: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	garar, err := r.client.GetAccountRoleArn(ctx, &datafy.GetAccountRoleArnRequest{
		AccountId: state.AccountId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error read account role arn",
			"Could not read account role arn: "+err.Error(),
		)
		return
	}

	state.Arn = types.StringValue(garar.AccountRoleArn.RoleArn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateAccountRoleArn(ctx, &datafy.UpdateAccountRoleArnRequest{
		AccountId: plan.AccountId.ValueString(),
		Arn:       plan.Arn.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update account role arn",
			"Could not update account role arn: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteAccountRoleArn(ctx, &datafy.DeleteAccountRoleArnRequest{
		AccountId: state.AccountId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error delete account role arn",
			"Could not delete account role arn: "+err.Error(),
		)
		return
	}
}
