// Package vcda provides a Terraform resource to manage vcda.
package vcda

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vcdaIPResource{}
	_ resource.ResourceWithConfigure   = &vcdaIPResource{}
	_ resource.ResourceWithImportState = &vcdaIPResource{}
)

// NewVcdaIPResource is a helper function to simplify the provider implementation.
func NewVCDAIPResource() resource.Resource {
	return &vcdaIPResource{}
}

// vcdaIPResource is the resource implementation.
type vcdaIPResource struct {
	client *client.CloudAvenue
	vcda   v1.VCDA
}

// Metadata returns the resource type name.
func (r *vcdaIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "ip"
}

// Schema defines the schema for the resource.
func (r *vcdaIPResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vcdaIPSchema().GetResource(ctx)
}

// Configure configures the resource.
func (r *vcdaIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.vcda = r.client.CAVSDK.V1.VCDA
}

// Create creates the resource and sets the initial Terraform state.
func (r *vcdaIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vcda_ip", r.client.GetOrgName(), metrics.Create)()

	plan := new(vcdaIPResourceModel)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	listOfIPS, err := r.vcda.List()
	if err != nil {
		resp.Diagnostics.AddError("Error on list VDCA IPs", err.Error())
		return
	}

	if listOfIPS.IsIPExists(plan.IPAddress.Get()) {
		resp.Diagnostics.AddError("IP address already registered", fmt.Sprintf("The IP address %s is already registered", plan.IPAddress.Get()))
		return
	}

	if err := r.vcda.RegisterIP(plan.IPAddress.Get()); err != nil {
		resp.Diagnostics.AddError("Error on register new VDCA IP", err.Error())
		return
	}

	stateRefreshed, found, diags := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vcdaIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vcda_ip", r.client.GetOrgName(), metrics.Read)()

	state := new(vcdaIPResourceModel)

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, diags := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vcdaIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vcdaIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vcda_ip", r.client.GetOrgName(), metrics.Delete)()

	state := new(vcdaIPResourceModel)

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	listOfIPs, err := r.vcda.List()
	if err != nil {
		resp.Diagnostics.AddError("Error on list VDCA IPs", err.Error())
		return
	}

	if listOfIPs.IsIPExists(state.IPAddress.Get()) {
		if err := listOfIPs.DeleteIP(state.IPAddress.Get()); err != nil {
			resp.Diagnostics.AddError("Error on delete VDCA IP", err.Error())
			return
		}
	}
}

func (r *vcdaIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip_address"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(urn.Normalize(
		urn.VCDA,
		utils.GenerateUUID(
			req.ID,
		).ValueString(),
	).String()))...)
}

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *vcdaIPResource) read(_ context.Context, planOrState *vcdaIPResourceModel) (stateRefreshed *vcdaIPResourceModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	listOfIps, err := r.vcda.List()
	if err != nil {
		diags.AddError("Error on list VDCA IPs", err.Error())
		return nil, true, diags
	}

	if !listOfIps.IsIPExists(planOrState.IPAddress.Get()) {
		return nil, false, diags
	}

	stateRefreshed.ID.Set(urn.Normalize(
		urn.VCDA,
		utils.GenerateUUID(planOrState.IPAddress.Get()).ValueString(),
	).String())

	return stateRefreshed, true, nil
}
