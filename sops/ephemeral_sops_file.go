package sops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &fileEphemeralResource{}

func newFileEphemeralResource() ephemeral.EphemeralResource {
	return &fileEphemeralResource{}
}

type fileEphemeralResource struct{}

type fileEphemeralResourceModel struct {
	InputType  types.String `tfsdk:"input_type"`
	SourceFile types.String `tfsdk:"source_file"`
	Data       types.Map    `tfsdk:"data"`
	Raw        types.String `tfsdk:"raw"`
}

func (d *fileEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "sops_file"
}

func (d *fileEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Decrypt sops-encrypted files",
		Attributes: map[string]schema.Attribute{
			"input_type": schema.StringAttribute{
				Description: "Type of the input file: json, yaml, dotenv, ini, raw",
				Optional:    true,
			},
			"source_file": schema.StringAttribute{
				Description: "Path to the encrypted file",
				Required:    true,
			},

			"data": schema.MapAttribute{
				Description: "Decrypted data",
				Computed:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
			"raw": schema.StringAttribute{
				Description: "Raw decrypted content",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (d *fileEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var config fileEphemeralResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, raw, err := getFileData(config.SourceFile, config.InputType)
	if err != nil {
		if detailedErr, ok := err.(summaryError); ok {
			resp.Diagnostics.AddError(detailedErr.Summary, detailedErr.Err.Error())
		} else {
			resp.Diagnostics.AddError("Failed to decrypt file", err.Error())
		}
		return
	}

	m, mapDiags := types.MapValueFrom(ctx, types.StringType, data)
	resp.Diagnostics.Append(mapDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	config.Data = m
	config.Raw = types.StringValue(raw)

	diags = resp.Result.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
