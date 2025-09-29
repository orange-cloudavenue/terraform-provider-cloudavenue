/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	ServicesModel struct {
		ID              supertypes.StringValue                                  `tfsdk:"id"`
		EdgeGatewayName supertypes.StringValue                                  `tfsdk:"edge_gateway_name"`
		EdgeGatewayID   supertypes.StringValue                                  `tfsdk:"edge_gateway_id"`
		Network         supertypes.StringValue                                  `tfsdk:"network"`
		IPAddress       supertypes.StringValue                                  `tfsdk:"ip_address"`
		Services        supertypes.MapNestedObjectValueOf[ServicesModelCatalog] `tfsdk:"services"`
	}

	ServicesModelCatalog struct {
		Network  supertypes.StringValue                                         `tfsdk:"network"`
		Category supertypes.StringValue                                         `tfsdk:"category"`
		Services supertypes.MapNestedObjectValueOf[ServicesModelCatalogService] `tfsdk:"services"`
	}

	ServicesModelCatalogService struct {
		Name        supertypes.StringValue                                               `tfsdk:"name"`
		Description supertypes.StringValue                                               `tfsdk:"description"`
		IPs         supertypes.ListValueOf[string]                                       `tfsdk:"ips"`
		FQDNs       supertypes.ListValueOf[string]                                       `tfsdk:"fqdns"`
		Ports       supertypes.ListNestedObjectValueOf[ServicesModelCatalogServicePorts] `tfsdk:"ports"`
	}

	ServicesModelCatalogServicePorts struct {
		Port     supertypes.Int32Value  `tfsdk:"port"`
		Protocol supertypes.StringValue `tfsdk:"protocol"`
	}
)

func (rm *ServicesModel) Copy() *ServicesModel {
	x := &ServicesModel{}
	utils.ModelCopy(rm, x)
	return x
}
