package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkResource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &edgeGatewaysResource{}
	_ resource.ResourceWithConfigure   = &edgeGatewaysResource{}
	_ resource.ResourceWithImportState = &edgeGatewaysResource{}
)

// NewEdgeGatewayResource is a helper function to simplify the provider implementation.
func NewEdgeGatewayResource() resource.Resource {
	return &edgeGatewaysResource{}
}

// edgeGatewaysResource is the resource implementation.
type edgeGatewaysResource struct {
	client *CloudAvenueClient
}

type edgeGatewaysResourceModel struct {
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
	ID                  types.String   `tfsdk:"id"`
	Tier0VrfID          types.String   `tfsdk:"tier0_vrf_name"`
	EdgeName            types.String   `tfsdk:"edge_name"`
	EdgeID              types.String   `tfsdk:"edge_id"`
	OwnerType           types.String   `tfsdk:"owner_type"`
	OwnerName           types.String   `tfsdk:"owner_name"`
	Description         types.String   `tfsdk:"description"`
	EnableLoadBalancing types.Bool     `tfsdk:"enable_load_balancing"`
}

// Metadata returns the resource type name.
func (r *edgeGatewaysResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_edge_gateway"
}

// Schema defines the schema for the resource.
func (r *edgeGatewaysResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Edge Gateway resource allows you to create and delete Edge Gateways in CloudAvenue.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
			"id": schema.StringAttribute{
				Computed: true,
			},
			"tier0_vrf_name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The name of the Tier0 VRF to which the Edge Gateway will be attached.\n" +
					ForceNewDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"edge_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Edge Gateway.",
				Computed:            true,
			},
			"edge_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the Edge Gateway.",
			},
			"owner_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The type of the owner of the Edge Gateway (vdc|vdc-group).\n" +
					ForceNewDescription,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(vdc|vdc-group)$`),
						"must be vdc or vdc-group",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"owner_name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The name of the owner of the Edge Gateway.\n" +
					ForceNewDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the Edge Gateway.",
			},
			"enable_load_balancing": schema.BoolAttribute{
				MarkdownDescription: "Enable load balancing on the Edge Gateway.\n" +
					"Always set to true for now.",
				Computed: true,
			},
		},
	}
}

func (r *edgeGatewaysResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*CloudAvenueClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *CloudAvenueClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *edgeGatewaysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan *edgeGatewaysResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, errTO := plan.Timeouts.Create(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctxTO, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	auth, errCtx := getAuthContextWithTO(r.client.auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Create new edge gateway
	body := apiclient.EdgeGatewayCreate{
		Tier0VrfId:          plan.Tier0VrfID.ValueString(),
		EnableLoadBalancing: plan.EnableLoadBalancing.ValueBool(),
	}

	var err error
	var job apiclient.Jobcreated
	var httpR *http.Response

	switch plan.OwnerType.ValueString() {
	case "vdc":
		job, httpR, err = r.client.EdgeGatewaysApi.ApiCustomersV20VdcsVdcNameEdgesPost(auth, body, plan.OwnerName.ValueString())
		if apiErr := CheckAPIError(err, httpR); apiErr != nil {
			resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
			if resp.Diagnostics.HasError() {
				return
			}
		}
	case "vdc-group":
		job, httpR, err = r.client.EdgeGatewaysApi.ApiCustomersV20VdcGroupsVdcGroupNameEdgesPost(auth, body, plan.OwnerName.ValueString())
		if apiErr := CheckAPIError(err, httpR); apiErr != nil {
			resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	// Wait for job to complete
	refreshF := func() (interface{}, string, error) {
		jobStatus, errGetJob := getJobStatus(auth, r.client, job.JobId)
		if errGetJob != nil {
			return nil, "", err
		}

		edgeID := ""

		if jobStatus.IsDone() {
			// get all edge gateways and find the one that matches the tier0_vrf_id and owner_name
			gateways, _, errEdgesGet := r.client.EdgeGatewaysApi.ApiCustomersV20EdgesGet(auth)
			if errEdgesGet != nil {
				return nil, "err", err
			}

			for _, gw := range gateways {
				if gw.Tier0VrfId == plan.Tier0VrfID.ValueString() && gw.OwnerName == plan.OwnerName.ValueString() {
					edgeID = gw.EdgeId
					break
				}
			}
		} else {
			return nil, jobStatus.String(), nil
		}
		return edgeID, jobStatus.String(), nil
	}

	createStateConf := &sdkResource.StateChangeConf{
		Delay:      10 * time.Second,
		Refresh:    refreshF,
		MinTimeout: 5 * time.Second,
		Timeout:    5 * time.Minute,
		Pending:    StatePending(),
		Target:     []string{DONE.String()},
	}

	edgeID, err := createStateConf.WaitForStateContext(ctxTO)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating edge gateway",
			"Could not create edge gateway, unexpected error: "+err.Error(),
		)
		return
	}

	// Set the ID

	id, idIsAString := edgeID.(string)
	if !idIsAString {
		resp.Diagnostics.AddError(
			"Error creating edge gateway",
			"Could not create edge gateway, unexpected error: edgeID is not a string",
		)
		return
	}

	plan = &edgeGatewaysResourceModel{
		ID:                  types.StringValue(id),
		EdgeID:              types.StringValue(id),
		Tier0VrfID:          plan.Tier0VrfID,
		OwnerName:           plan.OwnerName,
		OwnerType:           plan.OwnerType,
		Timeouts:            plan.Timeouts,
		EnableLoadBalancing: types.BoolValue(true),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *edgeGatewaysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *edgeGatewaysResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create timeout
	readTimeout, errTO := state.Timeouts.Read(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctxTO, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	auth, errCtx := getAuthContextWithTO(r.client.auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Get edge gateway
	gateway, httpR, err := r.client.EdgeGatewaysApi.ApiCustomersV20EdgesEdgeIdGet(auth, state.EdgeID.ValueString())
	if apiErr := CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}

		resp.State.RemoveResource(ctxTO)
		return
	}

	state = &edgeGatewaysResourceModel{
		ID:                  types.StringValue(gateway.EdgeId),
		Tier0VrfID:          types.StringValue(gateway.Tier0VrfId),
		EdgeName:            types.StringValue(gateway.EdgeName),
		EdgeID:              types.StringValue(gateway.EdgeId),
		OwnerType:           types.StringValue(gateway.OwnerType),
		OwnerName:           types.StringValue(gateway.OwnerName),
		Description:         types.StringValue(gateway.Description),
		EnableLoadBalancing: types.BoolValue(true),
		Timeouts:            state.Timeouts,
	}

	// Set refreshed state
	diags = resp.State.Set(ctxTO, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *edgeGatewaysResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *edgeGatewaysResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state edgeGatewaysResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create timeout
	deleteTimeout, errTO := state.Timeouts.Delete(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctxTO, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	auth, errCtx := getAuthContextWithTO(r.client.auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Delete the edge gateway
	job, httpR, err := r.client.EdgeGatewaysApi.ApiCustomersV20EdgesEdgeIdDelete(auth, state.EdgeID.ValueString())
	if apiErr := CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Wait for job to complete
	deleteStateConf := &sdkResource.StateChangeConf{
		Delay: 10 * time.Second,
		Refresh: func() (interface{}, string, error) {
			jobStatus, errGetJob := getJobStatus(auth, r.client, job.JobId)
			if errGetJob != nil {
				return nil, "", err
			}

			return jobStatus, jobStatus.String(), nil
		},
		MinTimeout: 5 * time.Second,
		Timeout:    5 * time.Minute,
		Pending:    StatePending(),
		Target:     []string{DONE.String()},
	}

	_, err = deleteStateConf.WaitForStateContext(ctxTO)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Edge Gateway",
			"Could not delete Edge Gateway, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *edgeGatewaysResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("edge_id"), req, resp)
}
