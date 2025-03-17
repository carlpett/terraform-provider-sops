package sops

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &fileDataSource{}

func newFileDataSource() datasource.DataSource {
	return &fileDataSource{}
}

type fileDataSource struct{}

type fileDataSourceModel struct {
	InputType  types.String `tfsdk:"input_type"`
	SourceFile types.String `tfsdk:"source_file"`
	Data       types.Map    `tfsdk:"data"`
	Raw        types.String `tfsdk:"raw"`
	Id         types.String `tfsdk:"id"`
}

func (d *fileDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "sops_file"
}

func (d *fileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Decrypt sops-encrypted files",
		Attributes: map[string]schema.Attribute{
			"input_type": schema.StringAttribute{
				Description: "Type of the input file (json, yaml, raw)",
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

func (d *fileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config fileDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceFile := config.SourceFile.ValueString()
	content, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading File", err.Error())
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
				"Unknown File Type",
				fmt.Sprintf("Don't know how to decode file with extension %s, set input_type to json, yaml or raw as appropriate", ext),
			)
			return
		}
	}

	if err := validateInputType(format); err != nil {
		resp.Diagnostics.AddError("Invalid Input Type", err.Error())
		return
	}

	data, raw, err := readData(content, format)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Data", err.Error())
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

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
