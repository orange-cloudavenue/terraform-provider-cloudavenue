// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2024-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// Package iam provides a Terraform datasource.
package iam

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &RolesDataSource{}
	_ datasource.DataSourceWithConfigure = &RolesDataSource{}
)

func NewRolesDataSource() datasource.DataSource {
	return &RolesDataSource{}
}

type RolesDataSource struct {
	client     *client.CloudAvenue
	roleClient roleClientV2
}

// Init Initializes the data source.
func (d *RolesDataSource) Init(_ context.Context, _ *RolesModel) (diags diag.Diagnostics) {
	d.roleClient, diags = InitRoleClientV2(d.client)
	return diags
}

func (d *RolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_roles"
}

func (d *RolesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rolesSchema(ctx).GetDataSource(ctx)
}

func (d *RolesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_iam_roles", d.client.GetOrgName(), metrics.Read)()

	config := &RolesModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all roles
	rolesList, err := d.roleClient.ListOrgRoles(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get roles", err.Error())
		return
	}

	roles, di := config.Roles.Get(ctx)
	if di.HasError() {
		resp.Diagnostics.Append(di...)
		return
	}

	idOfRoles := make([]string, 0)

	for _, role := range rolesList {
		rights, err := d.roleClient.ListOrgRoleRights(ctx, role.ID)
		if err != nil {
			continue
		}

		idOfRoles = append(idOfRoles, role.ID)

		roles[role.Name] = &RoleDataSourceModel{
			ID:          supertypes.NewStringNull(),
			Name:        supertypes.NewStringNull(),
			Description: supertypes.NewStringNull(),
			ReadOnly:    supertypes.NewBoolNull(),
			Rights:      supertypes.NewSetValueOfNull[string](ctx),
		}
		roles[role.Name].ID.Set(role.ID)
		roles[role.Name].Name.Set(role.Name)
		roles[role.Name].Description.Set(role.Description)
		roles[role.Name].ReadOnly.Set(role.ReadOnly)
		roles[role.Name].Rights.Set(ctx, func() []string {
			var r []string

			for _, right := range rights {
				r = append(r, right.Name)
			}

			return r
		}())
	}

	config.ID.Set(utils.GenerateUUID(idOfRoles...).String())
	resp.Diagnostics.Append(config.Roles.Set(ctx, roles)...)

	// Save data into Terraform data
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
