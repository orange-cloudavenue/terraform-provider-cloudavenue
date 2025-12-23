/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package publicip provides a Terraform resource to manage public IPs.
package publicip

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &publicIPDataSource{}
	_ datasource.DataSourceWithConfigure = &publicIPDataSource{}
)

// NewPublicIPDataSource returns a new resource implementing the public IP data source.
func NewPublicIPDataSource() datasource.DataSource {
	return &publicIPDataSource{}
}

type publicIPDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
}

// Init.
func (d *publicIPDataSource) Init(_ context.Context, _ *publicIPDataSourceModel) (diags diag.Diagnostics) {
	d.adminOrg, diags = adminorg.Init(d.client)
	return diags
}

func (d *publicIPDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *publicIPDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = publicIPsSchema(ctx).GetDataSource(ctx)
}

func (d *publicIPDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *publicIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_publicips", d.client.GetOrgName(), metrics.Read)()

	data := &publicIPDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ips, err := d.client.CAVSDK.V1.PublicIP.GetIPs()
	if err != nil {
		resp.Diagnostics.AddError("Error while getting public IPs", err.Error())
		return
	}

	listOfIps := make([]string, 0)
	ipPubs := make([]*publicIPNetworkConfigModel, 0)
	for _, cfg := range ips.NetworkConfig {
		edgeGateway, err := d.adminOrg.GetEdgeGateway(edgegw.BaseEdgeGW{
			Name: types.StringValue(cfg.EdgeGatewayName),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error while getting edge gateway", err.Error())
			return
		}

		x := &publicIPNetworkConfigModel{
			ID:              supertypes.NewStringNull(),
			EdgeGatewayName: supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringNull(),
			PublicIP:        supertypes.NewStringNull(),
		}

		x.ID.Set(cfg.GetIP())
		x.EdgeGatewayName.Set(edgeGateway.GetName())
		x.EdgeGatewayID.Set(edgeGateway.GetID())
		x.PublicIP.Set(cfg.GetIP())

		ipPubs = append(ipPubs, x)
		listOfIps = append(listOfIps, cfg.GetIP())
	}

	data.ID.Set(utils.GenerateUUID(listOfIps).String())
	resp.Diagnostics.Append(data.PublicIPs.Set(ctx, ipPubs)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
