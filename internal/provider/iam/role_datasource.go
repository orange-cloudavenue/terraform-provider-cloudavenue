// Package org provides a Terraform datasource.
package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &iamRoleDataSource{}
	_ datasource.DataSourceWithConfigure = &iamRoleDataSource{}
)

func NewIAMRoleDataSource() datasource.DataSource {
	return &iamRoleDataSource{}
}

type iamRoleDataSource struct {
	client *client.CloudAvenue
}

type iamRoleDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	BundleKey   types.String `tfsdk:"bundle_key"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	Rights      types.Set    `tfsdk:"rights"`
}

func (d *iamRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "role"
}

func (d *iamRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The CloudAvenue iam role datasource allows you to read roles.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the role.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A name for the role",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description for the role",
			},
			"bundle_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Key used for internationalization",
			},
			"read_only": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the role is read only",
			},
			"rights": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "A list of rights for the role",
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *iamRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data *iamRoleDataSourceModel
		err  error
		role *govcd.Role
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// role read is accessible only in administrator
	adminOrg, err := d.client.Vmware.GetAdminOrgByNameOrId(d.client.GetOrgName())
	if err != nil {
		resp.Diagnostics.AddError("[role read] Error retrieving Org", err.Error())
		return
	}

	// Get the role
	roleID := data.ID.ValueString()
	roleName := data.Name.ValueString()
	if roleID == "" {
		role, err = adminOrg.GetRoleByName(roleName)
	} else {
		role, err = adminOrg.GetRoleById(roleID)
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[role read] Error retrieving role", err.Error())
		return
	}

	// Get rights
	rights, err := role.GetRights(nil)
	if err != nil {
		resp.Diagnostics.AddError("[role read] Error while querying role rights", err.Error())
		return
	}
	assignedRights := []attr.Value{}
	for _, right := range rights {
		assignedRights = append(assignedRights, types.StringValue(right.Name))
	}
	var y diag.Diagnostics
	if len(assignedRights) > 0 {
		data.Rights, y = types.SetValue(types.StringType, assignedRights)
		resp.Diagnostics.Append(y...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set state to fully populated data
	data = &iamRoleDataSourceModel{
		ID:          types.StringValue(role.Role.ID),
		Name:        types.StringValue(role.Role.Name),
		BundleKey:   types.StringValue(role.Role.BundleKey),
		ReadOnly:    types.BoolValue(role.Role.ReadOnly),
		Description: types.StringValue(role.Role.Description),
		Rights:      data.Rights,
	}

	// Save data into Terraform data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
