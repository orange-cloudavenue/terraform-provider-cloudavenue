package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vcdaIPResource{}
	_ resource.ResourceWithConfigure   = &vcdaIPResource{}
	_ resource.ResourceWithImportState = &vcdaIPResource{}
)

// NewVcdaIPResource is a helper function to simplify the provider implementation.
func NewVcdaIPResource() resource.Resource {
	return &vcdaIPResource{}
}

// vcdaIPResource is the resource implementation.
type vcdaIPResource struct {
	client *CloudAvenueClient
}

type vcdaIPResourceModel struct {
	ID        types.String `tfsdk:"id"`
	IPAddress types.String `tfsdk:"ip_address"`
}

// Metadata returns the resource type name.
func (r *vcdaIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vcda_ip"
}

// Schema defines the schema for the resource.
func (r *vcdaIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "VCDa resource permit to declare or remove your On Premise IP address for DRaaS Service." +
			" -> Note: For more information, please refer to the [Cloud Avenue DRaaS documentation](https://wiki.cloudavenue.orange-business.com/w/index.php/DRaaS_avec_VCDA).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"ip_address": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "On Premise IP address. This is the IP address of our on premise infrastructure which run vCloud Extender.\n" +
					ForceNewDescription,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`),
						"must be a valide ipv4 address",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *vcdaIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*CloudAvenueClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *CloudAvenueClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	// Call API to create the resource and check for errors.
	_, httpR, err := r.client.VCDAApi.ApiCustomersV20VcdaIpsIpAddressPost(r.client.auth, plan.IPAddress.ValueString())
	if apiErr := CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set the ID
	plan.ID = plan.IPAddress

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

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
	vcdaIPList, httpR, err := r.client.VCDAApi.ApiCustomersV20VcdaIpsGet(r.client.auth)
	if apiErr := CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Check if the VCDA is in the list
	found := false
	for _, vcdaIP := range vcdaIPList {
		if vcdaIP == state.IPAddress.ValueString() {
			found = true
			state.ID = types.StringValue(vcdaIP)
			break
		}
	}

	// If the VCDA is not in the list, remove it from the state
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

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

	// Call API to delete the resource and check for errors.
	_, httpR, err := r.client.VCDAApi.ApiCustomersV20VcdaIpsIpAddressDelete(r.client.auth, state.IPAddress.ValueString())
	if apiErr := CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

func (r *vcdaIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("ip_address"), req, resp)
}
