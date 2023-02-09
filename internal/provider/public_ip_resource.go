package provider

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
	sdkResource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
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
	client *CloudAvenueClient
}

type publicIPResourceModel struct {
	Timeouts      timeouts.Value             `tfsdk:"timeouts"`
	ID            types.String               `tfsdk:"id"`
	NattedIP      types.String               `tfsdk:"natted_ip"`
	EdgeName      types.String               `tfsdk:"edge_name"`
	EdgeID        types.String               `tfsdk:"edge_id"`
	VdcName       types.String               `tfsdk:"vdc_name"`
	InternalIP    types.String               `tfsdk:"internal_ip"`
	NetworkConfig publicIPNetworkConfigModel `tfsdk:"network_config"`
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
			"natted_ip": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Target IP to configure NAT. If not provided, it configure Double NAT with a new IP on INET VDC Edge. If the IP or IP range provided is in `100.64.102.1-100.64.102.253` Double NAT is configured. If the IP is a private IP Direct NAT is configured.",
			},
			"edge_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Edge Gateway.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("vdc_name")),
					stringvalidator.AlsoRequires(path.MatchRoot("edge_id")),
				},
			},
			"edge_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Edge Gateway.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("edge_name")),
					stringvalidator.AlsoRequires(path.MatchRoot("vdc_name")),
				},
			},
			"vdc_name": schema.StringAttribute{
				MarkdownDescription: "Public IP is natted toward the INET VDC Edge in the specified VDC Name. This parameter helps to find target VDC Edge in case of multiples INET VDC Edges with same names",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("edge_name")),
					stringvalidator.AlsoRequires(path.MatchRoot("edge_id")),
				},
			},
			"internal_ip": schema.StringAttribute{
				Description: "Internal IP address.",
				Computed:    true,
			},
			"network_config": schema.ListNestedAttribute{
				Description: "List of networks.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uplink_ip": schema.StringAttribute{
							Description: "Uplink IP address.",
							Computed:    true,
						},
						"translated_ip": schema.StringAttribute{
							Description: "Translated IP address.",
							Computed:    true,
						},
						"edge_gateway_name": schema.StringAttribute{
							Description: "The name of the edge gateway related to the public ip.",
							Computed:    true,
						},
					},
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
func (r *publicIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan                   *publicIPResourceModel
		findIPNotAlreadyExists func(Ips apiclient.PublicIps) (interface{}, error)
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

	auth, errCtx := getAuthContextWithTO(r.client.auth, ctxTO)
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
		edgeGateway, httR, err := r.client.EdgeGatewaysApi.ApiCustomersV20EdgesEdgeIdGet(auth, plan.EdgeID.ValueString())
		if apiErr := CheckAPIError(err, httR); apiErr != nil {
			resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
			if resp.Diagnostics.HasError() || apiErr.IsNotFound() {
				return
			}
		}

		plan.EdgeName = types.StringValue(edgeGateway.EdgeName)
	}

	// Create new edge gateway
	body := apiclient.PublicIPApiApiCustomersV10IpPostOpts{
		XNattedIP:    optional.NewString(plan.NattedIP.ValueString()),
		XVDCEdgeName: optional.NewString(plan.EdgeName.ValueString()),
		XVDCName:     optional.NewString(plan.VdcName.ValueString()),
	}

	// Set vars
	var (
		err     error
		job     apiclient.Jobcreated
		httpR   *http.Response
		knowIps []publicIPNetworkConfigModel
	)

	// Find an ip that is not already existing in the vdc
	findIPNotAlreadyExists = func(Ips apiclient.PublicIps) (interface{}, error) {
		if len(Ips.NetworkConfig) == 0 {
			return nil, fmt.Errorf("no public ip found")
		}

		// knowIps is a list of ips that are already existing in the vdc
		// we need to find an ip that is not in this list
		for _, ip := range Ips.NetworkConfig {
			notFound := false
			for _, knownIP := range knowIps {
				if knownIP.UPLinkIP.Equal(types.StringValue(ip.UplinkIp)) && knownIP.EdgeGatewayName.Equal(types.StringValue(ip.EdgeGatewayName)) {
					continue
				} else {
					notFound = true
					break
				}
			}
			if notFound {
				return ip, nil
			}
		}

		return publicIPNetworkConfigModel{}, fmt.Errorf("no public ip found")
	}

	job, httpR, err = r.client.PublicIPApi.ApiCustomersV10IpPost(auth, &body)
	if apiErr := CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Wait for job to complete
	refreshF := func() (interface{}, string, error) {
		var publicIP apiclient.PublicIps

		jobStatus, errGetJob := getJobStatus(auth, r.client, job.JobId)
		if errGetJob != nil {
			return nil, "", err
		}

		if jobStatus == "DONE" {
			// get all edge gateways and find the one that matches the tier0_vrf_id and owner_name
			publicIps, _, errEdgesGet := r.client.PublicIPApi.ApiCustomersV20IpGet(auth)
			if errEdgesGet != nil {
				return nil, "error", err
			}

			pubIP, err := findIPNotAlreadyExists(publicIps)
			if err != nil {
				return nil, "error", err
			}

			publicIP.InternalIp = publicIps.InternalIp
			publicIP.NetworkConfig = append(publicIP.NetworkConfig, pubIP.(apiclient.PublicIpsNetworkConfig))

			return publicIP, "done", nil
		}

		return nil, "pending", nil
	}

	createStateConf := &sdkResource.StateChangeConf{
		Delay:      10 * time.Second,
		Refresh:    refreshF,
		MinTimeout: 5 * time.Second,
		Timeout:    5 * time.Minute,
		Pending:    []string{"pending"},
		Target:     []string{"done"},
	}

	publicIP, err := createStateConf.WaitForStateContext(ctxTO)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create vdc, unexpected error: "+err.Error(),
		)
		return
	}

	// Set the ID
	ip, ok := publicIP.(apiclient.PublicIps)
	if !ok {
		resp.Diagnostics.AddError(
			"Error creating edge gateway",
			"Could not create edge gateway, unexpected error: publicIP is not a publicIPNetworkConfigModel",
		)
		return
	}

	plan.ID = types.StringValue(ip.NetworkConfig[0].UplinkIp)
	plan.NetworkConfig.EdgeGatewayName = types.StringValue(ip.NetworkConfig[0].EdgeGatewayName)
	plan.NetworkConfig.UPLinkIP = types.StringValue(ip.NetworkConfig[0].UplinkIp)
	plan.NetworkConfig.TranslatedIP = types.StringValue(ip.NetworkConfig[0].TranslatedIp)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
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

	readTimeout, errTO := state.Timeouts.Read(ctx, 1*time.Minute)
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
	publicIps, httpR, err := r.client.PublicIPApi.ApiCustomersV20IpGet(auth)
	if apiErr := CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}

		resp.State.RemoveResource(ctx)
		return
	}

	for _, ip := range publicIps.NetworkConfig {
		if state.NetworkConfig.UPLinkIP.Equal(types.StringValue(ip.UplinkIp)) && state.NetworkConfig.EdgeGatewayName.Equal(types.StringValue(ip.EdgeGatewayName)) {
			state.ID = types.StringValue(ip.UplinkIp)
			state.NetworkConfig.EdgeGatewayName = types.StringValue(ip.EdgeGatewayName)
			state.NetworkConfig.UPLinkIP = types.StringValue(ip.UplinkIp)
			state.NetworkConfig.TranslatedIP = types.StringValue(ip.TranslatedIp)
		}
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
	var state publicIPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the public IP
	_, httpR, err := r.client.PublicIPApi.ApiCustomersV10IpPublicIpDelete(ctx, state.NetworkConfig.UPLinkIP.ValueString())
	if apiErr := CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

func (r *publicIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
