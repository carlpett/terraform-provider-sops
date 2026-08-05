package sops

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &fileResource{}

func newFileResource() resource.Resource {
	return &fileResource{}
}

type fileResource struct{}

type fileResourceModel struct {
	OutputType      types.String  `tfsdk:"output_type"`
	DestinationFile types.String  `tfsdk:"destination_file"`
	DataWo          types.Dynamic `tfsdk:"data_wo"`
	DataWoVersion   types.String  `tfsdk:"data_wo_version"`
	Id              types.String  `tfsdk:"id"`
}

func (r *fileResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "sops_file"
}

func (r *fileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Encrypt content with sops and write it to a file on disk. This is the counterpart to the " +
			"sops_file data source. Encryption keys are resolved via the `.sops.yaml` " +
			"creation rules found by searching destination_file's directory and its parents, no keys are " +
			"configured on this resource directly. Provide the content to encrypt via `data_wo`.",
		Attributes: map[string]schema.Attribute{
			"destination_file": schema.StringAttribute{
				Description: "Path to write the encrypted file to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"output_type": schema.StringAttribute{
				Description: "The provider will use the file extension of destination_file to determine how to " +
					"format the content before encrypting it. If destination_file does not have the usual " +
					"extension, set this argument to `yaml` or `json` accordingly.",
				Optional: true,
			},
			"data_wo": schema.DynamicAttribute{
				Description: "The value to encrypt - an object, list, or scalar - marshalled according to " +
					"output_type. Can be copied directly from any other Terraform value or resource attribute. " +
					"Never stored in plan or state artifacts; must be used together with data_wo_version, since " +
					"Terraform cannot detect changes to a write-only value on its own.",
				Required:  true,
				Sensitive: true,
				WriteOnly: true,
			},
			"data_wo_version": schema.StringAttribute{
				Description: "An arbitrary value that, when changed, tells Terraform that data_wo has changed and " +
					"the resource should be updated.",
				Required: true,
			},
			"id": schema.StringAttribute{
				Description: "Unique identifier for this resource, equal to destination_file.",
				Computed:    true,
			},
		},
	}
}

func (r *fileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan fileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config fileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := writeResourceContent(ctx, plan, config.DataWo); err != nil {
		addFileResourceError(&resp.Diagnostics, err)
		return
	}

	plan.Id = plan.DestinationFile
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *fileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state fileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := os.Stat(state.DestinationFile.ValueString()); err != nil {
		if os.IsNotExist(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to stat file", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *fileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan fileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config fileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := writeResourceContent(ctx, plan, config.DataWo); err != nil {
		addFileResourceError(&resp.Diagnostics, err)
		return
	}

	plan.Id = plan.DestinationFile
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *fileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state fileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := os.Remove(state.DestinationFile.ValueString()); err != nil && !os.IsNotExist(err) {
		resp.Diagnostics.AddError("Failed to delete file", err.Error())
	}
}

// writeResourceContent encrypts dataWo, passed separately since write-only values are only
// available via config, not plan or state.
func writeResourceContent(ctx context.Context, plan fileResourceModel, dataWo types.Dynamic) error {
	value, err := dynamicToGo(ctx, dataWo)
	if err != nil {
		return err
	}
	return writeEncryptedFileFromDynamic(plan.DestinationFile, plan.OutputType, value)
}

func addFileResourceError(diags *diag.Diagnostics, err error) {
	if detailedErr, ok := err.(summaryError); ok {
		diags.AddError(detailedErr.Summary, detailedErr.Err.Error())
	} else {
		diags.AddError("Failed to encrypt content", err.Error())
	}
}
