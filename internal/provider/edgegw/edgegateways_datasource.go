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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
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
	client *client.CloudAvenue
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

	d.client = client
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

	gateways, err := d.client.CAVSDK.V1.EdgeGateway.List()
	if err != nil {
		resp.Diagnostics.AddError("Unable to list edge gateways", err.Error())
		return
	}

	gws := make([]*edgeGatewayDataSourceModelEdgeGateway, 0)
	for _, edge := range *gateways {
		gw := new(edgeGatewayDataSourceModelEdgeGateway)
		gw.ID.Set(urn.Normalize(urn.Gateway, edge.GetID()).String())
		gw.Name.Set(edge.GetName())
		gw.Description.Set(edge.GetDescription())
		gw.OwnerType.Set(string(edge.GetOwnerType()))
		gw.OwnerName.Set(edge.GetOwnerName())
		gw.Tier0VrfName.Set(edge.GetT0())

		gws = append(gws, gw)
		names = append(names, edge.GetName())
	}

	data.ID.Set(utils.GenerateUUID(names).ValueString())
	data.EdgeGateways.Set(ctx, gws)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
