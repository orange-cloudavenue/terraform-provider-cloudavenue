// Package client is the main client for the CloudAvenue provider.
package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	netbackupclient "github.com/orange-cloudavenue/netbackup-sdk-go"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

var (
	// ErrAuthFailed is returned when the authentication failed.
	ErrAuthFailed = errors.New("authentication error")
	// ErrTokenEmpty is returned when the token is empty.
	ErrTokenEmpty = errors.New("token is empty")
	// ErrConfigureVmware is returned when the configuration of vmware failed.
	ErrConfigureVmware = errors.New("error configuring vmware")
	ErrVCDVersionEmpty = errors.New("empty vcd version")
	// ErrConfigureNetBackup is returned when the configuration of netbackup failed.
	ErrConfigureNetBackup = errors.New("error configuring netbackup")
)

// CloudAvenue is the main struct for the CloudAvenue client.
type CloudAvenue struct {
	Org                string
	VDC                string
	User               string
	Password           string
	URL                string
	TerraformVersion   string
	CloudAvenueVersion string

	// API CLOUDAVENUE
	APIClient *apiclient.APIClient
	Auth      context.Context

	// API VMWARE
	Vmware     *govcd.VCDClient
	urlVmware  *url.URL
	VCDVersion string

	// API NetBackup
	NetBackupClient   *netbackupclient.Client
	NetBackupURL      string
	NetBackupUser     string
	NetBackupPassword string
}

// New creates a new CloudAvenue client.
func (c *CloudAvenue) New() (*CloudAvenue, error) {
	// API CLOUDAVENUE
	auth := c.createBasicAuthContext()
	cfg := c.createConfiguration()

	c.APIClient = apiclient.NewAPIClient(cfg)
	_, ret, err := c.APIClient.AuthenticationApi.GetToken(auth)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrAuthFailed, err)
	}
	token := ret.Header.Get("x-vmware-vcloud-access-token")
	if token == "" {
		return nil, ErrTokenEmpty
	}

	c.Auth = createTokenInContext(token)

	// API VMWARE
	err = c.configureVmware()
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrConfigureVmware, err)
	}

	if c.VCDVersion == "" {
		return nil, fmt.Errorf("%w : %w", ErrVCDVersionEmpty, err)
	}

	c.Vmware = govcd.NewVCDClient(*c.urlVmware, false, govcd.WithAPIVersion(c.VCDVersion))
	err = c.Vmware.SetToken(c.Org, govcd.AuthorizationHeader, token)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrConfigureVmware, err)
	}

	// API NetBackup
	if c.NetBackupURL != "" && c.NetBackupUser != "" && c.NetBackupPassword != "" {
		c.NetBackupClient, err = netbackupclient.New(netbackupclient.Opts{
			APIEndpoint: c.NetBackupURL,
			Username:    c.NetBackupUser,
			Password:    c.NetBackupPassword,
			Debug:       false,
		})
		if err != nil {
			return nil, fmt.Errorf("%w : %w", ErrConfigureNetBackup, err)
		}
	}

	return c, nil
}

// createBasicAuthContext creates a new context with the basic auth values.
func (c *CloudAvenue) createBasicAuthContext() context.Context {
	// Create a new CloudAvenue client using the configuration values
	auth := context.WithValue(context.Background(), apiclient.ContextBasicAuth, apiclient.BasicAuth{
		UserName: c.User + "@" + c.Org,
		Password: c.Password,
	})

	return auth
}

// createConfiguration creates a new configuration for the CloudAvenue client.
func (c *CloudAvenue) createConfiguration() *apiclient.Configuration {
	cfg := &apiclient.Configuration{
		BasePath:      c.URL,
		DefaultHeader: make(map[string]string),
		UserAgent:     c.createUserAgent(),
	}

	return cfg
}

// configurVmware creates a new configuration for the Vmware client.
func (c *CloudAvenue) configureVmware() (err error) {
	c.urlVmware, err = url.Parse(fmt.Sprintf("%s/api", c.GetURL()))
	return err
}

// createUserAgent creates a new user agent for the CloudAvenue client.
func (c *CloudAvenue) createUserAgent() string {
	return fmt.Sprintf("Terraform/%s CloudAvenue/%s", c.TerraformVersion, c.CloudAvenueVersion)
}

// createTokenInContext creates a new context with the token value.
func createTokenInContext(token string) context.Context {
	return context.WithValue(context.Background(), apiclient.ContextAccessToken, token)
}

// DefaultVDCExist returns true if the default VDC exists.
func (c *CloudAvenue) DefaultVDCExist() bool {
	return c.VDC != ""
}

// GetDefaultVDC returns the default VDC.
func (c *CloudAvenue) GetDefaultVDC() string {
	return c.VDC
}

// GetURL returns the base path of the API.
func (c *CloudAvenue) GetURL() string {
	return c.URL
}

// GetOrgName() returns the name of the organization.
// Deprecated: use GetOrg instead.
func (c *CloudAvenue) GetOrgName() string {
	return c.Org
}
