// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure WandbLaunchProvider satisfies various provider interfaces.
var _ provider.Provider = &WandbLaunchProvider{}
var _ provider.ProviderWithFunctions = &WandbLaunchProvider{}

// WandbLaunchProvider defines the provider implementation.
type WandbLaunchProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// WandbLaunchProviderModel describes the provider data model.
type WandbLaunchProviderModel struct {
	BaseUrl types.String `tfsdk:"base_url"`
	ApiKey  types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WandbLaunchProvider{
			version: version,
		}
	}
}

func (p *WandbLaunchProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "wandb"
	resp.Version = p.version
}

func (p *WandbLaunchProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "The base URL of the W&B API. Defaults to WANDB_BASE_URL environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The API key for the W&B API. Defaults to WANDB_API_KEY environment variable.",
			},
		},
	}
}

func (p *WandbLaunchProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var config WandbLaunchProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseUrl := os.Getenv("WANDB_BASE_URL")
	apiKey := os.Getenv("WANDB_API_KEY")

	if !config.BaseUrl.IsNull() {
		baseUrl = config.BaseUrl.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if baseUrl == "" {
		resp.Diagnostics.AddError(
			"Missing W&B API base URL",
			"Set WANDB_BASE_URL environment variable or configure the base_url.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing W&B API key",
			"Set WANDB_API_KEY environment variable or configure the api_key.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !strings.HasSuffix(baseUrl, "/graphql") {
		baseUrl += "/graphql"
	}
	headers := http.Header{
		"User-Agent":    []string{"terraform-provider-wandb-launch"},
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte("api:"+apiKey))},
		"Content-Type":  []string{"application/json"},
	}
	client := NewGraphQLClientWithHeaders(baseUrl, headers)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *WandbLaunchProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRunQueueResource,
	}
}

func (p *WandbLaunchProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *WandbLaunchProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}
