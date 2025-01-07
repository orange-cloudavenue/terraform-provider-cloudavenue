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

	commonutils "github.com/orange-cloudavenue/common-go/utils"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

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
	client *client.CloudAvenue
}

// Init Initializes the resource.
func (d *vdcsDataSource) Init(ctx context.Context, rm *vdcsDataSourceModel) (diags diag.Diagnostics) {
	return
}

func (d *vdcsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *vdcsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vdcsSchema().GetDataSource(ctx)
}

func (d *vdcsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vdcsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdcs", d.client.GetOrgName(), metrics.Read)()

	var (
		state    = new(vdcsDataSourceModel)
		names    []string
		dataVDCs = make([]*vdcRef, 0)
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vdcs, err := d.client.CAVSDK.V1.Querier().List().VDC()
	if err != nil {
		resp.Diagnostics.AddError("Unable to list VDCs", err.Error())
		return
	}

	for _, v := range vdcs {
		x := &vdcRef{
			ID:   supertypes.NewStringNull(),
			Name: supertypes.NewStringNull(),
		}

		// Extract ID from href
		uuid, err := commonutils.GetUUIDFromHref(v.HREF, true)
		if err != nil {
			resp.Diagnostics.AddError("Unable to extract VDC UUID", err.Error())
			return
		}

		x.Name.Set(v.Name)
		x.ID.Set(uuid)

		dataVDCs = append(dataVDCs, x)
		names = append(names, v.Name)
	}

	state.ID.Set(utils.GenerateUUID(names).String())
	state.VDCs.Set(ctx, dataVDCs)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
