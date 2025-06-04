package sops

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &externalEphemeralResource{}

func newExternalEphemeral() ephemeral.EphemeralResource {
	return &externalEphemeralResource{}
}

type externalEphemeralResource struct{}

type externalEphemeralModel struct {
	InputType types.String `tfsdk:"input_type"`
	Source    types.String `tfsdk:"source"`
	Data      types.Map    `tfsdk:"data"`
	Raw       types.String `tfsdk:"raw"`
	Id        types.String `tfsdk:"id"`
}

func (d *externalEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "sops_external"
}

func (d *externalEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Decrypt sops-encrypted data from external commands",
		Attributes: map[string]schema.Attribute{
			"input_type": schema.StringAttribute{
				Description: "Type of the input data: json, yaml, dotenv, ini, raw",
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

func (d *externalEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var config externalEphemeralModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	source := config.Source.ValueString()
	content, err := ioutil.ReadAll(strings.NewReader(source))
	if err != nil {
		resp.Diagnostics.AddError("Error reading source", err.Error())
		return
	}

	format := config.InputType.ValueString()
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
