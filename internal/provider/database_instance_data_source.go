// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-turso/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DatabaseInstanceDataSource{}

func NewDatabaseInstanceDataSource() datasource.DataSource {
	return &DatabaseInstanceDataSource{}
}

type DatabaseInstanceDataSource struct {
	client *client.Client
}

type DatabaseInstanceDataSourceModel struct {
	OrganizationSlug types.String `tfsdk:"organization_slug"`
	DatabaseName     types.String `tfsdk:"database_name"`
	Name             types.String `tfsdk:"name"`
	UUID             types.String `tfsdk:"uuid"`
	Type             types.String `tfsdk:"type"`
	Region           types.String `tfsdk:"region"`
	Hostname         types.String `tfsdk:"hostname"`
}

func (d *DatabaseInstanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_instance"
}

func (d *DatabaseInstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves data from an individual database instance by name.",

		Attributes: map[string]schema.Attribute{
			"organization_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization or user account.",
				Required:            true,
			},
			"database_name": schema.StringAttribute{
				MarkdownDescription: "The name of the database.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the instance.",
				Required:            true,
			},
			"uuid": schema.StringAttribute{
				MarkdownDescription: "The instance universal unique identifier (UUID).",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of database instance this, will be `primary` or `replica`.",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("primary", "replica"),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The location code for the region this instance is located.",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The DNS hostname used for client libSQL and HTTP connections (specific to this instance only).",
				Computed:            true,
			},
		},
	}
}

func (d *DatabaseInstanceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DatabaseInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatabaseInstanceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.GetDatabaseInstance(ctx, client.GetDatabaseInstanceParams{
		OrganizationSlug: data.OrganizationSlug.ValueString(),
		DatabaseName:     data.DatabaseName.ValueString(),
		InstanceName:     data.Name.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Unable to read database instance", err.Error())
		return
	}

	data.UUID = types.StringValue(res.Instance.Value.UUID.Value)
	data.Type = types.StringValue(string(res.Instance.Value.Type.Value))
	data.Region = types.StringValue(res.Instance.Value.Region.Value)
	data.Hostname = types.StringValue(res.Instance.Value.Hostname.Value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
