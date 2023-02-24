// Package org provides a Terraform datasource.
package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &orgGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &orgGroupDataSource{}
)

func NewOrgGroupDataSource() datasource.DataSource {
	return &orgGroupDataSource{}
}

type orgGroupDataSource struct {
	client *client.CloudAvenue
}

type orgGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Role        types.String `tfsdk:"role"`
	UserNames   types.List   `tfsdk:"user_names"`
}

func (d *orgGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "group"
}

func (d *orgGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a CloudAvenue Org User data source. This can be used to read group.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID is a group `name`.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A name for the org group",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the org group",
			},
			"role": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The role to assign to the org group",
			},
			"user_names": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Set of user names that belong to the org group",
			},
		},
	}
}

func (d *orgGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *orgGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data orgGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// group creation is accessible only for administator account
	adminOrg, err := d.client.Vmware.GetAdminOrgByNameOrId(d.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	orgGroup, err := adminOrg.GetGroupByName(data.Name.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org Group", err.Error())
		return
	}

	var userNames []attr.Value
	for _, user := range orgGroup.Group.UsersList.UserReference {
		userNames = append(userNames, types.StringValue(user.Name))
	}

	data = orgGroupDataSourceModel{
		ID:          types.StringValue(orgGroup.Group.ID),
		Name:        types.StringValue(orgGroup.Group.Name),
		Description: types.StringValue(orgGroup.Group.Description),
		Role:        types.StringValue(orgGroup.Group.Role.Name),
		UserNames:   types.ListValueMust(types.StringType, userNames),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
