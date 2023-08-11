// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vpnIPSecResource{}
	_ resource.ResourceWithConfigure   = &vpnIPSecResource{}
	_ resource.ResourceWithImportState = &vpnIPSecResource{}
)

// NewVpnIpsecResource is a helper function to simplify the provider implementation.
func NewVPNIPSecResource() resource.Resource {
	return &vpnIPSecResource{}
}

// vpnIPSecResource is the resource implementation.
type vpnIPSecResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the resource.
func (r *vpnIPSecResource) Init(ctx context.Context, rm *VPNIPSecModel) (diags diag.Diagnostics) {
	var err error

	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(rm.EdgeGatewayID.Get()),
		Name: types.StringValue(rm.EdgeGatewayName.Get()),
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	return
}

// Metadata returns the resource type name.
func (r *vpnIPSecResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_vpn_ipsec"
}

// Schema defines the schema for the resource.
func (r *vpnIPSecResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vpnIPSecSchema(ctx).GetResource(ctx)
}

func (r *vpnIPSecResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpnIPSecResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &VPNIPSecModel{}

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource creation logic here.
	*/

	// Lock object EdgeGateway
	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}
	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}

	// Get data from plan
	ipSecVPNConfig, d := plan.ToNsxtIPSecVPNTunnel(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create VPN Config with default profile settings
	createdIPSecVPNConfig, err := r.edgegw.NsxtEdgeGateway.CreateIpSecVpnTunnel(ipSecVPNConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error creating IPsec VPN Tunnel configuration", err.Error())
		return
	}

	// Set ID from created VPN
	plan.ID.Set(createdIPSecVPNConfig.NsxtIpSecVpn.ID)

	// Security Type is Set to DEFAULT in stateRefreshed
	// Check if Tunnel Profile has custom settings and apply them
	if plan.SecurityProfile.IsKnown() {
		// Get Tunnel Profile from Plan
		vpnTunnelSecProfile, d := plan.GetNsxtIPSecVPNTunnelSecurityProfile(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Version increment
		(*createdIPSecVPNConfig.NsxtIpSecVpn.Version.Version)++

		// Set Security Type to CUSTOM
		vpnTunnelSecProfile.SecurityType = profileCustom

		// Update Tunnel Profile with custom settings
		if _, err := createdIPSecVPNConfig.UpdateTunnelConnectionProperties(vpnTunnelSecProfile); err != nil {
			resp.Diagnostics.AddError("Error updating IPsec VPN Tunnel Security Profile", err.Error())
			return
		}
	}
	// refresh state
	stateRefreshed, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vpnIPSecResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &VPNIPSecModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource read here
	*/

	// read IP Sec VPN and update state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpnIPSecResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan  = &VPNIPSecModel{}
		state = &VPNIPSecModel{}
	)

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource update here
	*/

	// Lock object EdgeGateway
	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}
	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}

	// Get VPN Config from plan (without profile settings)
	planVPNTunnel, d := plan.ToNsxtIPSecVPNTunnel(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Profile from plan
	planVPNIPSec, d := plan.GetNsxtIPSecVPNTunnelSecurityProfile(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get VPN from API
	existingIPSecVPNConfiguration, err := r.edgegw.NsxtEdgeGateway.GetIpSecVpnTunnelById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error Retrieving IPsec VPN Tunnel: %s", err.Error())
		return
	}

	// Update VPN Tunnel
	if _, err = existingIPSecVPNConfiguration.Update(planVPNTunnel); err != nil {
		resp.Diagnostics.AddError("Error updating VPN Tunnel configuration", err.Error())
		return
	}
	if !plan.SecurityProfile.IsUnknown() {
		// Version +1
		(*existingIPSecVPNConfiguration.NsxtIpSecVpn.Version.Version)++

		// Force SecurityType to CUSTOM
		planVPNIPSec.SecurityType = profileCustom

		// Update VPN Tunnel with CUSTOM Profile IPsec settings
		if _, err = existingIPSecVPNConfiguration.UpdateTunnelConnectionProperties(planVPNIPSec); err != nil {
			resp.Diagnostics.AddError("Error updating VPN Tunnel configuration", err.Error())
			return
		}
	}

	// read Update Config and update state
	stateRefreshed, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpnIPSecResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &VPNIPSecModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource deletion here
	*/
	// Lock object EdgeGateway
	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}
	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}

	// Get VPN
	ipSecVPNConfig, err := r.edgegw.NsxtEdgeGateway.GetIpSecVpnTunnelById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error Retrieving IPsec VPN Tunnel: %s", err.Error())
		return
	}

	// Delete VPN
	if err = ipSecVPNConfig.Delete(); err != nil {
		resp.Diagnostics.AddError("Error Deleting IPsec VPN Tunnel configuration: %s", err.Error())
		return
	}
}

func (r *vpnIPSecResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var (
		edgegwID, edgegwName string
		d                    diag.Diagnostics
		err                  error
		vpnIPSec             *govcd.NsxtIpSecVpnTunnel
	)

	// Split req.ID with dot. ID format is EdgeGatewayIDOrName.VPNIPSecIDOrName
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError("Invalid ID format", "ID format is EdgeGatewayIDOrName.VPNIPSecIDOrName. If Several Name are the same, please use ID instead")
		return
	}

	// Get Org to retrieve EdgeGateway
	r.org, d = org.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Get EdgeGW is ID or Name
	if uuid.IsEdgeGateway(idParts[0]) {
		edgegwID = idParts[0]
	} else {
		edgegwName = idParts[0]
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(edgegwID),
		Name: types.StringValue(edgegwName),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to import VPN IPSec.", err.Error())
		return
	}

	// Get VPN IPSec
	vpnIPSec, err = r.edgegw.GetIpSecVpnTunnelByName(idParts[1])
	if govcd.ContainsNotFound(err) {
		vpnIPSec, err = r.edgegw.GetIpSecVpnTunnelById(idParts[1])
	}

	if err != nil {
		allRules, err := r.edgegw.GetAllIpSecVpnTunnels(nil)
		if err != nil {
			resp.Diagnostics.AddError("Failed to Get ALL VPN IPSec.", err.Error())
			return
		}
		listStr := getVPNIPSecTunnelsList(idParts[1], allRules)
		resp.Diagnostics.AddError("Failed to Get VPN IPSec "+listStr, err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), vpnIPSec.NsxtIpSecVpn.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), vpnIPSec.NsxtIpSecVpn.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_id"), r.edgegw.GetID())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), r.edgegw.GetName())...)
}

func (r *vpnIPSecResource) read(ctx context.Context, planOrState *VPNIPSecModel) (stateRefreshed *VPNIPSecModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	// Get IpSecVpnTunnel by Name or ID
	var (
		vpnTunnel           *govcd.NsxtIpSecVpnTunnel
		vpnTunnelSecProfile *VPNIPSecModelSecurityProfile
		err                 error
	)
	if stateRefreshed.ID.IsKnown() {
		vpnTunnel, err = r.edgegw.NsxtEdgeGateway.GetIpSecVpnTunnelById(stateRefreshed.ID.Get())
	} else {
		vpnTunnel, err = r.edgegw.NsxtEdgeGateway.GetIpSecVpnTunnelByName(stateRefreshed.Name.Get())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, diags
		}
		diags.AddError("Error retrieving VPN Tunnel", err.Error())
		return nil, true, diags
	}

	stateRefreshed.Description.Set(vpnTunnel.NsxtIpSecVpn.Description)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())
	stateRefreshed.Enabled.Set(vpnTunnel.NsxtIpSecVpn.Enabled)
	stateRefreshed.ID.Set(vpnTunnel.NsxtIpSecVpn.ID)
	stateRefreshed.LocalIPAddress.Set(vpnTunnel.NsxtIpSecVpn.LocalEndpoint.LocalAddress)
	stateRefreshed.LocalNetworks.Set(ctx, vpnTunnel.NsxtIpSecVpn.LocalEndpoint.LocalNetworks)
	stateRefreshed.Name.Set(vpnTunnel.NsxtIpSecVpn.Name)
	stateRefreshed.PreSharedKey.Set(vpnTunnel.NsxtIpSecVpn.PreSharedKey)
	stateRefreshed.RemoteIPAddress.Set(vpnTunnel.NsxtIpSecVpn.RemoteEndpoint.RemoteAddress)
	stateRefreshed.RemoteNetworks.Set(ctx, vpnTunnel.NsxtIpSecVpn.RemoteEndpoint.RemoteNetworks)
	stateRefreshed.SecurityType.Set(vpnTunnel.NsxtIpSecVpn.SecurityType)

	// Get IPSec Profile
	vpnTunnelSecProfile = &VPNIPSecModelSecurityProfile{
		IkeDhGroups:                supertypes.NewStringNull(),
		IkeDigestAlgorithm:         supertypes.NewStringNull(),
		IkeEncryptionAlgorithm:     supertypes.NewStringNull(),
		IkeSaLifetime:              supertypes.NewInt64Null(),
		IkeVersion:                 supertypes.NewStringNull(),
		TunnelDfPolicy:             supertypes.NewStringNull(),
		TunnelDhGroups:             supertypes.NewStringNull(),
		TunnelDigestAlgorithms:     supertypes.NewStringNull(),
		TunnelDpd:                  supertypes.NewInt64Null(),
		TunnelEncryptionAlgorithms: supertypes.NewStringNull(),
		TunnelPfs:                  supertypes.NewBoolNull(),
		TunnelSaLifetime:           supertypes.NewInt64Null(),
	}
	secProfile, err := vpnTunnel.GetTunnelConnectionProperties()
	if err != nil {
		diags.AddError("Error retrieving VPN Tunnel Security Profile", err.Error())
		return nil, true, diags
	}

	// IKE DH Groups
	if len(secProfile.IkeConfiguration.DhGroups) > 0 {
		vpnTunnelSecProfile.IkeDhGroups.Set(secProfile.IkeConfiguration.DhGroups[0])
	}
	// IKE Digest Algorithm
	if len(secProfile.IkeConfiguration.DigestAlgorithms) > 0 {
		vpnTunnelSecProfile.IkeDigestAlgorithm.Set(secProfile.IkeConfiguration.DigestAlgorithms[0])
	}
	// IKE Encryption Algorithm
	if len(secProfile.IkeConfiguration.EncryptionAlgorithms) > 0 {
		vpnTunnelSecProfile.IkeEncryptionAlgorithm.Set(secProfile.IkeConfiguration.EncryptionAlgorithms[0])
	}
	// IKE SA Lifetime
	vpnTunnelSecProfile.IkeSaLifetime.Set(int64(*secProfile.IkeConfiguration.SaLifeTime))
	// IKE Version
	vpnTunnelSecProfile.IkeVersion.Set(secProfile.IkeConfiguration.IkeVersion)
	// IKE Dead Peer Detection
	vpnTunnelSecProfile.TunnelDfPolicy.Set(secProfile.TunnelConfiguration.DfPolicy)
	// Tunnel DH Groups
	if len(secProfile.TunnelConfiguration.DhGroups) > 0 {
		vpnTunnelSecProfile.TunnelDhGroups.Set(secProfile.TunnelConfiguration.DhGroups[0])
	}
	// Tunnel Digest Algorithms
	if len(secProfile.TunnelConfiguration.DigestAlgorithms) > 0 {
		vpnTunnelSecProfile.TunnelDigestAlgorithms.Set(secProfile.TunnelConfiguration.DigestAlgorithms[0])
	}
	// Tunnel Encryption Algorithm
	if len(secProfile.TunnelConfiguration.EncryptionAlgorithms) > 0 {
		vpnTunnelSecProfile.TunnelEncryptionAlgorithms.Set(secProfile.TunnelConfiguration.EncryptionAlgorithms[0])
	}
	// Tunnel SA Lifetime
	vpnTunnelSecProfile.TunnelSaLifetime.Set(int64(*secProfile.TunnelConfiguration.SaLifeTime))
	// Tunnel PFS
	vpnTunnelSecProfile.TunnelPfs.Set(secProfile.TunnelConfiguration.PerfectForwardSecrecyEnabled)
	// Tunnel DPD
	vpnTunnelSecProfile.TunnelDpd.Set(int64(secProfile.DpdConfiguration.ProbeInterval))

	diags.Append(stateRefreshed.SecurityProfile.Set(ctx, vpnTunnelSecProfile)...)
	if diags.HasError() {
		return nil, true, diags
	}

	return stateRefreshed, true, diags
}

func getVPNIPSecTunnelsList(name string, allTunnels []*govcd.NsxtIpSecVpnTunnel) string {
	var list []string
	list = append(list, "\n\tName\t\t\t\tID\t\tLocal IP\tRemote IP\n")
	for _, tunnel := range allTunnels {
		if strings.Contains(tunnel.NsxtIpSecVpn.Name, name) {
			list = append(list, tunnel.NsxtIpSecVpn.Name+"\t"+tunnel.NsxtIpSecVpn.ID+"\t"+tunnel.NsxtIpSecVpn.LocalEndpoint.LocalAddress+"\t"+tunnel.NsxtIpSecVpn.RemoteEndpoint.RemoteAddress+"\n")
		}
	}
	return strings.Join(list, " ")
}
