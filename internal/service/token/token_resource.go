package token

import (
	"context"
	"fmt"
	"time"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.ResourceWithConfigure = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *datafy.Client
}

type ResourceModel struct {
	AccountId   types.String         `tfsdk:"account_id"`
	TokenId     types.String         `tfsdk:"token_id"`
	Description types.String         `tfsdk:"description"`
	Ttl         timetypes.GoDuration `tfsdk:"ttl"`
	RoleIds     types.List           `tfsdk:"role_ids"`
	Secret      types.String         `tfsdk:"secret"`
	Expires     timetypes.RFC3339    `tfsdk:"expires"`
	CreatedAt   timetypes.RFC3339    `tfsdk:"created_at"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token"
}

func (r *Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {

}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create a Datafy token, which represents an access token associated with a Datafy account.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "The unique identifier of the Datafy account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description of the token.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ttl": schema.StringAttribute{
				CustomType:  timetypes.GoDurationType{},
				Description: "The expiration time of the token.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_ids": schema.ListAttribute{
				Description: "A list of role IDs associated with the token.",
				ElementType: types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"secret": schema.StringAttribute{
				Description: "The secret value of the token.",
				Computed:    true,
				Sensitive:   true,
			},
			"token_id": schema.StringAttribute{
				Description: "The unique identifier of the Datafy token.",
				Computed:    true,
			},
			"expires": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Description: "The time when the token will expire.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Description: "The time when the token was created.",
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

	elements := make([]types.String, 0, len(plan.RoleIds.Elements()))
	resp.Diagnostics.Append(plan.RoleIds.ElementsAs(ctx, &elements, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catr, err := r.client.CreateAccountToken(ctx, &datafy.CreateAccountTokenRequest{
		AccountId:   plan.AccountId.ValueString(),
		Description: plan.Description.ValueString(),
		Ttl: func() time.Duration {
			d, _ := time.ParseDuration(plan.Ttl.ValueString())
			return d
		}(),
		RoleIds: func() []string {
			res := make([]string, 0, len(elements))
			for _, e := range elements {
				res = append(res, e.ValueString())
			}
			return res
		}(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating account token",
			"Could not create account token: "+err.Error(),
		)
		return
	}

	plan.TokenId = types.StringValue(catr.AccountToken.TokenId)
	plan.Secret = types.StringValue(catr.AccountToken.Secret)
	plan.Expires = timetypes.NewRFC3339TimeValue(catr.AccountToken.Expires)
	plan.CreatedAt = timetypes.NewRFC3339TimeValue(catr.AccountToken.CreatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gat, err := r.client.GetAccountToken(ctx, &datafy.GetAccountTokenRequest{
		AccountId: state.AccountId.ValueString(),
		TokenId:   state.TokenId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error read account token",
			"Could not read account token: "+err.Error(),
		)
		return
	}

	state.Description = types.StringValue(gat.AccountToken.Description)
	state.RoleIds, _ = types.ListValueFrom(ctx, types.StringType, gat.AccountToken.RoleIds)
	state.Expires = timetypes.NewRFC3339TimeValue(gat.AccountToken.Expires)
	state.CreatedAt = timetypes.NewRFC3339TimeValue(gat.AccountToken.CreatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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

	_, err := r.client.DeleteAccountToken(ctx, &datafy.DeleteAccountTokenRequest{
		AccountId: state.AccountId.ValueString(),
		TokenId:   state.TokenId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error delete account token",
			"Could not delete account token: "+err.Error(),
		)
		return
	}
}
