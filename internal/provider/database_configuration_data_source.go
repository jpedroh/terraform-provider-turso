package provider

import (
	"context"
	"fmt"
	"terraform-provider-turso/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DatabaseConfigurationDataSource{}

func NewDatabaseConfigurationDataSource() datasource.DataSource {
	return &DatabaseConfigurationDataSource{}
}

type DatabaseConfigurationDataSource struct {
	client *client.Client
}

type DatabaseConfigurationDataSourceModel struct {
	OrganizationSlug types.String `tfsdk:"organization_slug"`
	DatabaseName     types.String `tfsdk:"database_name"`
	SizeLimit        types.String `tfsdk:"size_limit"`
	BlockReads       types.Bool   `tfsdk:"block_reads"`
	BlockWrites      types.Bool   `tfsdk:"block_writes"`
	DeleteProtection types.Bool   `tfsdk:"delete_protection"`
}

func (d *DatabaseConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_configuration"
}

func (d *DatabaseConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Database data source",

		Attributes: map[string]schema.Attribute{
			"organization_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization or user account.",
				Required:            true,
			},
			"database_name": schema.StringAttribute{
				MarkdownDescription: "The name of the database.",
				Required:            true,
			},
			"size_limit": schema.StringAttribute{
				MarkdownDescription: "The maximum size of the database in bytes. Values with units are also accepted, e.g. 1mb, 256mb, 1gb.",
				Computed:            true,
			},
			"block_reads": schema.BoolAttribute{
				MarkdownDescription: "Block all database reads.",
				Computed:            true,
			},
			"block_writes": schema.BoolAttribute{
				MarkdownDescription: "Block all database writes.",
				Computed:            true,
			},
			"delete_protection": schema.BoolAttribute{
				MarkdownDescription: "Prevent the database from being deleted.",
				Computed:            true,
			},
		},
	}
}

func (d *DatabaseConfigurationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DatabaseConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatabaseConfigurationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.GetDatabaseConfiguration(ctx, client.GetDatabaseConfigurationParams{
		OrganizationSlug: data.OrganizationSlug.ValueString(),
		DatabaseName:     data.DatabaseName.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Unable to read database configuration", err.Error())
		return
	}

	data.BlockReads = types.BoolValue(res.BlockReads.Value)
	data.BlockWrites = types.BoolValue(res.BlockWrites.Value)
	data.DeleteProtection = types.BoolValue(res.DeleteProtection.Value)
	data.SizeLimit = types.StringValue(res.SizeLimit.Value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
