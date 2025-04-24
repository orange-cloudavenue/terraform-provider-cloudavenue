/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package tier0 provides a Terraform resource to manage Tier-0 VRFs.
package vrf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &tier0VrfsDataSource{}
	_ datasource.DataSourceWithConfigure = &tier0VrfsDataSource{}
)

// NewTier0VrfsDataSource returns a new resource implementing the Tier-0 VRFs data source.
func NewTier0VrfsDataSource() datasource.DataSource {
	return &tier0VrfsDataSource{}
}

type tier0VrfsDataSource struct {
	client *client.CloudAvenue
}

func (d *tier0VrfsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "vrfs"
}

func (d *tier0VrfsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tier0VrfsSchema().GetDataSource(ctx)
}

func (d *tier0VrfsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tier0VrfsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_tier0_vrfs", d.client.GetOrgName(), metrics.Read)()

	var data tier0VrfsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the list of Tier-0 VRFs
	t0s, err := d.client.CAVSDK.V1.T0.GetT0s()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read list T0, got error: %s", err))
		return
	}

	var names []string

	for _, t0 := range *t0s {
		names = append(names, t0.GetName())
	}

	// Generate a UUID from the list of names
	data.ID.Set(utils.GenerateUUID(names...).String())

	// Save data into Terraform state
	resp.Diagnostics.Append(data.Names.Set(ctx, names)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
