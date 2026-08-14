package sops

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &fileMetadataDataSource{}

func newFileMetadataDataSource() datasource.DataSource {
	return &fileMetadataDataSource{}
}

type fileMetadataDataSource struct{}

type fileMetadataDataSourceModel struct {
	InputType        types.String `tfsdk:"input_type"`
	SourceFile       types.String `tfsdk:"source_file"`
	LastModified     types.String `tfsdk:"last_modified"`
	LastModifiedUnix types.Int64  `tfsdk:"last_modified_unix"`
	Id               types.String `tfsdk:"id"`
}

func (d *fileMetadataDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "sops_file_metadata"
}

func (d *fileMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Read the metadata of a sops-encrypted file on disk without decrypting it. Only the unencrypted " +
			"sops metadata block is parsed, so no data key access (KMS, age, PGP, …) is required and no secret values are " +
			"ever read or written to state. Useful for driving a write-only argument's `wo_version` from the `lastmodified` " +
			"timestamp while the secrets themselves are read separately (for example via the `sops_file` ephemeral resource).",
		Attributes: map[string]schema.Attribute{
			"input_type": schema.StringAttribute{
				Description: "The provider will use the file extension to determine how to read the metadata. If your file " +
					"does not have the usual extension, set this argument to `yaml`, `json`, `dotenv` (`.env`), `ini` accordingly, " +
					"or `raw` if the encrypted data is encoded differently.",
				Optional: true,
			},
			"source_file": schema.StringAttribute{
				Description: "Path to the encrypted file.",
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
				Description: "Unique identifier for this data source.",
				Computed:    true,
			},
		},
	}
}

func (d *fileMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config fileMetadataDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lastModified, err := getFileMetadata(config.SourceFile, config.InputType)
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
