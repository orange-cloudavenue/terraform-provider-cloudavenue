// Package catalog provides a Terraform resource to manage catalogs.
package catalog

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &catalogResource{}
	_ resource.ResourceWithConfigure   = &catalogResource{}
	_ resource.ResourceWithImportState = &catalogResource{}
)

// NewCatalogResource is a helper function to simplify the provider implementation.
func NewCatalogResource() resource.Resource {
	return &catalogResource{}
}

// catalogResource is the resource implementation.
type catalogResource struct {
	client *client.CloudAvenue
}

type catalogResourceModel struct {
	ID               types.String `tfsdk:"id"`
	CatalogName      types.String `tfsdk:"catalog_name"`
	Description      types.String `tfsdk:"description"`
	StorageProfileID types.String `tfsdk:"storage_profile_id"`
	CreatedAt        types.String `tfsdk:"created_at"`
	OwnerName        types.String `tfsdk:"owner_name"`
	DeleteForce      types.Bool   `tfsdk:"delete_force"`
	DeleteRecursive  types.Bool   `tfsdk:"delete_recursive"`
	Href             types.String `tfsdk:"href"`
}

// Metadata returns the resource type name.
func (r *catalogResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *catalogResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Catalog resource allows you to manage a catalog in CloudAvenue.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID is a unique identifier for the catalog",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"catalog_name": schema.StringAttribute{
				// Not ForceNew, to allow the resource name to be updated
				Required:            true,
				MarkdownDescription: "A name for the Catalog",
			},
			"description": schema.StringAttribute{
				// Not ForceNew, to allow the resource description to be updated
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the catalog",
			},
			"storage_profile_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The ID of the storage profile",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The creation date of the catalog",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"delete_force": schema.BoolAttribute{
				// Not ForceNew, to allow the resource name to be updated
				Required:            true,
				MarkdownDescription: "When destroying use `delete_force=True` with `delete_recursive=True` to remove a catalog and any objects it contains, regardless of their state.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"delete_recursive": schema.BoolAttribute{
				// Not ForceNew, to allow the resource name to be updated
				Required:            true,
				MarkdownDescription: "When destroying use `delete_recursive=True` to remove the catalog and any objects it contains that are in a state that normally allows removal.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"href": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The HREF of the catalog",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the owner of the catalog",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
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
		plan *catalogResourceModel
		err  error

		storageProfiles *govcdtypes.CatalogStorageProfiles
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// catalog creation is accessible only for administrator account
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	if !plan.StorageProfileID.IsNull() && !plan.StorageProfileID.IsUnknown() {
		// Get storage profile
		storageProfiles, err = r.getStorageProfile(adminOrg, plan.StorageProfileID.ValueString(), false)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving Storage Profile", err.Error())
			return
		}
	}

	// Create catalog
	c, err := r.createCatalogStorageProfile(adminOrg, plan, storageProfiles)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Catalog", err.Error())
		return
	}

	plan.ID = types.StringValue(c.AdminCatalog.ID)
	plan.Href = types.StringValue(c.AdminCatalog.HREF)
	plan.OwnerName = types.StringValue(c.AdminCatalog.Owner.User.Name)
	plan.CreatedAt = types.StringValue(c.AdminCatalog.DateCreated)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *catalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *catalogResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	// if state.ID is empty use catalog name
	var catalogID string
	if state.ID.IsNull() {
		catalogID = state.CatalogName.ValueString()
	} else {
		catalogID = state.ID.ValueString()
	}

	adminCatalog, err := r.getCatalog(adminOrg, catalogID, false)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			log.Printf("[DEBUG] Unable to find catalog. Removing from tfstate")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error retrieving Catalog", err.Error())
		return
	}

	// var storageProfileID string
	// // Check if storage profile is set. Although storage profile structure accepts a list, in UI only one can be picked
	// if adminCatalog.AdminCatalog.CatalogStorageProfiles != nil && len(adminCatalog.AdminCatalog.CatalogStorageProfiles.VdcStorageProfile) > 0 {
	// 	// By default, API does not return Storage Profile Name in response. It has ID and HREF, but not Name so name
	// 	// must be looked up
	// 	storageProfileID = adminCatalog.AdminCatalog.CatalogStorageProfiles.VdcStorageProfile[0].ID
	// }

	plan := state

	plan.ID = types.StringValue(adminCatalog.AdminCatalog.ID)
	plan.CatalogName = types.StringValue(adminCatalog.AdminCatalog.Name)
	plan.Description = types.StringValue(adminCatalog.AdminCatalog.Description)
	plan.CreatedAt = types.StringValue(adminCatalog.AdminCatalog.DateCreated)
	plan.Href = types.StringValue(adminCatalog.AdminCatalog.HREF)
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

	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	adminCatalog, err := r.getCatalog(adminOrg, state.ID.ValueString(), false)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newAdminCatalog := govcd.NewAdminCatalogWithParent(&r.client.Vmware.Client, adminOrg)
	newAdminCatalog.AdminCatalog.ID = adminCatalog.AdminCatalog.ID
	newAdminCatalog.AdminCatalog.HREF = adminCatalog.AdminCatalog.HREF
	newAdminCatalog.AdminCatalog.Name = plan.CatalogName.ValueString()
	newAdminCatalog.AdminCatalog.Description = plan.Description.ValueString()

	// Check if StorageProfileID has changed
	if !plan.StorageProfileID.Equal(state.StorageProfileID) {
		if plan.StorageProfileID.IsNull() || plan.StorageProfileID.IsUnknown() || plan.StorageProfileID.ValueString() == "" {
			// If StorageProfileID is empty, remove storage profile from catalog
			newAdminCatalog.AdminCatalog.CatalogStorageProfiles = &govcdtypes.CatalogStorageProfiles{VdcStorageProfile: []*govcdtypes.Reference{}}
		} else {
			// If StorageProfileID is not empty, add storage profile to catalog
			storageProfileReference, errGet := r.getStorageProfileReference(adminOrg, plan.StorageProfileID.ValueString(), false)
			if errGet != nil {
				resp.Diagnostics.AddError("Error retrieving Storage Profile", errGet.Error())
				return
			}

			newAdminCatalog.AdminCatalog.CatalogStorageProfiles = &govcdtypes.CatalogStorageProfiles{VdcStorageProfile: []*govcdtypes.Reference{storageProfileReference}}
		}
	}

	if !plan.StorageProfileID.Equal(state.StorageProfileID) || !plan.Description.Equal(state.Description) || !plan.CatalogName.Equal(state.CatalogName) {
		// If field has changed, update it
		err = newAdminCatalog.Update()
		if err != nil {
			resp.Diagnostics.AddError("Error updating Catalog", err.Error())
			return
		}
	}

	// c, err := r.getCatalog(adminOrg, state.ID.ValueString(), true)
	// if err != nil {
	// 	resp.State.RemoveResource(ctx)
	// 	return
	// }

	// var storageProfileID string
	// // Check if storage profile is set. Although storage profile structure accepts a list, in UI only one can be picked
	// if c.AdminCatalog.CatalogStorageProfiles != nil && len(c.AdminCatalog.CatalogStorageProfiles.VdcStorageProfile) > 0 {
	// 	// By default, API does not return Storage Profile Name in response. It has ID and HREF, but not Name so name
	// 	// must be looked up
	// 	storageProfileID = c.AdminCatalog.CatalogStorageProfiles.VdcStorageProfile[0].ID
	// }

	// plan = &catalogResourceModel{
	// 	ID:              types.StringValue(c.AdminCatalog.ID),
	// 	CatalogName:     types.StringValue(c.AdminCatalog.Name),
	// 	Description:     types.StringValue(c.AdminCatalog.Description),
	// 	CreatedAt:       types.StringValue(c.AdminCatalog.DateCreated),
	// 	Href:            types.StringValue(c.AdminCatalog.HREF),
	// 	OwnerName:       types.StringValue(c.AdminCatalog.Owner.User.Name),
	// 	DeleteForce:     plan.DeleteForce,
	// 	DeleteRecursive: plan.DeleteRecursive,
	// }

	// if storageProfileID != "" {
	// 	plan.StorageProfileID = types.StringValue(storageProfileID)
	// }

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *catalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *catalogResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	// get catalog
	adminCatalog, err := r.getCatalog(adminOrg, state.ID.ValueString(), false)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// Catalog not found, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving Catalog", err.Error())
		return
	}

	err = adminCatalog.Delete(state.DeleteForce.ValueBool(), state.DeleteRecursive.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Catalog", err.Error())
		return
	}
}

func (r *catalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("catalog_name"), req, resp)
}

// getStorageProfile returns the storage profile reference.
func (r *catalogResource) getStorageProfile(adminOrg *govcd.AdminOrg, storageProfilID string, refresh bool) (*govcdtypes.CatalogStorageProfiles, error) {
	if storageProfilID == "" {
		return nil, errors.New("storageProfilID is an empty string")
	}

	// Get the storage profile
	storageProfileReference, err := adminOrg.GetStorageProfileReferenceById(storageProfilID, refresh)
	if err != nil {
		return nil, err
	}

	return &govcdtypes.CatalogStorageProfiles{VdcStorageProfile: []*govcdtypes.Reference{storageProfileReference}}, nil
}

// getStorageProfileReference returns the storage profile reference.
func (r *catalogResource) getStorageProfileReference(adminOrg *govcd.AdminOrg, storageProfilID string, refresh bool) (*govcdtypes.Reference, error) {
	return adminOrg.GetStorageProfileReferenceById(storageProfilID, refresh)
}

// createCatalogStorageProfile creates a storage profile reference.
func (r *catalogResource) createCatalogStorageProfile(adminOrg *govcd.AdminOrg, plan *catalogResourceModel, storageProfiles *govcdtypes.CatalogStorageProfiles) (*govcd.AdminCatalog, error) {
	return adminOrg.CreateCatalogWithStorageProfile(plan.CatalogName.ValueString(), plan.Description.ValueString(), storageProfiles)
}

// getCatalog returns the catalog reference.
func (r *catalogResource) getCatalog(adminOrg *govcd.AdminOrg, catalogNameOrID string, refresh bool) (*govcd.AdminCatalog, error) {
	// Get the catalog
	return adminOrg.GetAdminCatalogByNameOrId(catalogNameOrID, refresh)
}
