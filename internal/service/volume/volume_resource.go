package volume

import (
	"context"
	"fmt"

	"github.com/datafy-io/terraform-provider-datafy/internal/datafy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *datafy.Client
}

type ResourceModel struct {
	Id               types.String `tfsdk:"id"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
	DiskSize         types.Int64  `tfsdk:"disk_size"`
	VolumeIops       types.Int64  `tfsdk:"volume_iops"`
	VolumeThroughput types.Int64  `tfsdk:"volume_throughput"`
	Encrypted        types.Bool   `tfsdk:"encrypted"`
	KmsKeyId         types.String `tfsdk:"kms_key_id"`
	Tags             types.Map    `tfsdk:"tags"`
	VolumeSizeGB     types.Int64  `tfsdk:"volume_size_gb"`
	TargetVolumeIds  types.List   `tfsdk:"target_volume_ids"`
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provisions a Datafy native datafied volume: two blank gp3 EBS target volumes pre-tagged for Datafy agent initialization. The agent completes setup on first attach.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The synthesized EBS volume ID used as the logical source identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"availability_zone": schema.StringAttribute{
				Description: "AWS availability zone in which to create the volume (e.g. us-east-1a).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disk_size": schema.Int64Attribute{
				Description: "Logical disk size in GiB presented to the OS.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"volume_iops": schema.Int64Attribute{
				Description: "Provisioned IOPS for the underlying gp3 EBS volumes. Defaults to 3000.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3000),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"volume_throughput": schema.Int64Attribute{
				Description: "Provisioned throughput in MiB/s for the underlying gp3 EBS volumes. Defaults to 125.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(125),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"encrypted": schema.BoolAttribute{
				Description: "Whether the EBS volumes are encrypted. Defaults to true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"kms_key_id": schema.StringAttribute{
				Description: "KMS key ID or ARN to use for volume encryption. Uses the default AWS-managed key if omitted.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.MapAttribute{
				Description: "Key-value tags to apply to the EBS volumes.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"volume_size_gb": schema.Int64Attribute{
				Description: "Actual size in GiB of each underlying EBS target volume as provisioned.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"target_volume_ids": schema.ListAttribute{
				Description: "EBS volume IDs of the two underlying target volumes.",
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	tags := map[string]string{}
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	result, err := r.client.CreateVolume(ctx, &datafy.CreateVolumeRequest{
		AvailabilityZone: plan.AvailabilityZone.ValueString(),
		DiskSize:         plan.DiskSize.ValueInt64(),
		VolumeIops:       plan.VolumeIops.ValueInt64(),
		VolumeThroughput: plan.VolumeThroughput.ValueInt64(),
		Encrypted:        plan.Encrypted.ValueBool(),
		KmsKeyId:         plan.KmsKeyId.ValueString(),
		Tags:             tags,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating volume", "Could not create volume: "+err.Error())
		return
	}

	targetList, diags := types.ListValueFrom(ctx, types.StringType, result.TargetVolumeIds)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(result.VolumeId)
	plan.TargetVolumeIds = targetList
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vol, err := r.client.GetVolume(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading volume", "Could not read volume: "+err.Error())
		return
	}
	if vol == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.AvailabilityZone = types.StringValue(vol.AvailabilityZone)
	state.DiskSize = types.Int64Value(int64(vol.DiskSize >> 30))
	state.VolumeIops = types.Int64Value(int64(vol.Iops))
	state.VolumeThroughput = types.Int64Value(int64(vol.Throughput))

	tagMap := make(map[string]attr.Value, len(vol.Tags))
	for _, t := range vol.Tags {
		tagMap[t.Key] = types.StringValue(t.Value)
	}
	tags, diags := types.MapValue(types.StringType, tagMap)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Tags = tags

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// All fields are RequiresReplace; Update is never called.
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteVolume(ctx, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting volume", "Could not delete volume: "+err.Error())
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
