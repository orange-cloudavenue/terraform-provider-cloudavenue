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

package edgegw

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	AppPortProfileModel struct {
		ID              supertypes.StringValue                                         `tfsdk:"id"`
		Name            supertypes.StringValue                                         `tfsdk:"name"`
		EdgeGatewayID   supertypes.StringValue                                         `tfsdk:"edge_gateway_id"`
		EdgeGatewayName supertypes.StringValue                                         `tfsdk:"edge_gateway_name"`
		Description     supertypes.StringValue                                         `tfsdk:"description"`
		AppPorts        supertypes.ListNestedObjectValueOf[AppPortProfileModelAppPort] `tfsdk:"app_ports"`
	}

	AppPortProfileModelDatasource struct {
		ID              supertypes.StringValue                                         `tfsdk:"id"`
		Name            supertypes.StringValue                                         `tfsdk:"name"`
		EdgeGatewayID   supertypes.StringValue                                         `tfsdk:"edge_gateway_id"`
		EdgeGatewayName supertypes.StringValue                                         `tfsdk:"edge_gateway_name"`
		Description     supertypes.StringValue                                         `tfsdk:"description"`
		AppPorts        supertypes.ListNestedObjectValueOf[AppPortProfileModelAppPort] `tfsdk:"app_ports"`
		Scope           supertypes.StringValue                                         `tfsdk:"scope"`
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
func (rm *AppPortProfileModel) toSDKAppPortProfilePorts(ctx context.Context) (appPortProfilePorts []types.ParamsAppPortProfilePort, diags diag.Diagnostics) {
	appPortProfilePorts = make([]types.ParamsAppPortProfilePort, 0)

	appPorts, d := rm.AppPorts.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	for _, appPort := range appPorts {
		destPorts, d := appPort.Ports.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		appPortProfilePorts = append(appPortProfilePorts, types.ParamsAppPortProfilePort{
			Protocol:         appPort.Protocol.Get(),
			DestinationPorts: destPorts,
		})
	}

	return appPortProfilePorts, diags
}
