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
package edgegw //nolint:dupl // This is a datasource, it is normal to have similar code to the other datasource.

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

var (
	_ datasource.DataSource              = &dhcpForwardingDataSource{}
	_ datasource.DataSourceWithConfigure = &dhcpForwardingDataSource{}
)

func NewDhcpForwardingDataSource() datasource.DataSource {
	return &dhcpForwardingDataSource{}
}

type dhcpForwardingDataSource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the data source.
func (d *dhcpForwardingDataSource) Init(_ context.Context, dm *DhcpForwardingModel) (diags diag.Diagnostics) {
	var err error

	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return diags
	}

	d.edgegw, err = d.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(dm.EdgeGatewayID.Get()),
		Name: types.StringValue(dm.EdgeGatewayName.Get()),
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return diags
	}

	return diags
}

func (d *dhcpForwardingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_dhcp_forwarding"
}

func (d *dhcpForwardingDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dhcpForwardingSchema(ctx).GetDataSource(ctx)
}

func (d *dhcpForwardingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dhcpForwardingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_edgegateway_dhcp_forwarding", d.client.GetOrgName(), metrics.Read)()

	config := &DhcpForwardingModel{}

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

	// If read function is identical to the resource, you can use the following code:
	s := &dhcpForwardingResource{
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

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
