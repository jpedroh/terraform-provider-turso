// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"terraform-provider-turso/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure TursoProvider satisfies various provider interfaces.
var _ provider.Provider = &TursoProvider{}
var _ provider.ProviderWithFunctions = &TursoProvider{}

// TursoProvider defines the provider implementation.
type TursoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TursoProviderModel describes the provider data model.
type TursoProviderModel struct {
	ApiToken types.String `tfsdk:"api_token"`
}

func (p *TursoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "turso"
	resp.Version = p.version
}

func (p *TursoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token to authenticate with Turso API",
				Required:            true,
			},
		},
	}
}

type AuthenticationRoundTripper struct {
	Token   string
	Proxied http.RoundTripper
	Diag    diag.Diagnostics
}

func (lrt AuthenticationRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", lrt.Token))
	return lrt.Proxied.RoundTrip(req)
}

func (p *TursoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config TursoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Example client configuration for data sources and resources
	httpClient := &http.Client{
		Transport: AuthenticationRoundTripper{
			Token:   config.ApiToken.ValueString(),
			Proxied: http.DefaultTransport,
			Diag:    resp.Diagnostics,
		},
	}

	client, err := client.NewClient("https://api.turso.tech", client.WithClient(httpClient))
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *TursoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDatabaseResource,
		NewDatabaseTokenResource,
	}
}

func (p *TursoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDatabaseDataSource,
		NewOrganizationDataSource,
	}
}

func (p *TursoProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TursoProvider{
			version: version,
		}
	}
}
