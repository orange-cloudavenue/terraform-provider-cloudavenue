/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vdc provides a resource to manage VDCs.
package vdc

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &vdcsDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcsDataSource{}
)

// NewVDCsDataSource returns a new resource implementing the vdcs data source.
func NewVDCsDataSource() datasource.DataSource {
	return &vdcsDataSource{}
}

type vdcsDataSource struct {
	// Client is a terraform Client
	client *client.CloudAvenue

	// vClient is the VDC client from the SDK V2
	vClient *vdc.Client
}

func (d *vdcsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *vdcsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vdcsSchema().GetDataSource(ctx)
}

func (d *vdcsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	vC, err := vdc.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create VDC client, got error: %s", err),
		)
		return
	}

	d.client = client
	d.vClient = vC
}

func (d *vdcsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdcs", d.client.GetOrgName(), metrics.Read)()

	var (
		state    = new(vdcsDataSourceModel)
		names    []string
		dataVDCs = make([]*vdcsDataSourceModelVDC, 0)
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vdcs, err := d.vClient.ListVDC(ctx, types.ParamsListVDC{})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list VDCs", err.Error())
		return
	}

	for _, v := range vdcs.VDCS {
		x := &vdcsDataSourceModelVDC{
			ID:          supertypes.NewStringNull(),
			Name:        supertypes.NewStringNull(),
			Description: supertypes.NewStringNull(),
		}

		x.ID.Set(v.ID)
		x.Name.Set(v.Name)
		x.Description.Set(v.Description)

		dataVDCs = append(dataVDCs, x)
		names = append(names, v.Name)
	}

	state.ID.Set(utils.GenerateUUID(names).String())
	state.VDCs.DiagsSet(ctx, resp.Diagnostics, dataVDCs)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
