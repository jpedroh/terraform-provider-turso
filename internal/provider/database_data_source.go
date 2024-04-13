package provider

import (
	"context"
	"fmt"
	"terraform-provider-turso/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DatabaseDataSource{}

func NewDatabaseDataSource() datasource.DataSource {
	return &DatabaseDataSource{}
}

// DatabaseDataSource defines the data source implementation.
type DatabaseDataSource struct {
	client *client.Client
}

// DatabaseDataSourceModel describes the data source data model.
type DatabaseDataSourceModel struct {
	OrganizationName types.String `tfsdk:"organization_name"`
	Name             types.String `tfsdk:"name"`
	Group            types.String `tfsdk:"group"`
	SizeLimit        types.String `tfsdk:"size_limit"`
	IsSchema         types.Bool   `tfsdk:"is_schema"`
	Schema           types.String `tfsdk:"schema"`

	// Computed
	DbId     types.String `tfsdk:"db_id"`
	Hostname types.String `tfsdk:"hostname"`
}

func (d *DatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *DatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Database data source",

		Attributes: map[string]schema.Attribute{
			"organization_name": schema.StringAttribute{
				MarkdownDescription: "Name of organization to create the database for",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the new database. Must contain only lowercase letters, numbers, dashes. No longer than 32 characters.",
				Required:            true,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "The name of the group where the database should be created. The group must already exist.",
				Optional:            true,
			},
			"size_limit": schema.StringAttribute{
				MarkdownDescription: "The maximum size of the database in bytes. Values with units are also accepted, e.g. 1mb, 256mb, 1gb.",
				Optional:            true,
			},
			"is_schema": schema.BoolAttribute{
				MarkdownDescription: "Mark this database as the parent schema database that updates child databases with any schema changes.",
				Optional:            true,
			},
			"schema": schema.StringAttribute{
				MarkdownDescription: "The name of the parent database to use as the schema.",
				Optional:            true,
			},
			"db_id": schema.StringAttribute{
				MarkdownDescription: "The database universal unique identifier (UUID).",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The DNS hostname used for client libSQL and HTTP connections.",
				Computed:            true,
			},
		},
	}
}

func (d *DatabaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatabaseDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.V1OrganizationsOrganizationNameDatabasesDatabaseNameGet(ctx, client.V1OrganizationsOrganizationNameDatabasesDatabaseNameGetParams{
		OrganizationName: data.OrganizationName.ValueString(),
		DatabaseName:     data.Name.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err.Error()))
		return
	}

	switch p := res.(type) {
	case *client.V1OrganizationsOrganizationNameDatabasesDatabaseNameGetOK:
		data.DbId = types.StringValue(string(p.Database.Value.DbId.Value))         //nolint:all
		data.Hostname = types.StringValue(string(p.Database.Value.Hostname.Value)) //nolint:all
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
