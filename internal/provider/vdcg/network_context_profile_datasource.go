/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdkv1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &networkContextProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &networkContextProfileDataSource{}
)

func NewNetworkContextProfileDataSource() datasource.DataSource {
	return &networkContextProfileDataSource{}
}

type networkContextProfileDataSource struct {
	client   *client.CloudAvenue
	vdcGroup *sdkv1.VDCGroup
}

func (d *networkContextProfileDataSource) Init(_ context.Context, dm *networkContextProfileModelDatasource) (diags diag.Diagnostics) {
	var err error

	idOrName := dm.VDCGroupName.Get()
	if dm.VDCGroupID.IsKnown() && dm.VDCGroupID.Get() != "" {
		idOrName = dm.VDCGroupID.Get()
	}

	d.vdcGroup, err = d.client.CAVSDK.V1.VDC().GetVDCGroup(idOrName)
	if err != nil {
		diags.AddError("Error retrieving VDC Group", err.Error())
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
	defer metrics.New("data.cloudavenue_vdcg_network_context_profile", d.client.GetOrgName(), metrics.Read)()

	config := &networkContextProfileModelDatasource{}
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		profile *sdkv1.NetworkContextProfile
		err     error
	)

	if config.ID.IsKnown() && config.ID.Get() != "" {
		profile, err = d.vdcGroup.GetNetworkContextProfileByID(config.ID.Get())
	} else {
		profile, err = d.vdcGroup.GetNetworkContextProfileByName(config.Name.Get())
	}

	if err != nil {
		nameOrID := config.Name.Get()
		if config.ID.Get() != "" {
			nameOrID = config.ID.Get()
		}
		resp.Diagnostics.AddError(
			"Network Context Profile not found",
			fmt.Sprintf("No Network Context Profile found with name or ID %q on VDC Group %q: %s", nameOrID, d.vdcGroup.GetName(), err),
		)
		return
	}

	stateRefreshed := config.Copy()
	stateRefreshed.ID.Set(profile.ID)
	stateRefreshed.Name.Set(profile.Name)
	stateRefreshed.Description.Set(profile.Description)
	stateRefreshed.Scope.Set(string(profile.Scope))
	stateRefreshed.VDCGroupID.Set(d.vdcGroup.GetID())
	stateRefreshed.VDCGroupName.Set(d.vdcGroup.GetName())

	appIDBlock, domainBlock, diags := attributesFromSDKProfile(ctx, profile)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(stateRefreshed.AppID.Set(ctx, appIDBlock)...)
	resp.Diagnostics.Append(stateRefreshed.DomainName.Set(ctx, domainBlock)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}
