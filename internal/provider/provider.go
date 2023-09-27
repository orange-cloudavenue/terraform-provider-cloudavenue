// Package provider provides the CloudAvenue Terraform Provider.
package provider

import (
	"context"
	"errors"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
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

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &cloudavenueProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *cloudavenueProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudavenue"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *cloudavenueProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerSchema(ctx)
}

func (p *cloudavenueProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) { //nolint:gocyclo
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
	netbackupURL := os.Getenv("NETBACKUP_URL")
	netbackupUser := os.Getenv("NETBACKUP_USER")
	netbackupPassword := os.Getenv("NETBACKUP_PASSWORD")

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
	if !config.NetBackupURL.IsNull() {
		netbackupURL = config.NetBackupURL.ValueString()
	}
	if !config.NetBackupUser.IsNull() {
		netbackupUser = config.NetBackupUser.ValueString()
	}
	if !config.NetBackupPassword.IsNull() {
		netbackupPassword = config.NetBackupPassword.ValueString()
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
				"Set the host value in the configuration or use the CLOUDAVENUE_PASSWORD environment variable. "+
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
	// Default URL to the public NetBackup API if not set.
	if netbackupURL == "" {
		netbackupURL = "https://backup1.cloudavenue.orange-business.com/NetBackupSelfServiceNetBackupPanels/Api"
	}
	if netbackupUser == "" && netbackupPassword != "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("netbackup_user"),
			"Missing NetBackup API User",
			"The provider cannot create the NetBackup API client as there is a missing or empty value for the NetBackup API user. "+
				"Set the host value in the configuration or use the NETBACKUP_USER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if netbackupPassword == "" && netbackupUser != "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("netbackup_password"),
			"Missing NetBackup API Password",
			"The provider cannot create the NetBackup API client as there is a missing or empty value for the NetBackup API password. "+
				"Set the host value in the configuration or use the NETBACKUP_PASSWORD environment variable. "+
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
	ctx = tflog.SetField(ctx, "netbackup_host", netbackupURL)
	ctx = tflog.SetField(ctx, "netbackup_username", netbackupUser)
	ctx = tflog.SetField(ctx, "netbackup_password", netbackupPassword)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "netbackup_password")

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
		NetBackupURL:       netbackupURL,
		NetBackupUser:      netbackupUser,
		NetBackupPassword:  netbackupPassword,
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
		case errors.Is(err, client.ErrConfigureNetBackup):
			resp.Diagnostics.AddError(
				"Unable to Configure NetBackup Client",
				"An unexpected error occurred when creating the NetBackup Client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"NetBackup Client Error: "+err.Error(),
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
