// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package edgegw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	edge "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &networkContextProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &networkContextProfileDataSource{}
)

// NewNetworkContextProfileDataSource returns a new context profile data source.
func NewNetworkContextProfileDataSource() datasource.DataSource {
	return &networkContextProfileDataSource{}
}

type networkContextProfileDataSource struct {
	client     *client.CloudAvenue
	edgeClient *edge.Client
	edgegw     *types.ModelEdgeGateway
}

// Init initializes the data source.
func (d *networkContextProfileDataSource) Init(ctx context.Context, dm *networkContextProfileModelDatasource) (diags diag.Diagnostics) {
	var err error
	var v2 cav.Client

	if v2, err = d.client.V2(); err == nil {
		d.edgeClient, err = edge.New(v2)
	}
	if err != nil {
		diags.AddError("Error creating Edge Gateway client", err.Error())
		return diags
	}

	d.edgegw, err = d.edgeClient.GetEdgeGateway(ctx, types.ParamsEdgeGateway{
		ID:   dm.EdgeGatewayID.Get(),
		Name: dm.EdgeGatewayName.Get(),
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return diags
	}

	return diags
}

func (d *networkContextProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_network_context_profile"
}

func (d *networkContextProfileDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = networkContextProfileSchema(ctx).GetDataSource(ctx)
}

func (d *networkContextProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *networkContextProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_edgegateway_network_context_profile", d.client.GetOrgName(), metrics.Read)()

	config := &networkContextProfileModelDatasource{}

	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := d.edgeClient.GetNetworkContextProfile(ctx, types.ParamsGetEdgeGatewayNetworkContextProfile{ID: config.ID.Get(), Name: config.ID.Get(), EdgeGatewayID: d.edgegw.ID, EdgeGatewayName: d.edgegw.Name})

	if err != nil {
		resp.Diagnostics.AddError(
			"Network Context Profile not found",
			fmt.Sprintf("No Network Context Profile found with Id: %s , Name: %s on Edge Gateway %q: %s", config.ID.Get(), config.ID.Get(), d.edgegw.Name, err),
		)
		return
	}

	stateRefreshed := config.Copy()

	resp.Diagnostics.Append(stateRefreshed.fromSDKProfile(ctx, profile, d.edgegw.ID, d.edgegw.Name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}
