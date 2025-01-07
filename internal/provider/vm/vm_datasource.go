/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vm provides a Terraform datasource.
package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminvdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

var (
	_ datasource.DataSource              = &vmDataSource{}
	_ datasource.DataSourceWithConfigure = &vmDataSource{}
)

func NewVMDataSource() datasource.DataSource {
	return &vmDataSource{}
}

type vmDataSource struct {
	client   *client.CloudAvenue
	vdc      vdc.VDC
	adminVDC adminvdc.AdminVDC
	vapp     vapp.VAPP
	vm       vm.VM
}

// Init Initializes the data source.
func (d *vmDataSource) Init(ctx context.Context, dm *VMDataSourceModel) (diags diag.Diagnostics) {
	var mydiag diag.Diagnostics

	d.vdc, mydiag = vdc.Init(d.client, dm.VDC)
	diags.Append(mydiag...)
	if diags.HasError() {
		return
	}

	d.adminVDC, mydiag = adminvdc.Init(d.client, dm.VDC)
	diags.Append(mydiag...)
	if diags.HasError() {
		return
	}

	d.vapp, mydiag = vapp.Init(d.client, d.vdc, dm.VappID, dm.VappName)
	diags.Append(mydiag...)
	if diags.HasError() {
		return
	}

	if d.vapp.VAPP == nil {
		diags.AddError("Vapp not found", fmt.Sprintf("Vapp %s not found in VDC %s", dm.VappName, dm.VDC))
		return
	}

	d.vm, mydiag = vm.Init(d.client, d.vapp, vm.GetVMOpts{
		ID:   dm.ID,
		Name: dm.Name,
	})
	diags.Append(mydiag...)

	return
}

func (d *vmDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *vmDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vmSuperSchema(ctx).GetDataSource(ctx)
}

func (d *vmDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vmDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vm", d.client.GetOrgName(), metrics.Read)()

	config := &VMDataSourceModel{}

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

	// Read data from API
	data, mydiag := d.read(ctx, config, config)
	resp.Diagnostics.Append(mydiag...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// read is a common function for VM read.
func (d *vmDataSource) read(ctx context.Context, dm, dmPlan *VMDataSourceModel) (plan *VMDataSourceModel, diags diag.Diagnostics) {
	if err := d.vm.Refresh(); err != nil {
		diags.AddError("Error refreshing VM", err.Error())
		return
	}

	// ? State
	stateStruct, err := d.vm.StateRead(ctx)
	if err != nil {
		diags.AddError(
			"Unable to get VM state",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	// ? Resource
	networks, err := d.vm.NetworksRead()
	if err != nil {
		diags.AddError(
			"Unable to get VM networks",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	// ? Settings
	settings, err := d.vm.SettingsRead(ctx, dmPlan.Settings.Attributes()["customization"])
	if err != nil {
		diags.AddError(
			"Unable to get VM settings",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	return &VMDataSourceModel{
		ID:          types.StringValue(d.vm.GetID()),
		VDC:         types.StringValue(d.vdc.GetName()),
		Name:        types.StringValue(d.vm.GetName()),
		VappID:      types.StringValue(d.vapp.GetID()),
		VappName:    types.StringValue(d.vapp.GetName()),
		Description: dm.Description,
		State:       stateStruct.ToPlan(ctx),
		Resource:    d.vm.ResourceRead(ctx).ToPlan(ctx, networks),
		Settings:    settings.ToPlan(ctx),
	}, nil
}
