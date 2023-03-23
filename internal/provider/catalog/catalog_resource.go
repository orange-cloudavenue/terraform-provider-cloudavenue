// Package catalog provides a Terraform resource to manage catalogs.
package catalog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &catalogResource{}
	_ resource.ResourceWithConfigure   = &catalogResource{}
	_ resource.ResourceWithImportState = &catalogResource{}
	_ catalog                          = &catalogResource{}
)

// NewCatalogResource is a helper function to simplify the provider implementation.
func NewCatalogResource() resource.Resource {
	return &catalogResource{}
}

// catalogResource is the resource implementation.
type catalogResource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	catalog  base
}

func (r *catalogResource) Init(_ context.Context, rm *catalogResourceModel) (diags diag.Diagnostics) {
	r.catalog = base{
		name: rm.Name.ValueString(),
		id:   "",
	}

	r.adminOrg, diags = adminorg.Init(r.client)

	return
}

// Metadata returns the resource type name.
func (r *catalogResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *catalogResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = catalogSchema().GetResource()
}

func (r *catalogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *catalogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan = &catalogResourceModel{}
		err  error

		storageProfiles *govcdtypes.CatalogStorageProfiles
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.StorageProfile.IsNull() && !plan.StorageProfile.IsUnknown() {
		// Get storage profile
		storageProfiles, err = r.adminOrg.GetStorageProfile(plan.StorageProfile.ValueString(), false)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving Storage Profile", err.Error())
			return
		}
	}

	// Create catalog
	c, err := r.createCatalogStorageProfile(plan, storageProfiles)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Catalog", err.Error())
		return
	}

	plan.ID = types.StringValue(c.AdminCatalog.ID)
	plan.OwnerName = types.StringValue(c.AdminCatalog.Owner.User.Name)
	plan.CreatedAt = types.StringValue(c.AdminCatalog.DateCreated)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *catalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &catalogResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminCatalog, err := r.GetCatalog()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error retrieving Catalog", err.Error())
		return
	}

	plan := state
	plan.ID = types.StringValue(adminCatalog.AdminCatalog.ID)
	plan.Name = types.StringValue(adminCatalog.AdminCatalog.Name)
	plan.Description = types.StringValue(adminCatalog.AdminCatalog.Description)
	plan.CreatedAt = types.StringValue(adminCatalog.AdminCatalog.DateCreated)
	plan.OwnerName = types.StringValue(adminCatalog.AdminCatalog.Owner.User.Name)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *catalogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *catalogResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminCatalog, err := r.GetCatalog()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Catalog", err.Error())
		return
	}

	newAdminCatalog := govcd.NewAdminCatalogWithParent(&r.client.Vmware.Client, r.adminOrg)
	newAdminCatalog.AdminCatalog.ID = adminCatalog.AdminCatalog.ID
	newAdminCatalog.AdminCatalog.HREF = adminCatalog.AdminCatalog.HREF
	newAdminCatalog.AdminCatalog.Name = plan.Name.ValueString()
	newAdminCatalog.AdminCatalog.Description = plan.Description.ValueString()

	// Check if StorageProfileID has changed
	if !plan.StorageProfile.Equal(state.StorageProfile) {
		if plan.StorageProfile.IsNull() || plan.StorageProfile.IsUnknown() || plan.StorageProfile.ValueString() == "" {
			// If StorageProfileID is empty, remove storage profile from catalog
			newAdminCatalog.AdminCatalog.CatalogStorageProfiles = &govcdtypes.CatalogStorageProfiles{VdcStorageProfile: []*govcdtypes.Reference{}}
		} else {
			// If StorageProfile is not empty, add storage profile to catalog
			storageProfileID, err := r.adminOrg.FindStorageProfileID(plan.StorageProfile.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Error retrieving Storage Profile", err.Error())
				return
			}

			storageProfileReference, errGet := r.adminOrg.GetStorageProfileReference(storageProfileID, false)
			if errGet != nil {
				resp.Diagnostics.AddError("Error retrieving Storage Profile", errGet.Error())
				return
			}
			newAdminCatalog.AdminCatalog.CatalogStorageProfiles = &govcdtypes.CatalogStorageProfiles{VdcStorageProfile: []*govcdtypes.Reference{storageProfileReference}}
		}
	}

	if !plan.StorageProfile.Equal(state.StorageProfile) || !plan.Description.Equal(state.Description) || !plan.Name.Equal(state.Name) {
		// If field has changed, update it
		err = newAdminCatalog.Update()
		if err != nil {
			resp.Diagnostics.AddError("Error updating Catalog", err.Error())
			return
		}
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *catalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &catalogResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get catalog
	adminCatalog, err := r.GetCatalog()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// Catalog not found, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving Catalog", err.Error())
		return
	}

	if err = adminCatalog.Delete(state.DeleteForce.ValueBool(), state.DeleteRecursive.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Error deleting Catalog", err.Error())
		return
	}
}

func (r *catalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// createCatalogStorageProfile creates a storage profile reference.
func (r *catalogResource) createCatalogStorageProfile(plan *catalogResourceModel, storageProfiles *govcdtypes.CatalogStorageProfiles) (*govcd.AdminCatalog, error) {
	return r.adminOrg.CreateCatalogWithStorageProfile(plan.Name.ValueString(), plan.Description.ValueString(), storageProfiles)
}

func (r *catalogResource) GetID() string {
	return r.catalog.id
}

// GetName returns the name of the catalog.
func (r *catalogResource) GetName() string {
	return r.catalog.name
}

// GetIDOrName returns the ID if it is set, otherwise it returns the name.
func (r *catalogResource) GetIDOrName() string {
	if r.GetID() != "" {
		return r.GetID()
	}
	return r.GetName()
}

// GetCatalog returns the govcd.Catalog.
func (r *catalogResource) GetCatalog() (*govcd.AdminCatalog, error) {
	return r.adminOrg.GetAdminCatalogByNameOrId(r.GetIDOrName(), true)
}
