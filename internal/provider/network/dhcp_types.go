package network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type dhcpModel struct {
	DNSServers        types.List   `tfsdk:"dns_servers"`
	ID                types.String `tfsdk:"id"`
	LeaseTime         types.Int64  `tfsdk:"lease_time"`
	ListenerIPAddress types.String `tfsdk:"listener_ip_address"`
	Mode              types.String `tfsdk:"mode"`
	OrgNetworkID      types.String `tfsdk:"org_network_id"`
	Pools             types.Set    `tfsdk:"pools"`
}

type dhcpModelPools []dhcpModelPool

type dhcpModelPool struct {
	End   types.String `tfsdk:"end_address"`
	Start types.String `tfsdk:"start_address"`
}

type dhcpModelDNSServers []string

// ObjectType() returns the object type for the nested object.
func (p *dhcpModelPools) ObjectType(ctx context.Context) types.ObjectType {
	return types.ObjectType{
		AttrTypes: p.AttrTypes(ctx),
	}
}

// AttrTypes() returns the attribute types for the nested object.
func (p *dhcpModelPools) AttrTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"end_address":   types.StringType,
		"start_address": types.StringType,
	}
}

// ToPlan() returns the plan representation of the nested object.
func (p *dhcpModelPools) ToPlan(ctx context.Context) (basetypes.SetValue, diag.Diagnostics) {
	if p == nil {
		return types.SetNull(p.ObjectType(ctx)), nil
	}

	return types.SetValueFrom(ctx, p.ObjectType(ctx), p)
}

func (rm *dhcpModel) PoolsFromPlan(ctx context.Context) (pools dhcpModelPools, diags diag.Diagnostics) {
	pools = make(dhcpModelPools, 0)
	diags.Append(rm.Pools.ElementsAs(ctx, &pools, false)...)
	if diags.HasError() {
		return
	}

	return pools, diags
}

// DNSServersFromPlan returns the DNSServers from the plan.
func (rm *dhcpModel) DNSServersFromPlan(ctx context.Context) (dnsServers dhcpModelDNSServers, diags diag.Diagnostics) {
	dnsServers = make(dhcpModelDNSServers, 0)
	diags.Append(rm.DNSServers.ElementsAs(ctx, &dnsServers, false)...)
	if diags.HasError() {
		return
	}

	return dnsServers, diags
}

// ToPlan converts a dhcpModelDNSServers to a plan representation.
func (dnsServers *dhcpModelDNSServers) ToPlan(ctx context.Context) (basetypes.ListValue, diag.Diagnostics) {
	return types.ListValueFrom(ctx, types.StringType, dnsServers)
}
