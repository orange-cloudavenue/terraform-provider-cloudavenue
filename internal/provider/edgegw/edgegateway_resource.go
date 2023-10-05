// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	apiclient "github.com/orange-cloudavenue/infrapi-sdk-go"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

const (
	defaultCheckJobDelayEdgeGateway = 10 * time.Second
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &edgeGatewaysResource{}
	_ resource.ResourceWithConfigure   = &edgeGatewaysResource{}
	_ resource.ResourceWithImportState = &edgeGatewaysResource{}

	// ConfigEdgeGateway is the default configuration for edge gateway.
	ConfigEdgeGateway setDefaultEdgeGateway = func() EdgeGatewayConfig {
		return EdgeGatewayConfig{
			CheckJobDelay: defaultCheckJobDelayEdgeGateway,
		}
	}
)

// NewEdgeGatewayResource returns a new resource implementing the edge_gateway data source.
func NewEdgeGatewayResource() resource.Resource {
	return &edgeGatewaysResource{}
}

type setDefaultEdgeGateway func() EdgeGatewayConfig

// EdgeGatewayConfig is the configuration for edge gateway.
type EdgeGatewayConfig struct {
	CheckJobDelay time.Duration
}

// edgeGatewaysResource is the resource implementation.
type edgeGatewaysResource struct {
	client *client.CloudAvenue
	EdgeGatewayConfig
}

// Metadata returns the resource type name.
func (r *edgeGatewaysResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *edgeGatewaysResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = edgegwSchema().GetResource(ctx)
}

func (r *edgeGatewaysResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = client
	r.EdgeGatewayConfig = ConfigEdgeGateway()
}

// Create creates the resource and sets the initial Terraform state.
func (r *edgeGatewaysResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan *edgeGatewaysResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Create)()

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

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
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

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	var err error
	var job apiclient.Jobcreated
	var httpR *http.Response

	switch plan.OwnerType.ValueString() {
	case "vdc":
		// Check if vDC exist
		if _, _, errGetOrg := r.client.GetOrgAndVDC(r.client.GetOrgName(), plan.OwnerName.ValueString()); errGetOrg != nil {
			resp.Diagnostics.AddError("Error retrieving VDC", errGetOrg.Error())
			return
		}

		// Create Edge Gateway
		job, httpR, err = r.client.APIClient.EdgeGatewaysApi.CreateVdcEdge(
			auth,
			body,
			plan.OwnerName.ValueString(),
		)

		if httpR != nil {
			defer func() {
				err = errors.Join(err, httpR.Body.Close())
			}()
		}

		if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
			resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
			return
		}
	case "vdc-group":
		// Check if vDC Group exist
		adminOrg, errGetAdminOrg := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
		if errGetAdminOrg != nil {
			resp.Diagnostics.AddError("Error retrieving Org", errGetAdminOrg.Error())
			return
		}
		if _, errGetVDCGroup := adminOrg.GetVdcGroupByName(plan.OwnerName.ValueString()); errGetVDCGroup != nil {
			resp.Diagnostics.AddError("Error retrieving vDC Group", errGetVDCGroup.Error())
			return
		}

		// Create Edge Gateway
		job, httpR, errGetAdminOrg = r.client.APIClient.EdgeGatewaysApi.CreateVdcGroupEdge(
			auth,
			body,
			plan.OwnerName.ValueString(),
		)

		if httpR != nil {
			defer func() {
				errGetAdminOrg = errors.Join(errGetAdminOrg, httpR.Body.Close())
			}()
		}

		if apiErr := helpers.CheckAPIError(errGetAdminOrg, httpR); apiErr != nil {
			resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
			return
		}
	}

	// Wait for job to complete
	errRetry := retry.RetryContext(ctxTO, createTimeout, func() *retry.RetryError {
		jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
		if errGetJob != nil {
			return retry.NonRetryableError(errGetJob)
		}
		if !slices.Contains(helpers.JobStateDone(), jobStatus.String()) {
			return retry.RetryableError(fmt.Errorf("expected job done but was %s", jobStatus))
		}

		return nil
	})

	if errRetry != nil {
		resp.Diagnostics.AddError("Error waiting job to complete", errRetry.Error())
		return
	}
	// Job done, retrieve edge gateway
	// get all edge gateways and find the one that matches the tier0_vrf_id and owner_name
	gateways, httpRc, errEdgesGet := r.client.APIClient.EdgeGatewaysApi.GetEdges(auth)
	if httpRc != nil {
		defer func() {
			err = errors.Join(err, httpRc.Body.Close())
		}()
	}
	if apiErr := helpers.CheckAPIError(errEdgesGet, httpRc); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	var newEdgeGW apiclient.EdgeGateway
	for _, gw := range gateways {
		if gw.Tier0VrfId == plan.Tier0VrfID.ValueString() && gw.OwnerName == plan.OwnerName.ValueString() {
			newEdgeGW = gw
			break
		}
	}

	if newEdgeGW == (apiclient.EdgeGateway{}) {
		resp.Diagnostics.AddError("Error retrieving edge gateway", "edge gateway not found after the create action")
		return
	}

	plan = &edgeGatewaysResourceModel{
		ID:                  types.StringValue(uuid.Normalize(uuid.Gateway, newEdgeGW.EdgeId).String()),
		Name:                types.StringValue(newEdgeGW.EdgeName),
		Description:         types.StringValue(newEdgeGW.Description),
		Tier0VrfID:          plan.Tier0VrfID,
		OwnerName:           plan.OwnerName,
		OwnerType:           plan.OwnerType,
		Timeouts:            plan.Timeouts,
		EnableLoadBalancing: plan.EnableLoadBalancing,
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *edgeGatewaysResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state *edgeGatewaysResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Read)()

	// Read timeout
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

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	var gateway apiclient.EdgeGateway
	// Get edge gateway
	if !state.ID.IsNull() {
		var (
			httpR *http.Response
			err   error
		)
		gateway, httpR, err = r.client.APIClient.EdgeGatewaysApi.GetEdgeById(
			auth,
			common.ExtractUUID(state.ID.ValueString()),
		)

		if httpR != nil {
			defer func() {
				err = errors.Join(err, httpR.Body.Close())
			}()
		}

		if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
			if !apiErr.IsNotFound() {
				resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
				return
			}

			resp.State.RemoveResource(ctxTO)
			return
		}
	} else {
		gateways, httpR, err := r.client.APIClient.EdgeGatewaysApi.GetEdges(auth)

		if httpR != nil {
			defer func() {
				err = errors.Join(err, httpR.Body.Close())
			}()
		}

		if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
			resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
			return
		}

		found := false
		for _, gateway = range gateways {
			if state.Name.Equal(types.StringValue(gateway.EdgeName)) {
				found = true
				break
			}
		}

		if !found {
			resp.State.RemoveResource(ctxTO)
			return
		}
	}

	// Get LoadBalancing state.
	gatewaysLoadBalancing, httpR, err := r.client.APIClient.EdgeGatewaysApi.GetEdgeLoadBalancing(auth, gateway.EdgeId)
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		defer httpR.Body.Close()
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	state = &edgeGatewaysResourceModel{
		ID:                  types.StringValue(uuid.Normalize(uuid.Gateway, gateway.EdgeId).String()),
		Tier0VrfID:          types.StringValue(gateway.Tier0VrfId),
		Name:                types.StringValue(gateway.EdgeName),
		OwnerType:           types.StringValue(gateway.OwnerType),
		OwnerName:           types.StringValue(gateway.OwnerName),
		Description:         types.StringValue(gateway.Description),
		EnableLoadBalancing: types.BoolValue(gatewaysLoadBalancing.Enabled),
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
func (r *edgeGatewaysResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan *edgeGatewaysResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Update)()

	// Update() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	updateTimeout, errTO := plan.Timeouts.Update(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctxTO, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Convert from Terraform data model into API data model
	body := apiclient.EdgeGatewayLoadBalancing{
		Enabled: plan.EnableLoadBalancing.ValueBool(),
	}

	var err error
	var job apiclient.Jobcreated
	var httpR *http.Response

	// Call API to update the resource and test for errors.
	job, httpR, err = r.client.APIClient.EdgeGatewaysApi.UpdateEdgeLoadBalancing(auth, body, common.ExtractUUID(plan.ID.ValueString()))
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		defer httpR.Body.Close()
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Wait for job to complete
	errRetry := retry.RetryContext(ctxTO, updateTimeout, func() *retry.RetryError {
		jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
		if errGetJob != nil {
			retry.NonRetryableError(err)
		}
		if !slices.Contains(helpers.JobStateDone(), jobStatus.String()) {
			return retry.RetryableError(fmt.Errorf("expected job done but was %s", jobStatus))
		}

		return nil
	})

	if errRetry != nil {
		resp.Diagnostics.AddError("Error waiting job to complete", errRetry.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *edgeGatewaysResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Get current state
	var state edgeGatewaysResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Delete)()

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// Delete timeout
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

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}
	// Delete the edge gateway
	job, httpR, err := r.client.APIClient.EdgeGatewaysApi.DeleteEdge(
		auth,
		common.ExtractUUID(state.ID.ValueString()),
	)

	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	// Wait for job to complete
	errRetry := retry.RetryContext(ctxTO, deleteTimeout, func() *retry.RetryError {
		jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
		if errGetJob != nil {
			retry.NonRetryableError(err)
		}
		if !slices.Contains(helpers.JobStateDone(), jobStatus.String()) {
			return retry.RetryableError(fmt.Errorf("expected job done but was %s", jobStatus))
		}

		return nil
	})

	if errRetry != nil {
		resp.Diagnostics.AddError("Error waiting job to complete", errRetry.Error())
	}
}

func (r *edgeGatewaysResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	defer metrics.New("cloudavenue_edgegateway", r.client.GetOrgName(), metrics.Import)()

	// Retrieve import Name and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
