package sops

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"

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
	Id         types.String `tfsdk:"id"`
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
			"id": schema.StringAttribute{
				Description: "Unique identifier for this data source",
				Computed:    true,
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

	sourceFile := config.SourceFile.ValueString()
	content, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		resp.Diagnostics.AddError("Error reading file", err.Error())
		return
	}

	var format string
	if !config.InputType.IsNull() {
		format = config.InputType.ValueString()
	} else {
		switch ext := path.Ext(sourceFile); ext {
		case ".json":
			format = "json"
		case ".yaml", ".yml":
			format = "yaml"
		case ".env":
			format = "dotenv"
		case ".ini":
			format = "ini"
		default:
			resp.Diagnostics.AddError(
				"Unknown file type",
				fmt.Sprintf("Don't know how to decode file with extension %s, set input_type as appropriate", ext),
			)
			return
		}
	}

	if err := validateInputType(format); err != nil {
		resp.Diagnostics.AddError("Invalid input type", err.Error())
		return
	}

	data, raw, err := readData(content, format)
	if err != nil {
		resp.Diagnostics.AddError("Error reading data", err.Error())
		return
	}

	m, mapDiags := types.MapValueFrom(ctx, types.StringType, data)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Data = m
	config.Raw = types.StringValue(raw)
	config.Id = types.StringValue("-")

	diags = resp.Result.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
