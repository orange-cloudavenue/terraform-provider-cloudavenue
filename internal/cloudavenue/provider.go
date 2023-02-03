package cloudavenue

import (
	"context"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

// Ensure the implementation satisfies the expected interfaces
var _ provider.Provider = &cloudavenueProvider{}

// cloudavenueProvider is the provider implementation.
type cloudavenueProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type cloudavenueProviderModel struct {
	Url      types.String `tfsdk:"url"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Org      types.String `tfsdk:"org"`
}

type CloudAvenueClient struct {
	*apiclient.APIClient
	auth context.Context
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
func (p *cloudavenueProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the CloudAvenue API. Can also be set with the `CLOUDAVENUE_URL` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`[a-z]$`),
						"must end with a letter",
					),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The username to use to connect to the CloudAvenue API. Can also be set with the `CLOUDAVENUE_USER` environment variable.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to use to connect to the CloudAvenue API. Can also be set with the `CLOUDAVENUE_PASSWORD` environment variable.",
				Required:            true,
				Sensitive:           true,
			},
			"org": schema.StringAttribute{
				MarkdownDescription: "The organization used on CloudAvenue API. Can also be set with the `CLOUDAVENUE_ORG` environment variable.",
				Required:            true,
			},
		},
	}
}

func (p *cloudavenueProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring CloudAvenue client")
	var config cloudavenueProviderModel
	var client CloudAvenueClient

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	url := os.Getenv("CLOUDAVENUE_URL")
	user := os.Getenv("CLOUDAVENUE_USER")
	password := os.Getenv("CLOUDAVENUE_PASSWORD")
	org := os.Getenv("CLOUDAVENUE_ORG")

	if !config.Url.IsNull() {
		url = config.Url.ValueString()
	}
	if !config.User.IsNull() {
		user = config.User.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}
	if !config.Org.IsNull() {
		org = config.Org.ValueString()
	}

	// Default URL to the public CloudAvenue API if not set.
	if url == "" {
		url = "https://console1.cloudavenue.orange-business.com"
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if user == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("user"),
			"Missing CloudAvenue API User",
			"The provider cannot create the CloudAvenue API client as there is a missing or empty value for the CloudAvenue API user. "+
				"Set the host value in the configuration or use the CLOUDAVENUE_USER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing CloudAvenue API Password",
			"The provider cannot create the CloudAvenue API client as there is a missing or empty value for the CloudAvenue API password. "+
				"Set the host value in the configuration or use the CLOUDAVENUE_PASWWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if org == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("org"),
			"Missing CloudAvenue API Org",
			"The provider cannot create the CloudAvenue API client as there is a missing or empty value for the CloudAvenue API org. "+
				"Set the host value in the configuration or use the CLOUDAVENUE_ORG environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cloudavenue_host", url)
	ctx = tflog.SetField(ctx, "cloudavenue_username", user)
	ctx = tflog.SetField(ctx, "cloudavenue_org", org)
	ctx = tflog.SetField(ctx, "cloudavenue_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "cloudavenue_password")

	tflog.Debug(ctx, "Creating CloudAvenue client")

	// Create a new CloudAvenue client using the configuration values
	auth := context.WithValue(context.Background(), apiclient.ContextBasicAuth, apiclient.BasicAuth{
		UserName: user + "@" + org,
		Password: password,
	})

	cfg := &apiclient.Configuration{
		BasePath:      url,
		DefaultHeader: make(map[string]string),
		UserAgent:     "Terraform/" + req.TerraformVersion + "CloudAvenue/" + p.version,
	}

	client.APIClient = apiclient.NewAPIClient(cfg)
	_, ret, err := client.APIClient.AuthenticationApi.Cloudapi100SessionsPost(auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create CloudAvenue API Client",
			"An unexpected error occurred when creating the CloudAvenue API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"CloudAvenue Client Error: "+err.Error(),
		)
		return
	} else {
		token := ret.Header.Get("x-vmware-vcloud-access-token")
		if token == "" {
			resp.Diagnostics.AddError(
				"Unable to Create CloudAvenue API Client",
				"An unexpected error occurred when creating the CloudAvenue API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"CloudAvenue Client Error: "+err.Error(),
			)
			return
		}
		client.auth = context.WithValue(context.Background(), apiclient.ContextAccessToken, token)
	}

	// Make the CloudAvenue client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = &client
	resp.ResourceData = &client

	tflog.Info(ctx, "Configured CloudAvenue client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *cloudavenueProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTier0VrfsDataSource,
		NewTier0VrfDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *cloudavenueProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
