package sops

import (
	"context"
	"time"

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
	InputType        types.String `tfsdk:"input_type"`
	SourceFile       types.String `tfsdk:"source_file"`
	Data             types.Map    `tfsdk:"data"`
	Raw              types.String `tfsdk:"raw"`
	LastModified     types.String `tfsdk:"last_modified"`
	LastModifiedUnix types.Int64  `tfsdk:"last_modified_unix"`
}

func (d *fileEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "sops_file"
}

func (d *fileEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Decrypt sops-encrypted files",
		Attributes: map[string]schema.Attribute{
			"input_type": schema.StringAttribute{
				Description: "The provider will use the file extension to determine how to unmarshal the data. If your file " +
					"does not have the usual extension, set this argument to `yaml`, `json`, `dotenv` (`.env`), `ini` accordingly, " +
					"or `raw` if the encrypted data is encoded differently.",
				Optional: true,
			},
			"source_file": schema.StringAttribute{
				Description: "Path to the encrypted file",
				Required:    true,
			},

			"data": schema.MapAttribute{
				Description: "The unmarshalled data as a dictionary. Use dot-separated keys to access nested data.",
				Computed:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
			"raw": schema.StringAttribute{
				Description: "Raw decrypted content",
				Computed:    true,
				Sensitive:   true,
			},
			"last_modified": schema.StringAttribute{
				Description: "The `lastmodified` timestamp recorded in the sops metadata, in RFC3339 format. " +
					"As an ephemeral resource attribute it is itself ephemeral, so — unlike the `sops_file`/`sops_external` " +
					"data source's `last_modified` — it cannot be assigned to a state-persisted `wo_version` argument.",
				Computed: true,
			},
			"last_modified_unix": schema.Int64Attribute{
				Description: "The `lastmodified` timestamp recorded in the sops metadata, as a Unix epoch (seconds). " +
					"Like `last_modified`, this is an ephemeral attribute and so cannot be assigned to a state-persisted " +
					"`wo_version` argument.",
				Computed: true,
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

	result, err := getFileData(config.SourceFile, config.InputType)
	if err != nil {
		if detailedErr, ok := err.(summaryError); ok {
			resp.Diagnostics.AddError(detailedErr.Summary, detailedErr.Err.Error())
		} else {
			resp.Diagnostics.AddError("Failed to decrypt file", err.Error())
		}
		return
	}

	m, mapDiags := types.MapValueFrom(ctx, types.StringType, result.data)
	resp.Diagnostics.Append(mapDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	config.Data = m
	config.Raw = types.StringValue(result.raw)
	config.LastModified = types.StringValue(result.lastModified.Format(time.RFC3339))
	config.LastModifiedUnix = types.Int64Value(result.lastModified.Unix())

	diags = resp.Result.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
