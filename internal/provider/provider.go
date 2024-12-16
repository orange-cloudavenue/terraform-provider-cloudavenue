// Package provider provides the CloudAvenue Terraform Provider.
package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	casdk "github.com/orange-cloudavenue/cloudavenue-sdk-go"
	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	caverrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

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

func (p *cloudavenueProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config cloudavenueProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudAvenue := client.CloudAvenue{
		// This is a new SDK Cloudavenue
		CAVSDKOpts: &casdk.ClientOpts{
			Netbackup: &clientnetbackup.Opts{
				Endpoint: emptyOrValue(config.NetBackupURL),
				Username: emptyOrValue(config.NetBackupUser),
				Password: emptyOrValue(config.NetBackupPassword),
			},
			CloudAvenue: &clientcloudavenue.Opts{
				URL:      emptyOrValue(config.URL),
				Username: emptyOrValue(config.User),
				Password: emptyOrValue(config.Password),
				Org:      emptyOrValue(config.Org),
				VDC:      emptyOrValue(config.VDC),
			},
		},
	}

	const (
		summaryErrorAPICAV = "Unable to Create Cloud Avenue API Client"
		summaryErrorVCD    = "Unable to Create VMWare VCD Client"
	)

	cA, err := cloudAvenue.New()
	if err != nil {
		switch {
		case errors.Is(err, caverrors.ErrEmpty):
			resp.Diagnostics.AddError("Attribute is empty", err.Error())
			return
		case errors.Is(err, caverrors.ErrInvalidFormat):
			resp.Diagnostics.AddError("Attribute has invalid format", err.Error())
			return
		default:
			resp.Diagnostics.AddError(summaryErrorAPICAV, "unknown error: "+err.Error())
			return
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Make the CloudAvenue client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = cA
	resp.ResourceData = cA
}

func emptyOrValue(value basetypes.StringValue) string {
	if value.IsNull() {
		return ""
	}
	return value.ValueString()
}
