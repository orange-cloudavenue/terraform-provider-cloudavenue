// Package org provides a Terraform datasource.
package iam

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

var (
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
	_ role                               = &roleDataSource{}
)

func NewRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

type roleDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	role     commonRole
}

func (d *roleDataSource) Init(_ context.Context, rm *RoleDataSourceModel) (diags diag.Diagnostics) {
	d.role = commonRole{
		ID:   rm.ID.StringValue,
		Name: rm.Name.StringValue,
	}

	d.adminOrg, diags = adminorg.Init(d.client)

	return
}

func (d *roleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "role"
}

func (d *roleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = roleSchema().GetDataSource(ctx)
}

func (d *roleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_iam_role", d.client.GetOrgName(), metrics.Read)()

	var data *RoleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Role
	role, err := d.GetRole()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving role", err.Error())
		return
	}

	// Get rights
	rights, err := role.GetRights(nil)
	if err != nil {
		return
	}
	assignedRights := []string{}
	for _, right := range rights {
		assignedRights = append(assignedRights, right.Name)
	}

	resp.Diagnostics.Append(data.Rights.Set(ctx, assignedRights)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID.Set(role.Role.ID)
	data.Name.Set(role.Role.Name)
	data.ReadOnly.Set(role.Role.ReadOnly)
	data.Description.Set(role.Role.Description)

	// Save data into Terraform data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *roleDataSource) GetRole() (*govcd.Role, error) {
	return d.role.GetRole(d.adminOrg)
}
