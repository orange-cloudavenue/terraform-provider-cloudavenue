// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// Package org provides a Terraform datasource.
package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
)

var (
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
)

func NewRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

type roleDataSource struct {
	client     *client.CloudAvenue
	roleClient roleClientV2
}

func (d *roleDataSource) Init(_ context.Context, _ *RoleDataSourceModel) (diags diag.Diagnostics) {
	d.roleClient, diags = InitRoleClientV2(d.client)
	return diags
}

func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "role"
}

func (d *roleDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = roleSchema().GetDataSource(ctx)
}

func (d *roleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	// Determine lookup key: prefer ID if set, otherwise use Name
	lookupKey := data.ID.ValueString()
	if lookupKey == "" {
		lookupKey = data.Name.ValueString()
	}

	// Get Role
	role, err := d.roleClient.GetOrgRole(ctx, lookupKey)
	if err != nil {
		if cerrs.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, "role", err)
		return
	}

	// Get rights
	rights, err := d.roleClient.ListOrgRoleRights(ctx, role.ID)
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, "role", err)
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

	data.ID.Set(role.ID)
	data.Name.Set(role.Name)
	data.ReadOnly.Set(role.ReadOnly)
	data.Description.Set(role.Description)

	// Save data into Terraform data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
