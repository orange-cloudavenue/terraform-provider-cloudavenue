// Package publicip provides a Terraform resource to manage public IPs.
package publicip

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkResource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &publicIPResource{}
	_ resource.ResourceWithConfigure   = &publicIPResource{}
	_ resource.ResourceWithImportState = &publicIPResource{}
)

// NewEdgeGatewayResource is a helper function to simplify the provider implementation.
func NewPublicIPResource() resource.Resource {
	return &publicIPResource{}
}

// publicIPResource is the resource implementation.
type publicIPResource struct {
	client *client.CloudAvenue
}

type publicIPResourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	ID       types.String   `tfsdk:"id"`
	PublicIP types.String   `tfsdk:"public_ip"`
	EdgeName types.String   `tfsdk:"edge_name"`
	EdgeID   types.String   `tfsdk:"edge_id"`
	Vdc      types.String   `tfsdk:"vdc"`
}

// Metadata returns the resource type name.
func (r *publicIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip"
}

// Schema defines the schema for the resource.
func (r *publicIPResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The public IP resource allows you to manage a public IP on your Organization.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID is the public IP address.",
				Computed:            true,
			},
			"public_ip": schema.StringAttribute{
				MarkdownDescription: "Public IP address.",
				Computed:            true,
			},
			"edge_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Edge Gateway.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vdc"), path.MatchRoot("edge_id")),
				},
			},
			"edge_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Edge Gateway.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("edge_name"), path.MatchRoot("vdc")),
				},
			},
			"vdc": schema.StringAttribute{
				MarkdownDescription: "Public IP is natted toward the INET VDC Edge in the specified VDC Name. This parameter helps to find target VDC Edge in case of multiples INET VDC Edges with same names",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("edge_name"), path.MatchRoot("edge_id")),
				},
			},
		},
	}
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
	var (
		plan                   *publicIPResourceModel
		findIPNotAlreadyExists func(IPs apiclient.PublicIps) (interface{}, error)
		body                   apiclient.PublicIPApiCreatePublicIPOpts
	)

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
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

	// if edge_id is provided, get edge_name
	if !plan.EdgeID.IsNull() {
		// Get Edge Gateway Name
		edgeGateway, httR, err := r.client.APIClient.EdgeGatewaysApi.GetEdgeById(auth, plan.EdgeID.ValueString())
		if apiErr := helpers.CheckAPIError(err, httR); apiErr != nil {
			resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
			if resp.Diagnostics.HasError() || apiErr.IsNotFound() {
				return
			}
		}

		plan.EdgeName = types.StringValue(edgeGateway.EdgeName)
		body.XVDCEdgeName = optional.NewString(plan.EdgeName.ValueString())
		body.XVDCName = optional.EmptyString()
	}

	// if vdc is provided
	if !plan.Vdc.IsNull() {
		body.XVDCName = optional.NewString(plan.Vdc.ValueString())
		body.XVDCEdgeName = optional.EmptyString()
	}

	body.XNattedIP = optional.EmptyString()

	// Create new Public IP
	// Set vars
	var (
		err     error
		job     apiclient.Jobcreated
		httpR   *http.Response
		knowIPs []apiclient.PublicIpsNetworkConfig
	)

	// Store existing Public IP
	// Get Public IP
	publicIPs, httpR, err := r.client.APIClient.PublicIPApi.GetPublicIPs(auth)
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	knowIPs = append(knowIPs, publicIPs.NetworkConfig...)

	// Find an ip that is not already existing in the vdc
	findIPNotAlreadyExists = func(IPs apiclient.PublicIps) (interface{}, error) {
		if len(IPs.NetworkConfig) == 0 {
			return nil, fmt.Errorf("no public ip found")
		}

		// knowIPs is a list of ips that are already existing in the vdc
		// we need to find an ip that is not in this list
		for _, IP := range IPs.NetworkConfig {
			found := false
			for _, knownIP := range knowIPs {
				if knownIP.UplinkIp == IP.UplinkIp {
					continue
				} else {
					found = true
					break
				}
			}
			if found {
				return IP, nil
			}
		}

		return apiclient.PublicIpsNetworkConfig{}, fmt.Errorf("no public ip found")
	}

	job, httpR, err = r.client.APIClient.PublicIPApi.CreatePublicIP(auth, &body)
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Wait for job to complete
	refreshF := func() (interface{}, string, error) {
		var publicIP apiclient.PublicIps

		jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
		if errGetJob != nil {
			return nil, "", err
		}

		if jobStatus.IsDone() {
			// get all Public IPs and find the new one
			publicIPs, httpR, errGet := r.client.APIClient.PublicIPApi.GetPublicIPs(auth)
			if apiErr := helpers.CheckAPIError(errGet, httpR); apiErr != nil {
				resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
				if resp.Diagnostics.HasError() {
					return nil, "error", apiErr
				}
			}

			pubIP, errFind := findIPNotAlreadyExists(publicIPs)
			if errFind != nil {
				return nil, "error", err
			}

			publicIP.NetworkConfig = append(publicIP.NetworkConfig, pubIP.(apiclient.PublicIpsNetworkConfig))

			return publicIP, jobStatus.String(), nil
		}

		return nil, jobStatus.String(), nil
	}

	createStateConf := &sdkResource.StateChangeConf{
		Delay:      10 * time.Second,
		Refresh:    refreshF,
		MinTimeout: 5 * time.Second,
		Timeout:    5 * time.Minute,
		Pending:    []string{helpers.PENDING.String()},
		Target:     []string{helpers.DONE.String()},
	}

	publicIP, err := createStateConf.WaitForStateContext(ctxTO)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Public IP",
			"Could not create Public IP, unexpected error: "+err.Error(),
		)
		return
	}

	// Set the ID
	IP, ok := publicIP.(apiclient.PublicIps)
	if !ok {
		resp.Diagnostics.AddError(
			"Error creating Public IP",
			"Could not create Public IP, unexpected error: publicIP is not a apiclient.PublicIps",
		)
		return
	}

	if len(IP.NetworkConfig) == 0 {
		resp.Diagnostics.AddError(
			"Error creating Public IP",
			"Could not create Public IP, unexpected error: no public IP find after creation",
		)
		return
	}
	plan.ID = types.StringValue(IP.NetworkConfig[0].UplinkIp)
	plan.EdgeName = types.StringValue(IP.NetworkConfig[0].EdgeGatewayName)
	plan.PublicIP = types.StringValue(IP.NetworkConfig[0].UplinkIp)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *publicIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *publicIPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
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
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	found := false
	for _, ip := range publicIPs.NetworkConfig {
		if state.ID.Equal(types.StringValue(ip.UplinkIp)) {
			state.EdgeName = types.StringValue(ip.EdgeGatewayName)
			state.PublicIP = types.StringValue(ip.UplinkIp)

			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctxTO)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctxTO, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *publicIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *publicIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state *publicIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

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

	// Delete the public IP
	job, httpR, err := r.client.APIClient.PublicIPApi.DeletePublicIP(ctx, state.PublicIP.ValueString())
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Wait for job to complete
	deleteStateConf := &sdkResource.StateChangeConf{
		Delay: 10 * time.Second,
		Refresh: func() (interface{}, string, error) {
			jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
			if errGetJob != nil {
				return nil, "", errGetJob
			}

			return jobStatus, jobStatus.String(), nil
		},
		MinTimeout: 5 * time.Second,
		Timeout:    5 * time.Minute,
		Pending:    helpers.JobStatePending(),
		Target:     helpers.JobStateDone(),
	}

	_, err = deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Public IP",
			"Could not delete Public IP, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "Public IP deleted")
}

func (r *publicIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
