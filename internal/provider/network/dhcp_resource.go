// Package network provides a Terraform resource.
package network

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dhcpResource{}
	_ resource.ResourceWithConfigure   = &dhcpResource{}
	_ resource.ResourceWithImportState = &dhcpResource{}
)

// NewDhcpResource is a helper function to simplify the provider implementation.
func NewDhcpResource() resource.Resource {
	return &dhcpResource{}
}

// dhcpResource is the resource implementation.
type dhcpResource struct {
	client *client.CloudAvenue
	org    org.Org
}

// Init Initializes the resource.
func (r *dhcpResource) Init(ctx context.Context, rm *dhcpModel) (diags diag.Diagnostics) {
	r.org, diags = org.Init(r.client)
	return
}

// Metadata returns the resource type name.
func (r *dhcpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_dhcp"
}

// Schema defines the schema for the resource.
func (r *dhcpResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = dhcpSchema(ctx).GetResource(ctx)
}

func (r *dhcpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dhcpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &dhcpModel{}

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

	mutex.GlobalMutex.KvLock(ctx, plan.OrgNetworkID.ValueString())
	defer mutex.GlobalMutex.KvUnlock(ctx, plan.OrgNetworkID.ValueString())

	resp.Diagnostics.Append(r.createUpdateDHCP(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, found, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.Diagnostics.AddError("DHCP not found after creation", "After creating the DHCP, the API returned entity not found")
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *dhcpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &dhcpModel{}

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

	mutex.GlobalMutex.KvLock(ctx, state.OrgNetworkID.ValueString())
	defer mutex.GlobalMutex.KvUnlock(ctx, state.OrgNetworkID.ValueString())

	stateRefreshed, found, d := r.read(ctx, state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dhcpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan  = &dhcpModel{}
		state = &dhcpModel{}
	)

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
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

	mutex.GlobalMutex.KvLock(ctx, plan.OrgNetworkID.ValueString())
	defer mutex.GlobalMutex.KvUnlock(ctx, plan.OrgNetworkID.ValueString())

	resp.Diagnostics.Append(r.createUpdateDHCP(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, found, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dhcpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &dhcpModel{}

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

	mutex.GlobalMutex.KvLock(ctx, state.OrgNetworkID.ValueString())
	defer mutex.GlobalMutex.KvUnlock(ctx, state.OrgNetworkID.ValueString())

	if err := r.org.DeleteNetworkDHCP(state.OrgNetworkID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting dhcp", err.Error())
	}
}

func (r *dhcpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_network_id"), req.ID)...)
}

// createUpdateDhcp The dhcp has no create method in the API, so we use the update method.
func (r *dhcpResource) createUpdateDHCP(ctx context.Context, rm *dhcpModel) (diags diag.Diagnostics) {
	if err := r.org.UpdateNetworkDHCP(rm.OrgNetworkID.ValueString(), rm.toNetworkDHCP(ctx)); err != nil {
		diags.AddError("Error updating dhcp", err.Error())
	}

	return
}

// toNetworkDHCP converts a dhcp to a govcdtypes.OpenApiOrgVdcNetworkDhcp object.
func (rm *dhcpModel) toNetworkDHCP(ctx context.Context) *govcdtypes.OpenApiOrgVdcNetworkDhcp {
	object := &govcdtypes.OpenApiOrgVdcNetworkDhcp{
		Mode:      rm.Mode.ValueString(),
		LeaseTime: utils.TakeIntPointer(int(rm.LeaseTime.ValueInt64())),
		IPAddress: rm.ListenerIPAddress.ValueString(),
	}

	if !rm.Pools.IsNull() && !rm.Pools.IsUnknown() {
		pools, err := rm.PoolsFromPlan(ctx)
		if err != nil {
			return nil
		}

		object.DhcpPools = make([]govcdtypes.OpenApiOrgVdcNetworkDhcpPools, len(pools))

		for i, pool := range pools {
			object.DhcpPools[i] = govcdtypes.OpenApiOrgVdcNetworkDhcpPools{
				IPRange: govcdtypes.OpenApiOrgVdcNetworkDhcpIpRange{
					StartAddress: pool.Start.ValueString(),
					EndAddress:   pool.End.ValueString(),
				},
			}
		}
	}

	if !rm.DNSServers.IsNull() && !rm.DNSServers.IsUnknown() {
		dnsServers, err := rm.DNSServersFromPlan(ctx)
		if err != nil {
			return nil
		}

		object.DnsServers = dnsServers
	}

	return object
}

// read() reads the resource from the API and updates the state with it.
func (r *dhcpResource) read(ctx context.Context, plan *dhcpModel) (state *dhcpModel, found bool, diags diag.Diagnostics) {
	var d diag.Diagnostics

	state = new(dhcpModel)
	found = true

	orgNetworkDhcp, err := r.org.GetNetworkDHCP(plan.OrgNetworkID.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			found = false
			return
		}
		diags.AddError("Error getting org network dhcp", err.Error())
		return
	}

	// DHCP resource don't have an ID, so we use the OrgNetworkID
	state.ID = types.StringValue(plan.OrgNetworkID.ValueString())
	state.OrgNetworkID = types.StringValue(plan.OrgNetworkID.ValueString())
	state.Mode = types.StringValue(orgNetworkDhcp.OpenApiOrgVdcNetworkDhcp.Mode)
	state.LeaseTime = types.Int64Value(int64(*orgNetworkDhcp.OpenApiOrgVdcNetworkDhcp.LeaseTime))
	state.ListenerIPAddress = utils.StringValueOrNull(orgNetworkDhcp.OpenApiOrgVdcNetworkDhcp.IPAddress)

	if orgNetworkDhcp.OpenApiOrgVdcNetworkDhcp.DhcpPools != nil {
		pools := make(dhcpModelPools, len(orgNetworkDhcp.OpenApiOrgVdcNetworkDhcp.DhcpPools))
		for i, pool := range orgNetworkDhcp.OpenApiOrgVdcNetworkDhcp.DhcpPools {
			pools[i] = dhcpModelPool{
				Start: types.StringValue(pool.IPRange.StartAddress),
				End:   types.StringValue(pool.IPRange.EndAddress),
			}
		}
		state.Pools, d = pools.ToPlan(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
	}

	if orgNetworkDhcp.OpenApiOrgVdcNetworkDhcp.DnsServers != nil {
		dnsServers := make(dhcpModelDNSServers, len(orgNetworkDhcp.OpenApiOrgVdcNetworkDhcp.DnsServers))
		copy(dnsServers, orgNetworkDhcp.OpenApiOrgVdcNetworkDhcp.DnsServers)

		state.DNSServers, d = dnsServers.ToPlan(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
	}

	return
}
