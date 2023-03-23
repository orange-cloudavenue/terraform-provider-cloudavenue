// Package org provides a Terraform datasource.
package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
)

func NewRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

type roleDataSource struct {
	client   *client.CloudAvenue
	adminOrg *govcd.AdminOrg
}

func (d *roleDataSource) Init(_ context.Context, rm *roleDataSourceModel) (diags diag.Diagnostics) {
	var err error

	d.adminOrg, err = d.client.Vmware.GetAdminOrgByNameOrId(d.client.GetOrgName())
	if err != nil {
		diags.AddError("[role create] Error retrieving Org", err.Error())
		return
	}

	return
}

func (d *roleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "role"
}

func (d *roleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = roleSchema().GetDataSource()
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
	var data *roleDataSourceModel

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
	data, err := getRole(d.adminOrg, data.Name, data.ID)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[role read] Error retrieving role", err.Error())
		return
	}

	// Save data into Terraform data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getRole(adminOrg *govcd.AdminOrg, name, id types.String) (r *roleDataSourceModel, err error) {
	var role *govcd.Role

	// Get the role
	if id.IsNull() {
		role, err = adminOrg.GetRoleByName(name.ValueString())
	} else {
		role, err = adminOrg.GetRoleById(id.ValueString())
	}
	if err != nil {
		return
	}

	// Get rights
	rights, err := role.GetRights(nil)
	if err != nil {
		return
	}
	assignedRights := []attr.Value{}
	for _, right := range rights {
		assignedRights = append(assignedRights, types.StringValue(right.Name))
	}

	r = &roleDataSourceModel{
		ID:          types.StringValue(role.Role.ID),
		Name:        types.StringValue(role.Role.Name),
		ReadOnly:    types.BoolValue(role.Role.ReadOnly),
		Description: types.StringValue(role.Role.Description),
		Rights:      types.SetNull(types.StringType),
	}

	var y diag.Diagnostics
	if len(assignedRights) > 0 {
		r.Rights, y = types.SetValue(types.StringType, assignedRights)
		if y.HasError() {
			return nil, errors.New("unable to set rights value")
		}
	}

	return r, nil
}
