// Package publicip provides a Terraform resource to manage public IPs.
package publicip

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/antihax/optional"
	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &publicIPResource{}
	_ resource.ResourceWithConfigure   = &publicIPResource{}
	_ resource.ResourceWithImportState = &publicIPResource{}
)

// NewPublicIPResource returns a new resource implementing the public_ip resource.
func NewPublicIPResource() resource.Resource {
	return &publicIPResource{}
}

// publicIPResource is the resource implementation.
type publicIPResource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
}

// Metadata returns the resource type name.
func (r *publicIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Init.
func (r *publicIPResource) Init(_ context.Context, rm *publicIPResourceModel) (diags diag.Diagnostics) {
	r.adminOrg, diags = adminorg.Init(r.client)

	return
}

// Schema defines the schema for the resource.
func (r *publicIPResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = publicIPSchema().GetResource(ctx)
}

func (r *publicIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *publicIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan := &publicIPResourceModel{}

	var (
		edgeGateway            edgegw.EdgeGateway
		findIPNotAlreadyExists func(IPs apiclient.PublicIps) (interface{}, error)
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, errTO := plan.Timeouts.Create(ctx, 5*time.Minute)
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

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	edgeGateway, err := r.adminOrg.GetEdgeGateway(edgegw.BaseEdgeGW{
		Name: plan.EdgeGatewayName,
		ID:   plan.EdgeGatewayID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error getting Edge Gateway", err.Error())
		return
	}

	body := apiclient.PublicIPApiCreatePublicIPOpts{
		XNattedIP:    optional.EmptyString(),
		XVDCEdgeName: optional.NewString(edgeGateway.GetName()),
		XVDCName:     optional.EmptyString(),
	}

	// Create new Public IP
	// Set vars
	var (
		job     apiclient.Jobcreated
		httpR   *http.Response
		knowIPs []apiclient.PublicIpsNetworkConfig
	)

	// Store existing Public IP
	// Get Public IP
	publicIPs, httpR, err := r.client.APIClient.PublicIPApi.GetPublicIPs(auth)
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	knowIPs = append(knowIPs, publicIPs.NetworkConfig...)

	// Find an ip that is not already existing in the vdc
	// this function var is called later in the code...
	findIPNotAlreadyExists = func(IPs apiclient.PublicIps) (interface{}, error) {
		if len(IPs.NetworkConfig) == 0 {
			return nil, fmt.Errorf("no public ip found")
		}

		// knowIPs is a list of ips that are already existing in the vdc
		// we need to find an ip that is not in this list
		// compare the list of public ip after the creation of public ip to the list of public ip before the creation of new public ip
		for _, IP := range IPs.NetworkConfig {
			for j, knownIP := range knowIPs {
				if knownIP.UplinkIp == IP.UplinkIp {
					// if ip is equal then go to next ip to compare
					break
				}
				// if ip is not found on list of public ip before the creation then we found the new one
				if j == (len(knowIPs) - 1) {
					return IP, nil
				}
			}
		}
		return apiclient.PublicIpsNetworkConfig{}, fmt.Errorf("no public ip found")
	}

	job, httpR, err = r.client.APIClient.PublicIPApi.CreatePublicIP(auth, &body)
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
	errRetry := retry.RetryContext(ctxTO, createTimeout, func() *retry.RetryError {
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

	// get all Public IPs and find the new one
	checkPublicIPs, httpRc, errGet := r.client.APIClient.PublicIPApi.GetPublicIPs(auth)
	if httpRc != nil {
		defer func() {
			err = errors.Join(err, httpRc.Body.Close())
		}()
	}

	if apiErr := helpers.CheckAPIError(errGet, httpRc); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	pubIP, errFind := findIPNotAlreadyExists(checkPublicIPs)
	if errFind != nil {
		resp.Diagnostics.AddError("Error finding Public IP", errFind.Error())
		return
	}

	var publicIP apiclient.PublicIps

	publicIP.NetworkConfig = append(publicIP.NetworkConfig, pubIP.(apiclient.PublicIpsNetworkConfig))
	if len(publicIP.NetworkConfig) == 0 {
		resp.Diagnostics.AddError(
			"Error creating Public IP",
			"Could not create Public IP, unexpected error: no public IP found after creation",
		)
		return
	}

	plan.ID = types.StringValue(publicIP.NetworkConfig[0].UplinkIp)
	plan.EdgeGatewayID = types.StringValue(edgeGateway.GetID())
	plan.EdgeGatewayName = types.StringValue(edgeGateway.GetName())
	plan.PublicIP = types.StringValue(publicIP.NetworkConfig[0].UplinkIp)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *publicIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &publicIPResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read timeout
	readTimeout, errTO := state.Timeouts.Read(ctx, 5*time.Minute)
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

	// Get Public IP
	publicIPs, httpR, err := r.client.APIClient.PublicIPApi.GetPublicIPs(auth)

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
	for _, ip := range publicIPs.NetworkConfig {
		if state.ID.Equal(types.StringValue(ip.UplinkIp)) {
			state.EdgeGatewayName = types.StringValue(ip.EdgeGatewayName)
			state.PublicIP = types.StringValue(ip.UplinkIp)

			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctxTO)
		return
	}

	edgeGateway, err := r.adminOrg.GetEdgeGateway(edgegw.BaseEdgeGW{
		Name: state.EdgeGatewayName,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Edge Gateway",
			fmt.Sprintf("Could not get Edge Gateway, unexpected error: %s", err),
		)
		return
	}

	state.EdgeGatewayID = types.StringValue(edgeGateway.GetID())

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctxTO, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *publicIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *publicIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &publicIPResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
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

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// Delete the public IP
	job, httpR, err := r.client.APIClient.PublicIPApi.DeletePublicIP(auth, state.PublicIP.ValueString())
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
		return
	}
}

func (r *publicIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
