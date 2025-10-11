// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-turso/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &DatabaseConfigurationResource{}
var _ resource.ResourceWithImportState = &DatabaseConfigurationResource{}

func NewDatabaseConfigurationResource() resource.Resource {
	return &DatabaseConfigurationResource{}
}

type DatabaseConfigurationResource struct {
	client *client.Client
}

type DatabaseConfigurationResourceModel struct {
	OrganizationSlug types.String `tfsdk:"organization_slug"`
	DatabaseName     types.String `tfsdk:"database_name"`
	SizeLimit        types.String `tfsdk:"size_limit"`
	BlockReads       types.Bool   `tfsdk:"block_reads"`
	BlockWrites      types.Bool   `tfsdk:"block_writes"`
	DeleteProtection types.Bool   `tfsdk:"delete_protection"`
}

func (r *DatabaseConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_configuration"
}

func (r *DatabaseConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a database configuration belonging to the organization or user.",

		Attributes: map[string]schema.Attribute{
			"organization_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization or user account.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"database_name": schema.StringAttribute{
				MarkdownDescription: "The name of the database.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size_limit": schema.StringAttribute{
				MarkdownDescription: "The maximum size of the database in bytes. Values with units are also accepted, e.g. 1mb, 256mb, 1gb.",
				Optional:            true,
			},
			"block_reads": schema.BoolAttribute{
				MarkdownDescription: "Block all database reads.",
				Optional:            true,
			},
			"block_writes": schema.BoolAttribute{
				MarkdownDescription: "Block all database writes.",
				Optional:            true,
			},
			"delete_protection": schema.BoolAttribute{
				MarkdownDescription: "Prevent the database from being deleted.",
				Optional:            true,
			},
		},
	}
}

func (r *DatabaseConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *DatabaseConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseConfigurationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.UpdateDatabaseConfiguration(ctx, &client.DatabaseConfigurationInput{
		SizeLimit:        client.NewOptString(data.SizeLimit.ValueString()),
		BlockReads:       client.NewOptBool(data.BlockReads.ValueBool()),
		BlockWrites:      client.NewOptBool(data.BlockWrites.ValueBool()),
		DeleteProtection: client.NewOptBool(data.DeleteProtection.ValueBool()),
	}, client.UpdateDatabaseConfigurationParams{
		OrganizationSlug: data.OrganizationSlug.ValueString(),
		DatabaseName:     data.DatabaseName.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set database configuration, got error: %s", err.Error()))
		return
	}

	data.BlockReads = types.BoolValue(res.BlockReads.Value)
	data.BlockWrites = types.BoolValue(res.BlockWrites.Value)
	data.DeleteProtection = types.BoolValue(res.DeleteProtection.Value)
	data.SizeLimit = types.StringValue(res.SizeLimit.Value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseConfigurationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetDatabaseConfiguration(ctx, client.GetDatabaseConfigurationParams{
		OrganizationSlug: data.OrganizationSlug.ValueString(),
		DatabaseName:     data.DatabaseName.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err.Error()))
		return
	}

	data.BlockReads = types.BoolValue(res.BlockReads.Value)
	data.BlockWrites = types.BoolValue(res.BlockWrites.Value)
	data.DeleteProtection = types.BoolValue(res.DeleteProtection.Value)
	data.SizeLimit = types.StringValue(res.SizeLimit.Value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseConfigurationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.UpdateDatabaseConfiguration(ctx, &client.DatabaseConfigurationInput{
		SizeLimit:        client.NewOptString(data.SizeLimit.ValueString()),
		BlockReads:       client.NewOptBool(data.BlockReads.ValueBool()),
		BlockWrites:      client.NewOptBool(data.BlockWrites.ValueBool()),
		DeleteProtection: client.NewOptBool(data.DeleteProtection.ValueBool()),
	}, client.UpdateDatabaseConfigurationParams{
		OrganizationSlug: data.OrganizationSlug.ValueString(),
		DatabaseName:     data.DatabaseName.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to set database configuration, got error: %s", err.Error()))
		return
	}

	data.BlockReads = types.BoolValue(res.BlockReads.Value)
	data.BlockWrites = types.BoolValue(res.BlockWrites.Value)
	data.DeleteProtection = types.BoolValue(res.DeleteProtection.Value)
	data.SizeLimit = types.StringValue(res.SizeLimit.Value)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No operation, as this resource only manages configuration of an existing database.
	// The database itself must be deleted via the turso_database resource.
}

func (r *DatabaseConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization_name, name, err := ExtractDbIdFromImportStateId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), req.ID)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Importing database configuration %s/%s", organization_name, name))

	res, err := r.client.GetDatabaseConfiguration(ctx, client.GetDatabaseConfigurationParams{
		OrganizationSlug: organization_name,
		DatabaseName:     name,
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, DatabaseConfigurationResourceModel{
		OrganizationSlug: types.StringValue(organization_name),
		DatabaseName:     types.StringValue(name),
		SizeLimit:        types.StringValue(res.SizeLimit.Value),
		BlockReads:       types.BoolValue(res.BlockReads.Value),
		BlockWrites:      types.BoolValue(res.BlockWrites.Value),
		DeleteProtection: types.BoolValue(res.DeleteProtection.Value),
	})...)
}
