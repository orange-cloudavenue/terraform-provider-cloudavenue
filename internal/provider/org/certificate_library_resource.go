package org

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &CertificateLibraryResource{}
	_ resource.ResourceWithConfigure   = &CertificateLibraryResource{}
	_ resource.ResourceWithImportState = &CertificateLibraryResource{}
	// _ resource.ResourceWithModifyPlan     = &CertificateLibraryResource{}
	// _ resource.ResourceWithUpgradeState   = &CertificateLibraryResource{}
	// _ resource.ResourceWithValidateConfig = &CertificateLibraryResource{}.
)

// NewCertificateLibraryResource is a helper function to simplify the provider implementation.
func NewCertificateLibraryResource() resource.Resource {
	return &CertificateLibraryResource{}
}

// CertificateLibraryResource is the resource implementation.
type CertificateLibraryResource struct {
	client *client.CloudAvenue
	org    org.Org
}

// Init Initializes the resource.
func (r *CertificateLibraryResource) Init(ctx context.Context, rm *CertificateLibraryModel) (diags diag.Diagnostics) {
	// Uncomment the following lines if you need to access to the Org
	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	return
}

// Metadata returns the resource type name.
func (r *CertificateLibraryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_certificate_library"
}

// Schema defines the schema for the resource.
func (r *CertificateLibraryResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = certificateLibrarySchema(ctx).GetResource(ctx)
}

func (r *CertificateLibraryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *CertificateLibraryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_org_certificate_library", r.client.GetOrgName(), metrics.Create)()

	plan := &CertificateLibraryModel{}

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

	// Lock the resource organization
	mutex.GlobalMutex.KvLock(ctx, r.org.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.org.GetID())

	// Set the SDK certificate library Model
	certLibrarySDKModel := plan.ToSDKCertificateLibraryModel()

	// Create the certificate library
	newCertificate, err := r.org.CreateOrgCertificateLibrary(certLibrarySDKModel)
	if err != nil {
		resp.Diagnostics.AddError("error while creating certificate library: %s", err.Error())
		return
	}

	// Set ID
	plan.ID.Set(newCertificate.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The certificate library '%s' was not found after creation.", plan.Name.Get()))
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *CertificateLibraryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_org_certificate_library", r.client.GetOrgName(), metrics.Read)()

	state := &CertificateLibraryModel{}

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

	// Refresh the state
	stateRefreshed, found, d := r.read(state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *CertificateLibraryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_org_certificate_library", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &CertificateLibraryModel{}
		state = &CertificateLibraryModel{}
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

	// Lock the resource organization
	mutex.GlobalMutex.KvLock(ctx, r.org.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.org.GetID())

	// Set the SDK certificate library Model
	certLibrarySDKModel := plan.ToSDKCertificateLibraryModel()

	// Update the certificate library
	_, err := certLibrarySDKModel.Update()
	if err != nil {
		resp.Diagnostics.AddError("error while updating certificate library: %s", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *CertificateLibraryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_org_certificate_library", r.client.GetOrgName(), metrics.Delete)()

	state := &CertificateLibraryModel{}

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

	// Lock the resource organization
	mutex.GlobalMutex.KvLock(ctx, r.org.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.org.GetID())

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(state)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Delete the certificate library
	err := stateRefreshed.ToSDKCertificateLibraryModel().Delete()
	if err != nil {
		resp.Diagnostics.AddError("error while deleting certificate library: %s", err.Error())
		return
	}
}

func (r *CertificateLibraryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_org_certificate_library", r.client.GetOrgName(), metrics.Import)()

	// * ID format is CertificateIdOrName

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")
	if len(idParts) > 1 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: CertificateIdOrName. Got: %q", req.ID),
		)
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, &CertificateLibraryModel{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the certificate library
	certificate, err := r.org.GetOrgCertificateLibrary(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error while importing certificate library: %s", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), certificate.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), certificate.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), certificate.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("certificate"), certificate.Certificate)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("private_key"), certificate.PrivateKey)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("passphrase"), certificate.Passphrase)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *CertificateLibraryResource) read(planOrState *CertificateLibraryModel) (stateRefreshed *CertificateLibraryModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	var (
		certificate *v1.CertificateLibraryModel
		err         error
	)

	// Get CertificateLibrary
	if planOrState.ID.IsKnown() {
		certificate, err = r.org.GetOrgCertificateLibrary(planOrState.ID.Get())
	} else {
		certificate, err = r.org.GetOrgCertificateLibrary(planOrState.Name.Get())
	}
	if err != nil {
		if commoncloudavenue.IsNotFound(err) || govcd.IsNotFound(err) {
			return nil, false, diags
		}
		diags.AddError("error while fetching certificate library: %s", err.Error())
		return nil, false, diags
	}

	// Set the refreshed state
	stateRefreshed.ID.Set(certificate.ID)
	stateRefreshed.Name.Set(certificate.Name)
	stateRefreshed.Description.Set(certificate.Description)
	stateRefreshed.Certificate.Set(certificate.Certificate)
	// No need to set the private key and passphrase as they are not returned by the API (security reasons)

	return stateRefreshed, true, nil
}
