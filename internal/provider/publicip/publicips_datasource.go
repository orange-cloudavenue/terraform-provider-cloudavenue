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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
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
	// Client is a terraform Client
	client *client.CloudAvenue

	// eClient is the Edge Gateway client from the SDK V2
	eClient *edgegateway.Client
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

func (d *publicIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_publicips", d.client.GetOrgName(), metrics.Read)()

	data := &publicIPDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ips, err := d.eClient.ListPublicIP(ctx, types.ParamsEdgeGateway{
		Name: data.EdgeGatewayName.Get(),
		ID:   data.EdgeGatewayID.Get(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error while getting public IPs", err.Error())
		return
	}

	data.EdgeGatewayID.Set(ips.EdgegatewayID)
	data.EdgeGatewayName.Set(ips.EdgegatewayName)

	listOfIps := make([]string, 0)
	ipPubs := make([]*publicIPNetworkConfigModel, 0)
	for _, ip := range ips.PublicIPs {
		x := &publicIPNetworkConfigModel{
			ID:       supertypes.NewStringNull(),
			PublicIP: supertypes.NewStringNull(),
		}

		x.ID.Set(ip.IP) // For maintaining compatibility use IP for ID and not ID
		x.PublicIP.Set(ip.IP)

		ipPubs = append(ipPubs, x)
		listOfIps = append(listOfIps, ip.IP)
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
