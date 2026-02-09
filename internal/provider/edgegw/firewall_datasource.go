/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package edgegw provides a Terraform datasource.
package edgegw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

var (
	_ datasource.DataSource              = &firewallDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallDataSource{}
)

func NewFirewallDataSource() datasource.DataSource {
	return &firewallDataSource{}
}

type firewallDataSource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the data source.
func (d *firewallDataSource) Init(_ context.Context, dm *firewallModel) (diags diag.Diagnostics) {
	var err error

	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return diags
	}

	d.edgegw, err = d.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   dm.EdgeGatewayID.StringValue,
		Name: dm.EdgeGatewayName.StringValue,
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return diags
	}

	return diags
}

func (d *firewallDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_firewall"
}

func (d *firewallDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = firewallSchema(ctx).GetDataSource(ctx)
}

func (d *firewallDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *firewallDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_edgegateway_firewall", d.client.GetOrgName(), metrics.Read)()

	config := &firewallModel{}

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

	/*
		Implement the data source read logic here.
	*/

	s := &firewallResource{
		client: d.client,
		org:    d.org,
		edgegw: d.edgegw,
	}

	// Read data from the API
	data, _, diags := s.read(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
