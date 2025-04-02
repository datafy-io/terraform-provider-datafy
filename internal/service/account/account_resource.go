package account

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
	Name            types.String `tfsdk:"name"`
	Id              types.String `tfsdk:"id"`
	ParentAccountId types.String `tfsdk:"parent_account_id"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create a Datafy account",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the Datafy account.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the Datafy account.",
				Computed:    true,
			},
			"parent_account_id": schema.StringAttribute{
				Description: "The unique identifier of the parent Datafy account",
				Computed:    true,
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

	car, err := r.client.CreateAccount(ctx, &datafy.CreateAccountRequest{
		AccountName: plan.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating account",
			"Could not create account: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(car.Account.AccountId)
	plan.ParentAccountId = types.StringValue(car.Account.ParentAccountId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gcr, err := r.client.GetAccount(ctx, &datafy.GetAccountRequest{
		AccountId: state.Id.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error read account",
			"Could not read account: "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(gcr.Account.AccountName)
	state.ParentAccountId = types.StringValue(gcr.Account.ParentAccountId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateAccount(ctx, &datafy.UpdateAccountRequest{
		AccountId:   plan.Id.ValueString(),
		AccountName: plan.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update account",
			"Could not update account: "+err.Error(),
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

	_, err := r.client.DeleteAccount(ctx, &datafy.DeleteAccountRequest{
		AccountId: state.Id.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error delete account",
			"Could not delete account: "+err.Error(),
		)
		return
	}
}
