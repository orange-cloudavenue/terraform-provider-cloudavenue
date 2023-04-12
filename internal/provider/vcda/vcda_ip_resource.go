// Package vcda provides a Terraform resource to manage vcda.
package vcda

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
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
}

type vcdaIPResourceModel struct {
	ID        types.String `tfsdk:"id"`
	IPAddress types.String `tfsdk:"ip_address"`
}

// Metadata returns the resource type name.
func (r *vcdaIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "ip"
}

// Schema defines the schema for the resource.
func (r *vcdaIPResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The VCDa",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource allows you to declare or remove your on-premises IP address for the DRaaS service..\n" +
				" -> Note: For more information, please refer to the [Cloud Avenue DRaaS documentation](https://wiki.cloudavenue.orange-business.com/w/index.php/DRaaS_avec_VCDA).",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the VCDa resource.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
				},
			},
			"ip_address": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The on-premises IP address refers to the IP address of your local infrastructure running vCloud Extender.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						fstringvalidator.IsIP(),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
		},
	}.GetResource(ctx)
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
}

// Create creates the resource and sets the initial Terraform state.
func (r *vcdaIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan *vcdaIPResourceModel

	// Read Terraform plan data into the model.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// Call API to create the resource and check for errors.
	_, httpR, err := r.client.APIClient.VCDAApi.CreateVcdaIP(r.client.Auth, plan.IPAddress.ValueString())
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	// Set the ID
	plan.ID = plan.IPAddress

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vcdaIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *vcdaIPResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call API to get list of VCDA and check for errors.
	vcdaIPList, httpR, err := r.client.APIClient.VCDAApi.GetVcdaIPs(r.client.Auth)
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	found := slices.Contains(vcdaIPList, state.IPAddress.ValueString())
	// Check if the VCDA is in the list
	if !found {
		// If the VCDA is not in the list, remove it from the state
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = state.IPAddress

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vcdaIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vcdaIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	/// Get current state
	var state *vcdaIPResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// Call API to delete the resource and check for errors.
	_, httpR, err := r.client.APIClient.VCDAApi.DeleteVcdaIP(r.client.Auth, state.IPAddress.ValueString())
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}
}

func (r *vcdaIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("ip_address"), req, resp)
}
