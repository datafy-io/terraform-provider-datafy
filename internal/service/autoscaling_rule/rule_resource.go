package autoscaling_rule

import (
	"context"
	"fmt"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Id        types.String         `tfsdk:"id"`
	AccountId types.String         `tfsdk:"account_id"`
	Active    types.Bool           `tfsdk:"active"`
	Mode      types.String         `tfsdk:"mode"`
	Rule      jsontypes.Normalized `tfsdk:"rule"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_autoscaling_rule"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create a Datafy Autoscaling Rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the Datafy Autoscaling Rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_id": schema.StringAttribute{
				Description: "The unique identifier of the Datafy account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"active": schema.BoolAttribute{
				Description: "Indicates whether the autoscaling rule is active or not.",
				Required:    true,
			},
			"mode": schema.StringAttribute{
				Description: "The mode of the autoscaling rule.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("include", "exclude"),
				},
			},
			"rule": schema.StringAttribute{
				CustomType:  jsontypes.NormalizedType{},
				Description: "The autoscaling rule policy.",
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

	caarr, err := r.client.CreateAccountAutoscalingRule(ctx, &datafy.CreateAccountAutoscalingRuleRequest{
		AccountId: plan.AccountId.ValueString(),
		Active:    plan.Active.ValueBool(),
		Mode:      plan.Mode.ValueString(),
		Rule:      plan.Rule.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating account autoscaling rule",
			"Could not create account autoscaling rule: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(caarr.AutoscalingRule.RuleId)
	plan.AccountId = types.StringValue(caarr.AutoscalingRule.AccountId)
	plan.Active = types.BoolValue(caarr.AutoscalingRule.Active)
	plan.Mode = types.StringValue(caarr.AutoscalingRule.Mode)
	plan.Rule = jsontypes.NewNormalizedValue(caarr.AutoscalingRule.Rule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gaarr, err := r.client.GetAccountAutoscalingRule(ctx, &datafy.GetAccountAutoscalingRuleRequest{
		AccountId: state.AccountId.ValueString(),
		RuleId:    state.Id.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error read account autoscaling rule",
			"Could not read account autoscaling rule: "+err.Error(),
		)
		return
	}

	state.Id = types.StringValue(gaarr.AutoscalingRule.RuleId)
	state.AccountId = types.StringValue(gaarr.AutoscalingRule.AccountId)
	state.Active = types.BoolValue(gaarr.AutoscalingRule.Active)
	state.Mode = types.StringValue(gaarr.AutoscalingRule.Mode)
	state.Rule = jsontypes.NewNormalizedValue(gaarr.AutoscalingRule.Rule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateAccountAutoscalingRule(ctx, &datafy.UpdateAccountAutoscalingRuleRequest{
		AccountId: plan.AccountId.ValueString(),
		RuleId:    plan.Id.ValueString(),
		Active:    plan.Active.ValueBool(),
		Mode:      plan.Mode.ValueString(),
		Rule:      plan.Rule.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update account autoscaling rule",
			"Could not update account autoscaling rule: "+err.Error(),
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

	_, err := r.client.DeleteAccountAutoscalingRule(ctx, &datafy.DeleteAccountAutoscalingRuleRequest{
		AccountId: state.AccountId.ValueString(),
		RuleId:    state.Id.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error delete account autoscaling rule",
			"Could not delete account autoscaling rule: "+err.Error(),
		)
		return
	}
}
