package sops

import (
	"context"

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
		Description: "Read data from a sops-encrypted file on disk.",
		Attributes: map[string]schema.Attribute{
			"input_type": schema.StringAttribute{
				Description: "The provider will use the file extension to determine how to unmarshal the data. If your file " +
					"does not have the usual extension, set this argument to `yaml`, `json`, `dotenv` (`.env`), `ini` accordingly, " +
					"or `raw` if the encrypted data is encoded differently.",
				Optional: true,
			},
			"source_file": schema.StringAttribute{
				Description: "Path to the encrypted file.",
				Required:    true,
			},

			"data": schema.MapAttribute{
				Description: "The unmarshalled data as a dictionary. Use dot-separated keys to access nested data.",
				Computed:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
			"raw": schema.StringAttribute{
				Description: "The entire unencrypted file as a string.",
				Computed:    true,
				Sensitive:   true,
			},
			"id": schema.StringAttribute{
				Description: "Unique identifier for this data source.",
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
	config.Id = types.StringValue("-")

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
