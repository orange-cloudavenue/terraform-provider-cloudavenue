/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package elb provides a Terraform datasource.
package elb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &VirtualServiceDataSource{}
	_ datasource.DataSourceWithConfigure = &VirtualServiceDataSource{}
)

func NewVirtualServiceDataSource() datasource.DataSource {
	return &VirtualServiceDataSource{}
}

type VirtualServiceDataSource struct {
	client *client.CloudAvenue
	elb    edgeloadbalancer.Client
	edge   *v1.EdgeClient
}

// Init Initializes the data source.
func (d *VirtualServiceDataSource) Init(_ context.Context, dm *VirtualServiceModel) (diags diag.Diagnostics) {
	var err error

	d.elb, err = edgeloadbalancer.NewClient()
	if err != nil {
		diags.AddError("Error creating elb client", err.Error())
	}

	eIDOrName := dm.EdgeGatewayID.Get()
	if eIDOrName == "" {
		eIDOrName = dm.EdgeGatewayName.Get()
	}
	d.edge, err = d.client.CAVSDK.V1.EdgeGateway.Get(eIDOrName)
	if err != nil {
		diags.AddError("Error creating edge client", err.Error())
	}

	dm.EdgeGatewayID.Set(urn.Normalize(urn.Gateway, d.edge.GetID()).String())
	dm.EdgeGatewayName.Set(d.edge.GetName())

	return
}

func (d *VirtualServiceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_virtual_service"
}

func (d *VirtualServiceDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = virtualServiceSchema(ctx).GetDataSource(ctx)
}

func (d *VirtualServiceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VirtualServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_elb_virtual_service", d.client.GetOrgName(), metrics.Read)()

	config := &VirtualServiceModel{}

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

	s := &VirtualServiceResource{
		client: d.client,
		elb:    d.elb,
		edge:   d.edge,
	}

	// Read data from the API
	data, found, diags := s.read(ctx, config)
	if !found {
		resp.Diagnostics.AddError("Resource not found", "The resource was not found")
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
