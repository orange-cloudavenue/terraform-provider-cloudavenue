/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package elb

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	VirtualServiceModel struct {
		ID                     supertypes.StringValue                                             `tfsdk:"id"`
		Name                   supertypes.StringValue                                             `tfsdk:"name"`
		EdgeGatewayName        supertypes.StringValue                                             `tfsdk:"edge_gateway_name"`
		EdgeGatewayID          supertypes.StringValue                                             `tfsdk:"edge_gateway_id"`
		Description            supertypes.StringValue                                             `tfsdk:"description"`
		Enabled                supertypes.BoolValue                                               `tfsdk:"enabled"`
		PoolName               supertypes.StringValue                                             `tfsdk:"pool_name"`
		PoolID                 supertypes.StringValue                                             `tfsdk:"pool_id"`
		ServiceEngineGroupName supertypes.StringValue                                             `tfsdk:"service_engine_group_name"`
		VirtualIP              supertypes.StringValue                                             `tfsdk:"virtual_ip"`
		ServiceType            supertypes.StringValue                                             `tfsdk:"service_type"`
		CertificateID          supertypes.StringValue                                             `tfsdk:"certificate_id"`
		ServicePorts           supertypes.ListNestedObjectValueOf[VirtualServiceModelServicePort] `tfsdk:"service_ports"`
	}

	VirtualServiceModelServicePort struct {
		Start supertypes.Int64Value `tfsdk:"start"`
		End   supertypes.Int64Value `tfsdk:"end"`
	}
)

func (rm *VirtualServiceModel) Copy() *VirtualServiceModel {
	x := &VirtualServiceModel{}
	utils.ModelCopy(rm, x)
	return x
}

// ToSDKVirtualServiceGroupModel converts the model to the SDK model.
func (rm *VirtualServiceModel) ToSDKVirtualServiceModelRequest(ctx context.Context, c edgeloadbalancer.Client) (*edgeloadbalancer.VirtualServiceModelRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	poolID := rm.PoolID.Get()
	if poolID == "" {
		pool, err := c.GetPool(ctx, rm.EdgeGatewayID.Get(), rm.PoolName.Get())
		if err != nil {
			diags.AddError("Error getting pool", err.Error())
			return nil, diags
		}

		poolID = pool.ID
	}

	vs := &edgeloadbalancer.VirtualServiceModelRequest{
		Name:               rm.Name.Get(),
		Description:        rm.Description.Get(),
		Enabled:            rm.Enabled.GetPtr(),
		ApplicationProfile: edgeloadbalancer.VirtualServiceModelApplicationProfile(rm.ServiceType.Get()),
		VirtualIPAddress:   rm.VirtualIP.Get(),
		EdgeGatewayID:      rm.EdgeGatewayID.Get(),
		PoolID:             poolID,
		CertificateID:      rm.CertificateID.GetPtr(),
		ServicePorts: func() []edgeloadbalancer.VirtualServiceModelServicePort {
			var ports []edgeloadbalancer.VirtualServiceModelServicePort
			sps, d := rm.ServicePorts.Get(ctx)
			if d.HasError() {
				diags = append(diags, d...)
				return nil
			}

			for _, port := range sps {
				ports = append(ports, edgeloadbalancer.VirtualServiceModelServicePort{
					Start: port.Start.GetIntPtr(),
					End:   port.End.GetIntPtr(),
				})
			}
			return ports
		}(),
	}
	if rm.ServiceEngineGroupName.IsKnown() {
		seg, err := c.GetVirtualService(ctx, rm.EdgeGatewayID.Get(), rm.ServiceEngineGroupName.Get())
		if err != nil {
			diags.AddError("Error getting service engine group", err.Error())
			return nil, diags
		}

		vs.ServiceEngineGroupID = &seg.ID
	}

	return vs, diags
}
