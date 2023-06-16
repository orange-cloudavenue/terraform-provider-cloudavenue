// Package provider provides the CloudAvenue Terraform Provider.
package provider

import (
	"context"
	"errors"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/alb"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/catalog"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/iam"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/publicip"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vrf"
)

const VCDVersion = "37.1"

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &cloudavenueProvider{}

// cloudavenueProvider is the provider implementation.
type cloudavenueProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type cloudavenueProviderModel struct {
	URL      types.String `tfsdk:"url"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Org      types.String `tfsdk:"org"`
	VDC      types.String `tfsdk:"vdc"`
}

// DataSources defines the data sources implemented in the provider.
func (p *cloudavenueProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// ALB
		alb.NewAlbPoolDataSource,

		// TIER0
		vrf.NewTier0VrfsDataSource,
		vrf.NewTier0VrfDataSource,

		// PUBLICIP
		publicip.NewPublicIPDataSource,

		// EDGE GATEWAY
		edgegw.NewEdgeGatewayDataSource,
		edgegw.NewEdgeGatewaysDataSource,

		// VDC
		vdc.NewVDCsDataSource,
		vdc.NewVDCDataSource,
		vdc.NewVDCGroupDataSource,

		// VAPP
		vapp.NewVappDataSource,
		vapp.NewOrgNetworkDataSource,
		vapp.NewIsolatedNetworkDataSource,
		// CATALOG
		catalog.NewCatalogsDataSource,
		catalog.NewCatalogDataSource,
		catalog.NewVAppTemplateDataSource,
		catalog.NewCatalogMediaDataSource,
		catalog.NewCatalogMediasDataSource,

		// IAM
		iam.NewUserDataSource,
		iam.NewRoleDataSource,
		iam.NewIAMRightDataSource,

		// VM
		vm.NewVMAffinityRuleDatasourceDataSource,

		// NETWORK
		network.NewNetworkIsolatedDataSource,
		network.NewNetworkRoutedDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *cloudavenueProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// ALB
		alb.NewAlbPoolResource,

		// EDGE GATEWAY
		edgegw.NewEdgeGatewayResource,

		// VDC
		vdc.NewVDCResource,
		vdc.NewACLResource,

		// VCDA
		// https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/pull/232
		// vcda.NewVCDAIPResource,

		// PUBLICIP
		// Bug in the API, we can't delete a public IP. Waiting for a fix.
		// publicip.NewPublicIPResource,

		// VAPP
		vapp.NewVappResource,
		vapp.NewOrgNetworkResource,
		vapp.NewIsolatedNetworkResource,
		vapp.NewACLResource,

		// CATALOG
		catalog.NewCatalogResource,

		// IAM
		iam.NewIAMUserResource,
		iam.NewRoleResource,

		// VM
		vm.NewDiskResource,
		vm.NewVMResource,
		vm.NewVMInsertedMediaResource,
		vm.NewVMAffinityRuleResource,

		// NETWORK
		network.NewNetworkRoutedResource,
		network.NewNetworkIsolatedResource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &cloudavenueProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *cloudavenueProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "cloudavenue"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *cloudavenueProvider) Schema(
	_ context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This provider offers utilities for working with the Cloud Avenue platform.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the Cloud Avenue API. Can also be set with the `CLOUDAVENUE_URL` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^https?:\/\/\S+\w$`),
						"must end with a letter",
					),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The username to use to connect to the Cloud Avenue API. Can also be set with the `CLOUDAVENUE_USER` environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to use to connect to the Cloud Avenue API. Can also be set with the `CLOUDAVENUE_PASSWORD` environment variable.",
				Sensitive:           true,
				Optional:            true,
			},
			"org": schema.StringAttribute{
				MarkdownDescription: "The organization used on Cloud Avenue API. Can also be set with the `CLOUDAVENUE_ORG` environment variable.",
				Optional:            true,
			},
			"vdc": schema.StringAttribute{
				MarkdownDescription: "The VDC used on Cloud Avenue API. Can also be set with the `CLOUDAVENUE_VDC` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *cloudavenueProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	tflog.Info(ctx, "Configuring Cloud Avenue client")
	var config cloudavenueProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	urlCloudAvenue := os.Getenv("CLOUDAVENUE_URL")
	user := os.Getenv("CLOUDAVENUE_USER")
	password := os.Getenv("CLOUDAVENUE_PASSWORD")
	org := os.Getenv("CLOUDAVENUE_ORG")
	vdc := os.Getenv("CLOUDAVENUE_VDC")

	if !config.URL.IsNull() && config.URL.ValueString() != "" {
		urlCloudAvenue = config.URL.ValueString()
	}
	if !config.User.IsNull() && config.User.ValueString() != "" {
		user = config.User.ValueString()
	}
	if !config.Password.IsNull() && config.Password.ValueString() != "" {
		password = config.Password.ValueString()
	}
	if !config.Org.IsNull() && config.Org.ValueString() != "" {
		org = config.Org.ValueString()
	}
	if !config.VDC.IsNull() && config.VDC.ValueString() != "" {
		vdc = config.VDC.ValueString()
	}
	if !config.VDC.IsNull() && config.VDC.ValueString() != "" {
		vdc = config.VDC.ValueString()
	}

	// Default URL to the public Cloud Avenue API if not set.
	if urlCloudAvenue == "" {
		urlCloudAvenue = "https://console1.cloudavenue.orange-business.com"
	}
	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if user == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("user"),
			"Missing Cloud Avenue API User",
			"The provider cannot create the Cloud Avenue API client as there is a missing or empty value for the Cloud Avenue API user. "+
				"Set the host value in the configuration or use the CLOUDAVENUE_USER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Cloud Avenue API Password",
			"The provider cannot create the Cloud Avenue API client as there is a missing or empty value for the Cloud Avenue API password. "+
				"Set the host value in the configuration or use the CLOUDAVENUE_PASWWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if org == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("org"),
			"Missing Cloud Avenue API Org",
			"The provider cannot create the Cloud Avenue API client as there is a missing or empty value for the Cloud Avenue API org. "+
				"Set the host value in the configuration or use the CLOUDAVENUE_ORG environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cloudavenue_host", urlCloudAvenue)
	ctx = tflog.SetField(ctx, "cloudavenue_username", user)
	ctx = tflog.SetField(ctx, "cloudavenue_org", org)
	ctx = tflog.SetField(ctx, "cloudavenue_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "cloudavenue_password")

	tflog.Debug(ctx, "Creating CloudAvenue client")

	cloudAvenue := client.CloudAvenue{
		URL:                urlCloudAvenue,
		User:               user,
		Password:           password,
		Org:                org,
		VDC:                vdc,
		TerraformVersion:   req.TerraformVersion,
		CloudAvenueVersion: p.version,
		VCDVersion:         VCDVersion,
	}

	cA, err := cloudAvenue.New()
	if err != nil {
		switch {
		case errors.Is(err, client.ErrAuthFailed):
			resp.Diagnostics.AddError(
				"Unable to Create Cloud Avenue API Client",
				"An unexpected error occurred when creating the Cloud Avenue API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"Cloud Avenue Client Error: "+err.Error(),
			)
			return
		case errors.Is(err, client.ErrTokenEmpty):
			resp.Diagnostics.AddError(
				"Unable to Create Cloud Avenue API Client",
				"An unexpected error occurred when creating the Cloud Avenue API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"Cloud Avenue Client Error: empty token",
			)
			return
		case errors.Is(err, client.ErrConfigureVmware):
			resp.Diagnostics.AddError(
				"Unable to Configure VMWare VCD Client",
				"An unexpected error occurred when creating the VMWare VCD Client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"VMWare VCD Client Error: "+err.Error(),
			)
			return
		case errors.Is(err, client.ErrVCDVersionEmpty):
			resp.Diagnostics.AddError(
				"Unable to Configure VMWare VCD Client",
				"An unexpected error occurred when creating the VMWare VCD Client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"VMWare VCD version is empty",
			)
			return
		default:
			resp.Diagnostics.AddError(
				"Unable to Create Cloud Avenue API Client",
				"An unexpected error occurred when creating the Cloud Avenue API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"unknown error: "+err.Error(),
			)
			return
		}
	}

	// Make the CloudAvenue client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = cA
	resp.ResourceData = cA

	tflog.Info(ctx, "Configured Cloud Avenue client", map[string]any{"success": true})
}
