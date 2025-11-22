package sops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &externalDataSource{}

func newExternalDataSource() datasource.DataSource {
	return &externalDataSource{}
}

type externalDataSource struct{}

type externalDataSourceModel struct {
	InputType types.String `tfsdk:"input_type"`
	Source    types.String `tfsdk:"source"`
	Data      types.Map    `tfsdk:"data"`
	Raw       types.String `tfsdk:"raw"`
	Id        types.String `tfsdk:"id"`
}

func (d *externalDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "sops_external"
}

func (d *externalDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read data from a sops-encrypted string. Useful if the data does not reside on disk locally (otherwise use `sops_file`).",
		Attributes: map[string]schema.Attribute{
			"input_type": schema.StringAttribute{
				Description: "`yaml`, `json` `dotenv` (`.env`), `ini` or `raw`, depending on the structure of the un-encrypted data.",
				Optional:    true,
			},
			"source": schema.StringAttribute{
				Description: "A string with sops-encrypted data",
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

func (d *externalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config externalDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, raw, err := getExternalData(config.Source, config.InputType)
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
