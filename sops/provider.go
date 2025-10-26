package sops

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &SopsProvider{}

type SopsProvider struct{}

type ProviderConfig struct {
	Env types.Map `tfsdk:"env"`
}

func New() provider.Provider {
	return &SopsProvider{}
}

func (p *SopsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sops"
}

func (p *SopsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"env": schema.MapAttribute{
				Description: "Environment variables to use",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (p *SopsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg ProviderConfig
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !cfg.Env.IsNull() {
		env := make(map[string]types.String, len(cfg.Env.Elements()))
		diags := cfg.Env.ElementsAs(ctx, &env, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for k, v := range env {
			if err := os.Setenv(k, v.ValueString()); err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error setting environment variable %q", k),
					err.Error(),
				)
			}
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

func (p *SopsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newFileDataSource,
		newExternalDataSource,
	}
}

func (p *SopsProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *SopsProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		newFileEphemeralResource,
		newExternalEphemeral,
	}
}
