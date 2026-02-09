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

// Package bms provides a Terraform datasource.
package bms

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &BMSDataSource{}
	_ datasource.DataSourceWithConfigure = &BMSDataSource{}
)

func NewBMSDataSource() datasource.DataSource {
	return &BMSDataSource{}
}

type BMSDataSource struct { //nolint: revive
	client *client.CloudAvenue
}

// Init Initializes the data source.
func (d *BMSDataSource) Init(_ context.Context, _ *bmsModelDatasource) (diags diag.Diagnostics) {
	return diags
}

func (d *BMSDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *BMSDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = bmsSchema(ctx).GetDataSource(ctx)
}

func (d *BMSDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BMSDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_bms", d.client.GetOrgName(), metrics.Read)()

	config := &bmsModelDatasource{}

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

	// Set default timeouts
	readTimeout, diags := config.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	// Read data from the API
	bms, err := d.client.CAVSDK.V1.BMS.List()
	if err != nil {
		resp.Diagnostics.AddError("error on list BMS(s)", err.Error())
		return
	}
	data := []*bmsModelDatasourceEnv{}
	for _, b := range *bms {
		// Set Network
		net := networkToTerraform(&b)

		// Set BMS
		bms := BMSToTerraform(ctx, &b)

		// Set data
		x := newBMSModelDatasourceEnv(ctx)
		x.Network.Set(ctx, net)
		x.BMS.Set(ctx, bms)
		data = append(data, x)
	}
	// Set List
	config.Env.Set(ctx, data)

	// Set ID
	config.ID.Set(d.client.GetOrgName())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
