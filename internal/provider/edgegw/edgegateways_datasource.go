/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &edgeGatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &edgeGatewaysDataSource{}
)

// NewEdgeGatewaysDataSource returns a new resource implementing the edge_gateways data source.
func NewEdgeGatewaysDataSource() datasource.DataSource {
	return &edgeGatewaysDataSource{}
}

type edgeGatewaysDataSource struct {
	// Client is a terraform Client
	client *client.CloudAvenue

	// eClient is the Edge Gateway client from the SDK V2
	eClient *edgegateway.Client
}

func (d *edgeGatewaysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *edgeGatewaysDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = edgeGatewaysSuperSchema(ctx).GetDataSource(ctx)
}

func (d *edgeGatewaysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *edgeGatewaysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_edgegateways", d.client.GetOrgName(), metrics.Read)()
	var (
		data  = new(edgeGatewaysDataSourceModel)
		names = make([]string, 0)
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gateways, err := d.eClient.ListEdgeGateway(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list edge gateways", err.Error())
		return
	}

	gws := make([]*edgeGatewaysDataSourceModelEdgeGateway, 0)
	for _, edge := range gateways.EdgeGateways {
		gw := new(edgeGatewaysDataSourceModelEdgeGateway)
		gw.ID.Set(edge.ID)
		gw.Name.Set(edge.Name)
		gw.Description.Set(edge.Description)

		gw.OwnerName.Set(edge.OwnerRef.Name)
		gw.OwnerID.Set(edge.OwnerRef.ID)

		gw.T0ID.Set(edge.UplinkT0.ID)
		gw.T0Name.Set(edge.UplinkT0.Name)

		gw.Tier0VRFName.Set(edge.UplinkT0.Name)
		gws = append(gws, gw)
		names = append(names, edge.Name)
	}

	data.ID.Set(utils.GenerateUUID(names).ValueString())

	if len(gws) == 0 {
		data.EdgeGateways.SetNull(ctx)
	} else {
		resp.Diagnostics.Append(data.EdgeGateways.Set(ctx, gws)...)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
