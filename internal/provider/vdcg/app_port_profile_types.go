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

package vdcg

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/helpers"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	AppPortProfileModel struct {
		ID           supertypes.StringValue                                         `tfsdk:"id"`
		Name         supertypes.StringValue                                         `tfsdk:"name"`
		VDCGroupID   supertypes.StringValue                                         `tfsdk:"vdc_group_id"`
		VDCGroupName supertypes.StringValue                                         `tfsdk:"vdc_group_name"`
		Description  supertypes.StringValue                                         `tfsdk:"description"`
		AppPorts     supertypes.ListNestedObjectValueOf[AppPortProfileModelAppPort] `tfsdk:"app_ports"`
	}

	AppPortProfileModelDatasource struct {
		ID           supertypes.StringValue                                         `tfsdk:"id"`
		Name         supertypes.StringValue                                         `tfsdk:"name"`
		VDCGroupID   supertypes.StringValue                                         `tfsdk:"vdc_group_id"`
		VDCGroupName supertypes.StringValue                                         `tfsdk:"vdc_group_name"`
		Description  supertypes.StringValue                                         `tfsdk:"description"`
		AppPorts     supertypes.ListNestedObjectValueOf[AppPortProfileModelAppPort] `tfsdk:"app_ports"`
		Scope        supertypes.StringValue                                         `tfsdk:"scope"`
	}

	AppPortProfileModelAppPort struct {
		Protocol supertypes.StringValue        `tfsdk:"protocol"`
		Ports    supertypes.SetValueOf[string] `tfsdk:"ports"`
	}
)

func (rm *AppPortProfileModelDatasource) Copy() *AppPortProfileModelDatasource {
	x := &AppPortProfileModelDatasource{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *AppPortProfileModel) Copy() *AppPortProfileModel {
	x := &AppPortProfileModel{}
	utils.ModelCopy(rm, x)
	return x
}

// toSDKAppPortProfile converts the AppPortProfileModel to the SDK representation.
func (rm *AppPortProfileModel) toSDKAppPortProfile(ctx context.Context) (nsxtAppPortProfilePorts *v1.FirewallGroupAppPortProfileModel, diags diag.Diagnostics) {
	nsxtAppPortProfilePorts = &v1.FirewallGroupAppPortProfileModel{}

	nsxtAppPortProfilePorts.ID = rm.ID.Get()
	nsxtAppPortProfilePorts.Name = rm.Name.Get()
	nsxtAppPortProfilePorts.Description = rm.Description.Get()

	appPorts, d := rm.AppPorts.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	nsxtAppPortProfilePorts.ApplicationPorts = make(v1.FirewallGroupAppPortProfileModelPorts, 0)
	for _, appPort := range appPorts {
		destPorts, d := appPort.Ports.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		protocol, err := helpers.ParseFirewallAppPortProfileProtocol(appPort.Protocol.Get())
		if err != nil {
			diags.AddError("Error parsing protocol", err.Error())
			return nil, diags
		}

		nsxtAppPortProfilePorts.ApplicationPorts = append(nsxtAppPortProfilePorts.ApplicationPorts, v1.FirewallGroupAppPortProfileModelPort{
			Protocol:         protocol,
			DestinationPorts: destPorts,
		})
	}

	return nsxtAppPortProfilePorts, diags
}
