// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-turso/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DatabaseResource{}
var _ resource.ResourceWithImportState = &DatabaseResource{}

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}

type DatabaseResource struct {
	client *client.Client
}

type DatabaseResourceModel struct {
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

func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Database resource",

		Attributes: map[string]schema.Attribute{
			"organization_name": schema.StringAttribute{
				MarkdownDescription: "Name of organization to create the database for",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the new database. Must contain only lowercase letters, numbers, dashes. No longer than 32 characters.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "The name of the group where the database should be created. The group must already exist.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The DNS hostname used for client libSQL and HTTP connections.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.CreateDatabase(ctx, &client.CreateDatabaseInput{
		Name:      data.Name.ValueString(),
		Group:     data.Group.ValueString(),
		SizeLimit: client.NewOptString(data.SizeLimit.ValueString()),
		Schema:    client.NewOptString(data.Schema.ValueString()),
		IsSchema:  client.NewOptBool(data.IsSchema.ValueBool()),
	}, client.CreateDatabaseParams{
		OrganizationSlug: data.OrganizationName.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create database, got error: %s", err.Error()))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.

	switch p := res.(type) {
	case *client.CreateDatabaseOK:
		data.DbId = types.StringValue(string(p.Database.Value.DbId.Value))
		data.Hostname = types.StringValue(string(p.Database.Value.Hostname.Value))
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created database resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetDatabase(ctx, client.GetDatabaseParams{
		OrganizationSlug: data.OrganizationName.ValueString(),
		DatabaseName:     data.Name.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err.Error()))
		return
	}

	switch p := res.(type) {
	case *client.GetDatabaseOK:
		data.DbId = types.StringValue(string(p.Database.Value.DbId.Value))         //nolint:all
		data.Hostname = types.StringValue(string(p.Database.Value.Hostname.Value)) //nolint:all
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateDatabaseConfiguration(ctx, &client.DatabaseConfigurationInput{
		SizeLimit: client.NewOptString(data.SizeLimit.ValueString()),
	}, client.UpdateDatabaseConfigurationParams{
		OrganizationSlug: data.OrganizationName.ValueString(),
		DatabaseName:     data.Name.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update database, got error: %s", err.Error()))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteDatabase(ctx, client.DeleteDatabaseParams{
		OrganizationSlug: data.OrganizationName.ValueString(),
		DatabaseName:     data.Name.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete database, got error: %s", err.Error()))
		return
	}
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization/name. Got: %q", req.ID),
		)
		return
	}

	organization_name := idParts[0]
	name := idParts[1]

	tflog.Debug(ctx, fmt.Sprintf("Importing database %s/%s", organization_name, name))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_name"), organization_name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)

	res, err := r.client.GetDatabase(ctx, client.GetDatabaseParams{
		OrganizationSlug: organization_name,
		DatabaseName:     name,
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err.Error()))
		return
	}

	switch p := res.(type) {
	case *client.GetDatabaseOK:
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group"), types.StringValue(p.Database.Value.Group.Value))...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("db_id"), types.StringValue(p.Database.Value.DbId.Value))...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("hostname"), types.StringValue(p.Database.Value.Hostname.Value))...)
	}
}
