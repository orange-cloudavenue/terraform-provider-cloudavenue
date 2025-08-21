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

	edgegateway "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
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

	// eClient is the Edge Gateway client from the SDK V2
	eClient *edgegateway.Client
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

	eC, err := edgegateway.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Edge Gateway client, got error: %s", err),
		)
		return
	}

	d.client = client
	d.eClient = eC
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
	t0s, err := d.eClient.ListT0(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read list T0, got error: %s", err))
		return
	}

	var names []string

	for _, t0 := range t0s.T0s {
		names = append(names, t0.Name)
	}

	// Generate a UUID from the list of names
	data.ID.Set(utils.GenerateUUID(names...).String())

	// Save data into Terraform state
	resp.Diagnostics.Append(data.Names.Set(ctx, names)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataT0s := make([]*tier0VrfDataSourceModel, 0, len(t0s.T0s))

	for _, t0 := range t0s.T0s {
		t0Model := &tier0VrfDataSourceModel{}
		resp.Diagnostics.Append(t0Model.fromAPI(ctx, &t0)...)
		if resp.Diagnostics.HasError() {
			return
		}
		dataT0s = append(dataT0s, t0Model)
	}

	resp.Diagnostics.Append(data.T0s.Set(ctx, dataT0s)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
