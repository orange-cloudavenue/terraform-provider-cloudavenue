/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vrf provides a Terraform resource to manage Tier-0 VRFs.
package vrf

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &tier0VrfDataSource{}
	_ datasource.DataSourceWithConfigure = &tier0VrfDataSource{}
)

// NewTier0VrfDataSource returns a new datasource implementing the tier0_vrf data source.
func NewTier0VrfDataSource() datasource.DataSource {
	return &tier0VrfDataSource{}
}

type tier0VrfDataSource struct {
	client *client.CloudAvenue
}

func (d *tier0VrfDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "vrf"
}

func (d *tier0VrfDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tier0VrfSchema().GetDataSource(ctx)
}

func (d *tier0VrfDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tier0VrfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_tier0_vrf", d.client.GetOrgName(), metrics.Read)()

	var data tier0VrfDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	t0, err := d.client.CAVSDK.V1.T0.GetT0(data.Name.Get())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Tier-0, got error: %s", err))
		return
	}

	data.ID.Set(utils.GenerateUUID(t0.GetName()).String())
	data.Provider.Set(t0.Tier0Provider)
	data.Name.Set(t0.GetName())
	data.ClassService.Set(string(t0.ClassService))

	var services []*segmentModel

	for _, segment := range t0.Services {
		s := &segmentModel{
			Service: supertypes.NewStringNull(),
			VLANID:  supertypes.NewStringNull(),
		}

		s.Service.Set(segment.Service)
		switch segment.VLANID.(type) { //nolint:gocritic
		case int:
			s.VLANID.Set(fmt.Sprintf("%d", segment.VLANID))
		case string:
			s.VLANID.Set(segment.VLANID.(string))
		}

		services = append(services, s)
	}

	resp.Diagnostics.Append(data.Services.Set(ctx, services)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
