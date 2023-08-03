package edgegw

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

type VPNIPSecModel struct {
	// AuthenticationMode supertypes.StringValue       `tfsdk:"authentication_mode"`
	// CACertificateID    supertypes.StringValue       `tfsdk:"ca_certificate_id"`
	// CertificateID      supertypes.StringValue       `tfsdk:"certificate_id"`
	Description     supertypes.StringValue       `tfsdk:"description"`
	EdgeGatewayID   supertypes.StringValue       `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue       `tfsdk:"edge_gateway_name"`
	Enabled         supertypes.BoolValue         `tfsdk:"enabled"`
	ID              supertypes.StringValue       `tfsdk:"id"`
	LocalIPAddress  supertypes.StringValue       `tfsdk:"local_ip_address"`
	LocalNetworks   supertypes.SetValue          `tfsdk:"local_networks"`
	Name            supertypes.StringValue       `tfsdk:"name"`
	PreSharedKey    supertypes.StringValue       `tfsdk:"pre_shared_key"`
	RemoteIPAddress supertypes.StringValue       `tfsdk:"remote_ip_address"`
	RemoteNetworks  supertypes.SetValue          `tfsdk:"remote_networks"`
	SecurityProfile supertypes.SingleNestedValue `tfsdk:"security_profile"`
	SecurityType    supertypes.StringValue       `tfsdk:"security_type"`
}

type VPNIPSecModelLocalNetworks []supertypes.StringValue

type VPNIPSecModelRemoteNetworks []supertypes.StringValue

// * SecurityProfile.
type VPNIPSecModelSecurityProfile struct {
	IkeDhGroups                supertypes.StringValue `tfsdk:"ike_dh_groups"`
	IkeDigestAlgorithm         supertypes.StringValue `tfsdk:"ike_digest_algorithm"`
	IkeEncryptionAlgorithm     supertypes.StringValue `tfsdk:"ike_encryption_algorithm"`
	IkeSaLifetime              supertypes.Int64Value  `tfsdk:"ike_sa_lifetime"`
	IkeVersion                 supertypes.StringValue `tfsdk:"ike_version"`
	TunnelDfPolicy             supertypes.StringValue `tfsdk:"tunnel_df_policy"`
	TunnelDhGroups             supertypes.StringValue `tfsdk:"tunnel_dh_groups"`
	TunnelDigestAlgorithms     supertypes.StringValue `tfsdk:"tunnel_digest_algorithms"`
	TunnelDpd                  supertypes.Int64Value  `tfsdk:"tunnel_dpd"`
	TunnelEncryptionAlgorithms supertypes.StringValue `tfsdk:"tunnel_encryption_algorithms"`
	TunnelPfs                  supertypes.BoolValue   `tfsdk:"tunnel_pfs"`
	TunnelSaLifetime           supertypes.Int64Value  `tfsdk:"tunnel_sa_lifetime"`
}

const (
	vpnAuthentication string = "PSK"
	vpnModeInit       string = "INITIATOR"
	profileDefault    string = "DEFAULT"
	profileCustom     string = "CUSTOM"
)

func NewVPNIPSec(t any) *VPNIPSecModel {
	switch x := t.(type) {
	case tfsdk.State: //nolint:dupl
		return &VPNIPSecModel{
			// AuthenticationMode: supertypes.NewStringUnknown(),
			// CACertificateID:    supertypes.NewStringNull(),
			// CertificateID:      supertypes.NewStringNull(),
			Description:     supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			Enabled:         supertypes.NewBoolUnknown(),
			ID:              supertypes.NewStringUnknown(),
			LocalIPAddress:  supertypes.NewStringNull(),
			LocalNetworks:   supertypes.NewSetNull(x.Schema.GetAttributes()["local_networks"].GetType().(supertypes.SetType).ElementType()),
			Name:            supertypes.NewStringNull(),
			PreSharedKey:    supertypes.NewStringNull(),
			RemoteIPAddress: supertypes.NewStringNull(),
			RemoteNetworks:  supertypes.NewSetNull(x.Schema.GetAttributes()["remote_networks"].GetType().(supertypes.SetType).ElementType()),
			SecurityProfile: supertypes.NewSingleNestedNull(x.Schema.GetAttributes()["security_profile"].GetType().(supertypes.SingleNestedType).AttributeTypes()),
			SecurityType:    supertypes.NewStringUnknown(),
		}

	case tfsdk.Plan: //nolint:dupl
		return &VPNIPSecModel{
			// AuthenticationMode: supertypes.NewStringUnknown(),
			// CACertificateID:    supertypes.NewStringNull(),
			// CertificateID:      supertypes.NewStringNull(),
			Description:     supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			Enabled:         supertypes.NewBoolUnknown(),
			ID:              supertypes.NewStringUnknown(),
			LocalIPAddress:  supertypes.NewStringNull(),
			LocalNetworks:   supertypes.NewSetNull(x.Schema.GetAttributes()["local_networks"].GetType().(supertypes.SetType).ElementType()),
			Name:            supertypes.NewStringNull(),
			PreSharedKey:    supertypes.NewStringNull(),
			RemoteIPAddress: supertypes.NewStringNull(),
			RemoteNetworks:  supertypes.NewSetNull(x.Schema.GetAttributes()["remote_networks"].GetType().(supertypes.SetType).ElementType()),
			SecurityProfile: supertypes.NewSingleNestedNull(x.Schema.GetAttributes()["security_profile"].GetType().(supertypes.SingleNestedType).AttributeTypes()),
			SecurityType:    supertypes.NewStringUnknown(),
		}

	case tfsdk.Config: //nolint:dupl
		return &VPNIPSecModel{
			// AuthenticationMode: supertypes.NewStringUnknown(),
			// CACertificateID:    supertypes.NewStringNull(),
			// CertificateID:      supertypes.NewStringNull(),
			Description:     supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			Enabled:         supertypes.NewBoolUnknown(),
			ID:              supertypes.NewStringUnknown(),
			LocalIPAddress:  supertypes.NewStringNull(),
			LocalNetworks:   supertypes.NewSetNull(x.Schema.GetAttributes()["local_networks"].GetType().(supertypes.SetType).ElementType()),
			Name:            supertypes.NewStringNull(),
			PreSharedKey:    supertypes.NewStringNull(),
			RemoteIPAddress: supertypes.NewStringNull(),
			RemoteNetworks:  supertypes.NewSetNull(x.Schema.GetAttributes()["remote_networks"].GetType().(supertypes.SetType).ElementType()),
			SecurityProfile: supertypes.NewSingleNestedNull(x.Schema.GetAttributes()["security_profile"].GetType().(supertypes.SingleNestedType).AttributeTypes()),
			SecurityType:    supertypes.NewStringUnknown(),
		}

	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *VPNIPSecModel) Copy() *VPNIPSecModel {
	x := &VPNIPSecModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetLocalNetworks returns the value of the LocalNetworks field.
func (rm *VPNIPSecModel) GetLocalNetworks(ctx context.Context) (values VPNIPSecModelLocalNetworks, diags diag.Diagnostics) {
	values = make(VPNIPSecModelLocalNetworks, 0)
	d := rm.LocalNetworks.Get(ctx, &values, false)
	return values, d
}

// GetRemoteNetworks returns the value of the RemoteNetworks field.
func (rm *VPNIPSecModel) GetRemoteNetworks(ctx context.Context) (values VPNIPSecModelRemoteNetworks, diags diag.Diagnostics) {
	values = make(VPNIPSecModelRemoteNetworks, 0)
	d := rm.RemoteNetworks.Get(ctx, &values, false)
	return values, d
}

// GetSecurityProfile returns the value of the SecurityProfile field.
func (rm *VPNIPSecModel) GetSecurityProfile(ctx context.Context) (values VPNIPSecModelSecurityProfile, diags diag.Diagnostics) {
	values = VPNIPSecModelSecurityProfile{}
	d := rm.SecurityProfile.Get(ctx, &values, basetypes.ObjectAsOptions{})
	return values, d
}

func (rm *VPNIPSecModel) GetNsxtIPSecVPNTunnelSecurityProfile(ctx context.Context) (values *govcdtypes.NsxtIpSecVpnTunnelSecurityProfile, diags diag.Diagnostics) {
	if !rm.SecurityProfile.IsKnown() {
		return &govcdtypes.NsxtIpSecVpnTunnelSecurityProfile{}, nil
	}
	values = &govcdtypes.NsxtIpSecVpnTunnelSecurityProfile{}
	securityProfile, d := rm.GetSecurityProfile(ctx)
	if d.HasError() {
		return &govcdtypes.NsxtIpSecVpnTunnelSecurityProfile{}, d
	}

	if securityProfile.IkeDhGroups.IsKnown() {
		values.IkeConfiguration.DhGroups = append(values.IkeConfiguration.DhGroups, securityProfile.IkeDhGroups.Get())
	}
	if securityProfile.IkeDigestAlgorithm.IsKnown() {
		values.IkeConfiguration.DigestAlgorithms = append(values.IkeConfiguration.DigestAlgorithms, securityProfile.IkeDigestAlgorithm.Get())
	}
	if securityProfile.IkeEncryptionAlgorithm.IsKnown() {
		values.IkeConfiguration.EncryptionAlgorithms = append(values.IkeConfiguration.EncryptionAlgorithms, securityProfile.IkeEncryptionAlgorithm.Get())
	}
	if securityProfile.IkeSaLifetime.IsKnown() {
		values.IkeConfiguration.SaLifeTime = utils.TakeIntPointer(int(securityProfile.IkeSaLifetime.Get()))
	}
	if securityProfile.IkeVersion.IsKnown() {
		values.IkeConfiguration.IkeVersion = securityProfile.IkeVersion.Get()
	}
	if securityProfile.TunnelDfPolicy.IsKnown() {
		values.TunnelConfiguration.DfPolicy = securityProfile.TunnelDfPolicy.Get()
	}
	if securityProfile.TunnelDhGroups.IsKnown() {
		values.TunnelConfiguration.DhGroups = append(values.TunnelConfiguration.DhGroups, securityProfile.TunnelDhGroups.Get())
	}
	if securityProfile.TunnelDigestAlgorithms.IsKnown() {
		values.TunnelConfiguration.DigestAlgorithms = append(values.TunnelConfiguration.DigestAlgorithms, securityProfile.TunnelDigestAlgorithms.Get())
	}
	if securityProfile.TunnelEncryptionAlgorithms.IsKnown() {
		values.TunnelConfiguration.EncryptionAlgorithms = append(values.TunnelConfiguration.EncryptionAlgorithms, securityProfile.TunnelEncryptionAlgorithms.Get())
	}
	if securityProfile.TunnelPfs.IsKnown() {
		values.TunnelConfiguration.PerfectForwardSecrecyEnabled = securityProfile.TunnelPfs.Get()
	}
	if securityProfile.TunnelSaLifetime.IsKnown() {
		values.TunnelConfiguration.SaLifeTime = utils.TakeIntPointer(int(securityProfile.TunnelSaLifetime.Get()))
	}
	if securityProfile.TunnelDpd.IsKnown() {
		values.DpdConfiguration.ProbeInterval = int(securityProfile.TunnelDpd.Get())
	}

	return
}

func (rm VPNIPSecModelLocalNetworks) Get() []string {
	return utils.SuperSliceTypesStringToSliceString(rm)
}

func (rm VPNIPSecModelRemoteNetworks) Get() []string {
	return utils.SuperSliceTypesStringToSliceString(rm)
}

func (rm *VPNIPSecModel) ToNsxtIPSecVPNTunnel(ctx context.Context) (values *govcdtypes.NsxtIpSecVpnTunnel, err error) {
	values = &govcdtypes.NsxtIpSecVpnTunnel{
		Name:                    rm.Name.Get(),
		Description:             rm.Description.Get(),
		Enabled:                 rm.Enabled.Get(),
		PreSharedKey:            rm.PreSharedKey.Get(),
		Logging:                 false,             // not available on cloudavenue
		AuthenticationMode:      vpnAuthentication, // Only PSK is supported on cloudavenue
		ConnectorInitiationMode: vpnModeInit,
	}

	if rm.ID.IsKnown() {
		values.ID = rm.ID.Get()
	}

	// Get local networks
	localNet, diags := rm.GetLocalNetworks(ctx)
	if diags.HasError() {
		return
	}
	values.LocalEndpoint.LocalId = rm.LocalIPAddress.Get()
	values.LocalEndpoint.LocalAddress = rm.LocalIPAddress.Get()
	values.LocalEndpoint.LocalNetworks = localNet.Get()

	// Get remote networks
	remoteNet, diags := rm.GetRemoteNetworks(ctx)
	if diags.HasError() {
		return
	}
	values.RemoteEndpoint.RemoteId = rm.RemoteIPAddress.Get()
	values.RemoteEndpoint.RemoteAddress = rm.RemoteIPAddress.Get()
	values.RemoteEndpoint.RemoteNetworks = remoteNet.Get()

	return
}
