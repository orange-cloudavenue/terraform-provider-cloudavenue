// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	edge "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &edgeGatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &edgeGatewaysDataSource{}
)

// NewEdgeGatewayDataSource returns a new datasource implementing the edge_gateway data source.
func NewEdgeGatewayDataSource() datasource.DataSource {
	return &edgeGatewayDataSource{}
}

type edgeGatewayDataSource struct {
	client     *client.CloudAvenue
	edgeClient *edge.Client
}

func (d *edgeGatewayDataSource) Init(_ context.Context) (diags diag.Diagnostics) {
	var err error
	var v2 cav.Client

	if v2, err = d.client.V2(); err == nil {
		d.edgeClient, err = edge.New(v2)
	}
	if err != nil {
		diags.AddError("Error initializing EdgeGateway client", err.Error())
		return diags
	}

	return diags
}

func (d *edgeGatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *edgeGatewayDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = edgegwSchema().GetDataSource(ctx)
}

func (d *edgeGatewayDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *edgeGatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_edgegateway", d.client.GetOrgName(), metrics.Read)()

	config := &edgeGatewayDatasourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	s := &edgeGatewayResource{
		client:     d.client,
		edgeClient: d.edgeClient,
	}

	configResource := &edgeGatewayResourceModel{
		ID:   config.ID,
		Name: config.Name,
	}

	data, found, diags := s.read(ctx, configResource)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The edge gateway %s(%s) was not found.", config.Name.Get(), config.ID.Get()))
		return
	}

	data.ID.Set(data.ID.Get())
	data.Name.Set(data.Name.Get())
	data.Tier0VRFName.Set(data.Tier0VRFName.Get())
	data.OwnerName.Set(data.OwnerName.Get())
	data.Description.Set(data.Description.Get())
	data.Bandwidth.SetInt(data.Bandwidth.GetInt())

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
