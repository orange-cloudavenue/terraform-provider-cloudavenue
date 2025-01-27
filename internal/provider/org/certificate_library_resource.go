package org

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/org"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &CertificateLibraryResource{}
	_ resource.ResourceWithConfigure   = &CertificateLibraryResource{}
	_ resource.ResourceWithImportState = &CertificateLibraryResource{}
)

// NewCertificateLibraryResource is a helper function to simplify the provider implementation.
func NewCertificateLibraryResource() resource.Resource {
	return &CertificateLibraryResource{}
}

// CertificateLibraryResource is the resource implementation.
type CertificateLibraryResource struct {
	client    *client.CloudAvenue
	orgClient *org.Client
}

// Init Initializes the resource.
func (r *CertificateLibraryResource) Init(ctx context.Context, rm *CertificateLibraryModel) (diags diag.Diagnostics) {
	var err error

	org, err := r.client.CAVSDK.V1.Org()
	if err != nil {
		diags.AddError("Error initializing ORG client", err.Error())
	}

	r.orgClient = org.Client

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

	// Create the certificate library
	newCertificate, err := r.orgClient.CreateCertificateInLibrary(plan.ToSDKCertificateLibraryModel())
	if err != nil {
		resp.Diagnostics.AddError("error while creating certificate %s in library", err.Error())
		return
	}

	plan.ID.Set(newCertificate.Certificate.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The certificate '%s' was not found after creation.", plan.Name.Get()))
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
	stateRefreshed, found, d := r.read(ctx, state)
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

	certificate, err := r.orgClient.GetCertificateFromLibrary(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("error while fetching certificate : %s", err.Error())
		return
	}

	certificate.Certificate = plan.ToSDKCertificateLibraryModel()

	// Update the certificate library
	if err := certificate.Update(); err != nil {
		resp.Diagnostics.AddError("error while updating certificate %s", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
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

	certificate, err := r.orgClient.GetCertificateFromLibrary(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("error while fetching certificate : %s", err.Error())
		return
	}

	if err := certificate.Delete(); err != nil {
		resp.Diagnostics.AddError("error while deleting certificate : %s", err.Error())
		return
	}
}

func (r *CertificateLibraryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_org_certificate_library", r.client.GetOrgName(), metrics.Import)()

	// * ID format is CertificateIdOrName

	config := &CertificateLibraryModel{}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if urn.IsCertificateLibraryItem(req.ID) {
		config.ID.Set(req.ID)
	} else {
		config.Name.Set(req.ID)
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, config)
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The certificate '%s' was not found after import.", req.ID))
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *CertificateLibraryResource) read(_ context.Context, planOrState *CertificateLibraryModel) (stateRefreshed *CertificateLibraryModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	var (
		certificate *org.CertificateClient
		err         error
	)

	if planOrState.ID.IsKnown() {
		certificate, err = r.orgClient.GetCertificateFromLibrary(planOrState.ID.Get())
	} else {
		certificate, err = r.orgClient.GetCertificateFromLibrary(planOrState.Name.Get())
	}
	if err != nil {
		if govcd.IsNotFound(err) {
			return nil, false, diags
		}
		diags.AddError("error while fetching certificate library: %s", err.Error())
		return nil, false, diags
	}

	// Set the refreshed state
	stateRefreshed.ID.Set(certificate.Certificate.ID)
	stateRefreshed.Name.Set(certificate.Certificate.Name)
	stateRefreshed.Description.Set(certificate.Certificate.Description)
	stateRefreshed.Certificate.Set(certificate.Certificate.Certificate)

	return stateRefreshed, true, nil
}
