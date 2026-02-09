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

package elb

import (
	"context"
	"strconv"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	PoolModel struct {
		ID              supertypes.StringValue                                     `tfsdk:"id"`
		Name            supertypes.StringValue                                     `tfsdk:"name"`
		Description     supertypes.StringValue                                     `tfsdk:"description"`
		EdgeGatewayID   supertypes.StringValue                                     `tfsdk:"edge_gateway_id"`
		EdgeGatewayName supertypes.StringValue                                     `tfsdk:"edge_gateway_name"`
		Enabled         supertypes.BoolValue                                       `tfsdk:"enabled"`
		Algorithm       supertypes.StringValue                                     `tfsdk:"algorithm"`
		DefaultPort     supertypes.Int64Value                                      `tfsdk:"default_port"`
		Members         supertypes.SingleNestedObjectValueOf[PoolModelMembers]     `tfsdk:"members"`
		Health          supertypes.SingleNestedObjectValueOf[PoolModelHealth]      `tfsdk:"health"`
		TLS             supertypes.SingleNestedObjectValueOf[PoolModelTLS]         `tfsdk:"tls"`
		Persistence     supertypes.SingleNestedObjectValueOf[PoolModelPersistence] `tfsdk:"persistence"`
	}

	PoolModelMembers struct {
		GracefulTimeoutPeriod supertypes.StringValue                                        `tfsdk:"graceful_timeout_period"`
		TargetGroup           supertypes.StringValue                                        `tfsdk:"target_group"`
		Targets               supertypes.ListNestedObjectValueOf[PoolModelMembersIPAddress] `tfsdk:"targets"`
	}

	PoolModelMembersIPAddress struct {
		Enabled   supertypes.BoolValue   `tfsdk:"enabled"`
		IPAddress supertypes.StringValue `tfsdk:"ip_address"`
		Port      supertypes.Int64Value  `tfsdk:"port"`
		Ratio     supertypes.Int64Value  `tfsdk:"ratio"`
	}

	PoolModelHealth struct {
		PassiveMonitoringEnabled supertypes.BoolValue           `tfsdk:"passive_monitoring_enabled"`
		Monitors                 supertypes.ListValueOf[string] `tfsdk:"monitors"`
	}

	PoolModelTLS struct {
		Enabled                supertypes.BoolValue           `tfsdk:"enabled"`
		DomainNames            supertypes.ListValueOf[string] `tfsdk:"domain_names"`
		CaCertificateRefs      supertypes.ListValueOf[string] `tfsdk:"ca_certificate_refs"`
		CommonNameCheckEnabled supertypes.BoolValue           `tfsdk:"common_name_check_enabled"`
	}

	PoolModelPersistence struct {
		Type  supertypes.StringValue `tfsdk:"type"`
		Value supertypes.StringValue `tfsdk:"value"`
	}
)

func (rm *PoolModel) Copy() *PoolModel {
	x := &PoolModel{}
	utils.ModelCopy(rm, x)
	return x
}

// ToSDKPoolGroupModel converts the model to the SDK model.
func (rm *PoolModel) ToSDKPoolModelRequest(ctx context.Context, _ *client.CloudAvenue) (*edgeloadbalancer.PoolModelRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	pool := &edgeloadbalancer.PoolModelRequest{
		Name:        rm.Name.Get(),
		Description: rm.Description.Get(),
		Enabled:     rm.Enabled.GetPtr(),
		Algorithm:   edgeloadbalancer.PoolAlgorithm(rm.Algorithm.Get()),
		DefaultPort: rm.DefaultPort.GetIntPtr(),
		GatewayRef:  govcdtypes.OpenApiReference{ID: rm.EdgeGatewayID.Get(), Name: rm.EdgeGatewayName.Get()},
	}

	if rm.Members.IsKnown() {
		x, d := rm.Members.Get(ctx)
		diags = append(diags, d...)
		if diags.HasError() {
			return nil, diags
		}

		if x.GracefulTimeoutPeriod.IsKnown() {
			i, err := strconv.Atoi(x.GracefulTimeoutPeriod.Get())
			if err != nil {
				diags.AddError("Error converting GracefulTimeoutPeriod to int", err.Error())
				return nil, diags
			}
			pool.GracefulTimeoutPeriod = &i
		}

		if x.TargetGroup.IsKnown() {
			pool.MemberGroupRef = &govcdtypes.OpenApiReference{
				ID: x.TargetGroup.Get(),
			}
		}

		if x.Targets.IsKnown() {
			ipAddrs, d := x.Targets.Get(ctx)
			diags = append(diags, d...)
			if diags.HasError() {
				return nil, diags
			}

			members := make([]edgeloadbalancer.PoolModelMember, 0)
			for _, m := range ipAddrs {
				members = append(members, edgeloadbalancer.PoolModelMember{
					Enabled:   m.Enabled.Get(),
					IPAddress: m.IPAddress.Get(),
					Port:      m.Port.GetInt(),
					Ratio:     m.Ratio.GetIntPtr(),
				})
			}
			pool.Members = members
		}
	}

	if rm.Health.IsKnown() {
		x, d := rm.Health.Get(ctx)
		diags = append(diags, d...)
		if diags.HasError() {
			return nil, diags
		}

		pool.PassiveMonitoringEnabled = x.PassiveMonitoringEnabled.GetPtr()
		if x.Monitors.IsKnown() {
			pool.HealthMonitors = make([]edgeloadbalancer.PoolModelHealthMonitor, 0)
			monitors, d := x.Monitors.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			for _, m := range monitors {
				pool.HealthMonitors = append(pool.HealthMonitors, edgeloadbalancer.PoolModelHealthMonitor{
					Type: edgeloadbalancer.PoolHealthMonitorType(m),
				})
			}
		}
	}

	if rm.TLS.IsKnown() {
		x, d := rm.TLS.Get(ctx)
		diags = append(diags, d...)
		if diags.HasError() {
			return nil, diags
		}

		pool.SSLEnabled = x.Enabled.GetPtr()
		pool.CommonNameCheckEnabled = x.CommonNameCheckEnabled.GetPtr()
		if x.DomainNames.IsKnown() {
			pool.DomainNames, d = x.DomainNames.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}
		}

		if x.CaCertificateRefs.IsKnown() {
			refs, d := x.CaCertificateRefs.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}

			for _, ref := range refs {
				pool.CaCertificateRefs = append(pool.CaCertificateRefs, govcdtypes.OpenApiReference{
					ID: ref,
				})
			}
		}
	}

	if rm.Persistence.IsKnown() {
		x, d := rm.Persistence.Get(ctx)
		diags = append(diags, d...)
		if diags.HasError() {
			return nil, diags
		}

		pool.PersistenceProfile = &edgeloadbalancer.PoolModelPersistenceProfile{
			Type:  edgeloadbalancer.PoolPersistenceProfileType(x.Type.Get()),
			Value: x.Value.Get(),
		}
	}

	return pool, diags
}
