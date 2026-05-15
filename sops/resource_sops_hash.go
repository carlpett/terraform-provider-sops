package sops

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type hashResourceModel struct {
	InputWO tftypes.String `tfsdk:"input_wo"`
	Hash    tftypes.String `tfsdk:"hash"`
}

var _ resource.Resource = &hashResource{}
var _ resource.ResourceWithModifyPlan = &hashResource{}

func newHashResource() resource.Resource {
	return &hashResource{}
}

type hashResource struct {
}

func (r *hashResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hash"
}

func (r *hashResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Accepts an ephemeral input and returns a hash of the input. Useful for triggering a write-only resource update when an ephemeral value changes.",
		Attributes: map[string]schema.Attribute{
			"input_wo": schema.StringAttribute{
				MarkdownDescription: "Write-only input (not stored in state).",
				Required:            true,
				WriteOnly:           true,
			},
			"hash": schema.StringAttribute{
				MarkdownDescription: "Returns sha256 hash of `input_wo`.",
				Computed:            true,
			},
		},
	}
}

func (r *hashResource) Configure(_ context.Context, _ resource.ConfigureRequest, _ *resource.ConfigureResponse) {
}

func (r *hashResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// create or destroy
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var config hashResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// An unknown write-only input should trigger a plan
	if config.InputWO.IsUnknown() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("hash"), tftypes.StringUnknown())...)
		return
	}

	// Plan values for hash must always contain the hash of the write-only input
	newHash := hashInput(config.InputWO.ValueString())
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("hash"), tftypes.StringValue(newHash))...)
}

func (r *hashResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config hashResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hash := tftypes.StringValue(hashInput(config.InputWO.ValueString()))
	config.Hash = hash

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (r *hashResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *hashResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan hashResourceModel
	var config hashResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hash := tftypes.StringValue(hashInput(config.InputWO.ValueString()))
	plan.Hash = hash

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *hashResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func hashInput(s string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(s))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
