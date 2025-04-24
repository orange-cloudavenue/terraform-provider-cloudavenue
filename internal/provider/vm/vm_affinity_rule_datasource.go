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
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

var (
	_ datasource.DataSource              = &vmAffinityRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &vmAffinityRuleDataSource{}
)

func NewVMAffinityRuleDatasource() datasource.DataSource {
	return &vmAffinityRuleDataSource{}
}

type vmAffinityRuleDataSource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
}

func (d *vmAffinityRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_affinity_rule"
}

func (d *vmAffinityRuleDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vmAffinityRuleSchema().GetDataSource(ctx)
}

func (d *vmAffinityRuleDataSource) Init(_ context.Context, rm *vmAffinityRuleDataSourceModel) (diags diag.Diagnostics) {
	d.vdc, diags = vdc.Init(d.client, rm.VDC)

	return
}

func (d *vmAffinityRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vmAffinityRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vm_affinity_rule", d.client.GetOrgName(), metrics.Read)()

	var data vmAffinityRuleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(d.Init(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmAffinityRule, err := getVMAffinityRule(d.vdc, data.Name.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read affinity rule", err.Error())
		return
	}

	data = vmAffinityRuleDataSourceModel{
		ID:       types.StringValue(vmAffinityRule.VmAffinityRule.ID),
		VDC:      types.StringValue(d.vdc.GetName()),
		Name:     types.StringValue(vmAffinityRule.VmAffinityRule.Name),
		Required: types.BoolValue(*vmAffinityRule.VmAffinityRule.IsMandatory),
		Enabled:  types.BoolValue(*vmAffinityRule.VmAffinityRule.IsEnabled),
		Polarity: types.StringValue(vmAffinityRule.VmAffinityRule.Polarity),
	}

	endpointVMs := vmReferencesToListValue(vmAffinityRule.VmAffinityRule.VmReferences)
	data.VMIDs = types.SetValueMust(types.StringType, endpointVMs)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
