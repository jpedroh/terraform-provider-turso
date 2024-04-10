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

var _ resource.Resource = &DatabaseTokenResource{}
var _ resource.ResourceWithImportState = &DatabaseTokenResource{}

func NewDatabaseTokenResource() resource.Resource {
	return &DatabaseTokenResource{}
}

type DatabaseTokenResource struct {
	client *client.Client
}

type DatabaseTokenResourceModel struct {
	OrganizationName types.String `tfsdk:"organization_name"`
	DatabaseName     types.String `tfsdk:"database_name"`
	Expiration       types.String `tfsdk:"expiration"`
	Authorization    types.String `tfsdk:"authorization"`

	JWT types.String `tfsdk:"jwt"`
}

func (r *DatabaseTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_token"
}

func (r *DatabaseTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Database Token resource",

		Attributes: map[string]schema.Attribute{
			"organization_name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization or user.",
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
			"expiration": schema.StringAttribute{
				MarkdownDescription: "Expiration time for the token (e.g., 2w1d30m).",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"authorization": schema.StringAttribute{
				MarkdownDescription: "Authorization level for the token (full-access or read-only).",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"jwt": schema.StringAttribute{
				MarkdownDescription: "The generated authorization token (JWT).",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DatabaseTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseTokenResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.V1OrganizationsOrganizationNameDatabasesDatabaseNameAuthTokensPost(ctx, client.OptCreateTokenInput{}, client.V1OrganizationsOrganizationNameDatabasesDatabaseNameAuthTokensPostParams{
		OrganizationName: data.OrganizationName.ValueString(),
		DatabaseName:     data.DatabaseName.ValueString(),
		Expiration:       client.NewOptString(data.Expiration.ValueString()),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create database, got error: %s", err.Error()))
		return
	}

	switch p := res.(type) {
	case *client.V1OrganizationsOrganizationNameDatabasesDatabaseNameAuthTokensPostOK:
		data.JWT = types.StringValue(p.Jwt.Value)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created database resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseTokenResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Currently, it's not possible to read a token

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseTokenResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Currently, it's not possible to update a token

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabaseTokenResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Currently it's not possible to granularly revoke a token.
}

func (r *DatabaseTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO: Currently, it's not possible to import a token.
}
