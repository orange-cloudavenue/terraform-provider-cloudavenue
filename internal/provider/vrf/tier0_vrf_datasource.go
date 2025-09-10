/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vrf provides a Terraform resource to manage Tier-0 VRFs.
package vrf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	edgegateway "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &tier0VrfDataSource{}
	_ datasource.DataSourceWithConfigure = &tier0VrfDataSource{}
)

// NewTier0VrfDataSource returns a new datasource implementing the tier0_vrf data source.
func NewTier0VrfDataSource() datasource.DataSource {
	return &tier0VrfDataSource{}
}

type tier0VrfDataSource struct {
	// Client is a terraform Client
	client *client.CloudAvenue

	// eClient is the Edge Gateway client from the SDK V2
	eClient *edgegateway.Client
}

func (d *tier0VrfDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "vrf"
}

func (d *tier0VrfDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tier0VrfSchema().GetDataSource(ctx)
}

func (d *tier0VrfDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tier0VrfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_tier0_vrf", d.client.GetOrgName(), metrics.Read)()

	var data tier0VrfDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	t0, err := d.eClient.GetT0(ctx, types.ParamsGetT0{
		T0Name:          data.Name.Get(),
		EdgegatewayID:   data.EdgeGatewayID.Get(),
		EdgegatewayName: data.EdgeGatewayName.Get(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Tier-0, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(data.fromAPI(ctx, t0)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
