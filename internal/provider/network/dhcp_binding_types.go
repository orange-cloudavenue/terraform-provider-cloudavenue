package network

import (
	"context"
	"fmt"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type DHCPBindingModel struct {
	Description  supertypes.StringValue       `tfsdk:"description"`
	DhcpV4Config supertypes.SingleNestedValue `tfsdk:"dhcp_v4_config"`
	DNSServers   supertypes.ListValue         `tfsdk:"dns_servers"`
	ID           supertypes.StringValue       `tfsdk:"id"`
	IPAddress    supertypes.StringValue       `tfsdk:"ip_address"`
	LeaseTime    supertypes.Int64Value        `tfsdk:"lease_time"`
	MacAddress   supertypes.StringValue       `tfsdk:"mac_address"`
	Name         supertypes.StringValue       `tfsdk:"name"`
	OrgNetworkID supertypes.StringValue       `tfsdk:"org_network_id"`
}

// * DhcpV4Config.
type DHCPBindingModelDhcpV4Config struct {
	GatewayAddress supertypes.StringValue `tfsdk:"gateway_address"`
	Hostname       supertypes.StringValue `tfsdk:"hostname"`
}

type DHCPBindingModelDNSServers []supertypes.StringValue

func NewDhcpBinding(t any) *DHCPBindingModel {
	switch x := t.(type) {
	case tfsdk.State:
		return &DHCPBindingModel{
			Description: supertypes.NewStringNull(),

			DhcpV4Config: supertypes.NewSingleNestedNull(x.Schema.GetAttributes()["dhcp_v_4_config"].GetType().(supertypes.SingleNestedType).AttributeTypes()),
			DNSServers:   supertypes.NewListNull(x.Schema.GetAttributes()["dns_servers"].GetType().(supertypes.ListType).ElementType()),
			ID:           supertypes.NewStringUnknown(),

			IPAddress: supertypes.NewStringNull(),

			LeaseTime: supertypes.NewInt64Unknown(),

			MacAddress: supertypes.NewStringNull(),

			Name: supertypes.NewStringNull(),

			OrgNetworkID: supertypes.NewStringNull(),
		}

	case tfsdk.Plan:
		return &DHCPBindingModel{
			Description: supertypes.NewStringNull(),

			DhcpV4Config: supertypes.NewSingleNestedNull(x.Schema.GetAttributes()["dhcp_v_4_config"].GetType().(supertypes.SingleNestedType).AttributeTypes()),
			DNSServers:   supertypes.NewListNull(x.Schema.GetAttributes()["dns_servers"].GetType().(supertypes.ListType).ElementType()),
			ID:           supertypes.NewStringUnknown(),

			IPAddress: supertypes.NewStringNull(),

			LeaseTime: supertypes.NewInt64Unknown(),

			MacAddress: supertypes.NewStringNull(),

			Name: supertypes.NewStringNull(),

			OrgNetworkID: supertypes.NewStringNull(),
		}

	case tfsdk.Config:
		return &DHCPBindingModel{
			Description: supertypes.NewStringNull(),

			DhcpV4Config: supertypes.NewSingleNestedNull(x.Schema.GetAttributes()["dhcp_v_4_config"].GetType().(supertypes.SingleNestedType).AttributeTypes()),
			DNSServers:   supertypes.NewListNull(x.Schema.GetAttributes()["dns_servers"].GetType().(supertypes.ListType).ElementType()),
			ID:           supertypes.NewStringUnknown(),

			IPAddress: supertypes.NewStringNull(),

			LeaseTime: supertypes.NewInt64Unknown(),

			MacAddress: supertypes.NewStringNull(),

			Name: supertypes.NewStringNull(),

			OrgNetworkID: supertypes.NewStringNull(),
		}

	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *DHCPBindingModel) Copy() *DHCPBindingModel {
	x := &DHCPBindingModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetDhcpV4Config returns the value of the DhcpV4Config field.
func (rm *DHCPBindingModel) GetDhcpV4Config(ctx context.Context) (values DHCPBindingModelDhcpV4Config, diags diag.Diagnostics) {
	values = DHCPBindingModelDhcpV4Config{}
	d := rm.DhcpV4Config.Get(ctx, &values, basetypes.ObjectAsOptions{})
	return values, d
}

// GetDNSServers returns the value of the DNSServers field.
func (rm *DHCPBindingModel) GetDNSServers(ctx context.Context) (values DHCPBindingModelDNSServers, diags diag.Diagnostics) {
	values = make(DHCPBindingModelDNSServers, 0)
	d := rm.DNSServers.Get(ctx, &values, false)
	return values, d
}

// Get returns the value of the given attribute.
func (ds DHCPBindingModelDNSServers) Get() []string {
	return utils.SuperSliceStringToSliceString(ds)
}

// ToNetworkDhcpBindingType converts a DHCPBindingModel to govcdtypes.OpenApiOrgVdcNetworkDhcpBinding.
func (rm *DHCPBindingModel) ToNetworkDhcpBindingType(ctx context.Context) (values *govcdtypes.OpenApiOrgVdcNetworkDhcpBinding, diags diag.Diagnostics) {
	values = &govcdtypes.OpenApiOrgVdcNetworkDhcpBinding{
		Name:        rm.Name.Get(),
		Description: rm.Description.Get(),
		BindingType: "IPV4",
		MacAddress:  rm.MacAddress.Get(),
		IpAddress:   rm.IPAddress.Get(),
		LeaseTime:   utils.TakeIntPointer(int(rm.LeaseTime.Get())),
	}

	if rm.ID.IsKnown() {
		values.ID = rm.ID.Get()
	}

	if rm.DNSServers.IsKnown() {
		dnsServers, d := rm.GetDNSServers(ctx)
		if d.HasError() {
			return nil, d
		}

		values.DnsServers = dnsServers.Get()
	}

	if rm.DhcpV4Config.IsKnown() {
		dhcpConfig, d := rm.GetDhcpV4Config(ctx)
		if d.HasError() {
			return nil, d
		}

		values.DhcpV4BindingConfig = &govcdtypes.DhcpV4BindingConfig{
			GatewayIPAddress: dhcpConfig.GatewayAddress.Get(),
			HostName:         dhcpConfig.Hostname.Get(),
		}
	}

	return values, diags
}
