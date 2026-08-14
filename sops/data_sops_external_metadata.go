package sops

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &externalMetadataDataSource{}

func newExternalMetadataDataSource() datasource.DataSource {
	return &externalMetadataDataSource{}
}

type externalMetadataDataSource struct{}

type externalMetadataDataSourceModel struct {
	InputType        types.String `tfsdk:"input_type"`
	Source           types.String `tfsdk:"source"`
	LastModified     types.String `tfsdk:"last_modified"`
	LastModifiedUnix types.Int64  `tfsdk:"last_modified_unix"`
	Id               types.String `tfsdk:"id"`
}

func (d *externalMetadataDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "sops_external_metadata"
}

func (d *externalMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read the metadata of a sops-encrypted string without decrypting it. Only the unencrypted sops " +
			"metadata block is parsed, so no data key access (KMS, age, PGP, …) is required and no secret values are ever " +
			"read or written to state. Useful when the encrypted data does not reside on disk locally (otherwise use " +
			"`sops_file_metadata`).",
		Attributes: map[string]schema.Attribute{
			"input_type": schema.StringAttribute{
				Description: "`yaml`, `json` `dotenv` (`.env`), `ini` or `raw`, depending on the structure of the un-encrypted data.",
				Optional:    true,
			},
			"source": schema.StringAttribute{
				Description: "A string with sops-encrypted data",
				Required:    true,
			},

			"last_modified": schema.StringAttribute{
				Description: "The `lastmodified` timestamp recorded in the sops metadata, in RFC3339 format. " +
					"Useful as a version identifier, for example as the `wo_version` of a write-only argument " +
					"(see also `last_modified_unix`).",
				Computed: true,
			},
			"last_modified_unix": schema.Int64Attribute{
				Description: "The `lastmodified` timestamp recorded in the sops metadata, as a Unix epoch (seconds). " +
					"Directly usable as the integer `wo_version` of a write-only argument, without needing to parse " +
					"the RFC3339 `last_modified` value.",
				Computed: true,
			},
			"id": schema.StringAttribute{
				Description: "Unique identifier for this data source",
				Computed:    true,
			},
		},
	}
}

func (d *externalMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config externalMetadataDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lastModified, err := getExternalMetadata(config.Source, config.InputType)
	if err != nil {
		if detailedErr, ok := err.(summaryError); ok {
			resp.Diagnostics.AddError(detailedErr.Summary, detailedErr.Err.Error())
		} else {
			resp.Diagnostics.AddError("Failed to read metadata", err.Error())
		}
		return
	}

	config.LastModified = types.StringValue(lastModified.Format(time.RFC3339))
	config.LastModifiedUnix = types.Int64Value(lastModified.Unix())
	config.Id = types.StringValue("-")

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
