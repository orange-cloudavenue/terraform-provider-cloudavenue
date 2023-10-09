// Package client is the main client for the CloudAvenue provider.
package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientca "github.com/orange-cloudavenue/cloudavenue-sdk-go"
	apiclient "github.com/orange-cloudavenue/infrapi-sdk-go"
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
	// NetBackup *NetBackup
	NetBackupClient *clientca.Client
	NetBackupOpts   *clientca.ClientOpts
}

// New creates a new CloudAvenue client.
func (c *CloudAvenue) New() (*CloudAvenue, error) {
	// API CLOUDAVENUE
	c.APIClient = apiclient.NewAPIClient(c.createConfiguration())
	_, ret, err := c.APIClient.AuthenticationApi.GetToken(c.createBasicAuthContext())
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrAuthFailed, err)
	}
	token := ret.Header.Get("x-vmware-vcloud-access-token")
	if token == "" {
		return nil, ErrTokenEmpty
	}

	c.Auth = createTokenInContext(token)

	// API VMWARE
	if err = c.configureVmware(); err != nil {
		return nil, fmt.Errorf("%w : %w", ErrConfigureVmware, err)
	}

	if c.VCDVersion == "" {
		return nil, fmt.Errorf("%w : %w", ErrVCDVersionEmpty, err)
	}

	c.Vmware = govcd.NewVCDClient(*c.urlVmware, false, govcd.WithAPIVersion(c.VCDVersion))
	if err = c.Vmware.SetToken(c.Org, govcd.AuthorizationHeader, token); err != nil {
		return nil, fmt.Errorf("%w : %w", ErrConfigureVmware, err)
	}

	// API NetBackup
	// if c.NetBackup.Netbackup.Endpoint != "" || c.NetBackup.Netbackup.Username != "" || c.NetBackup.Netbackup.Password != "" {
	c.NetBackupClient, err = clientca.New(*c.NetBackupOpts)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrConfigureNetBackup, err)
	}
	// if err := c.NewNetBackupClient(); err != nil {
	// 	return nil, err
	// }
	// }

	return c, nil
}

// createBasicAuthContext creates a new context with the basic auth values.
func (c *CloudAvenue) createBasicAuthContext() context.Context {
	// Create a new CloudAvenue client using the configuration values
	return context.WithValue(context.Background(), apiclient.ContextBasicAuth, apiclient.BasicAuth{
		UserName: c.User + "@" + c.Org,
		Password: c.Password,
	})
}

// createConfiguration creates a new configuration for the CloudAvenue client.
func (c *CloudAvenue) createConfiguration() *apiclient.Configuration {
	return &apiclient.Configuration{
		BasePath:      c.URL,
		DefaultHeader: make(map[string]string),
		UserAgent:     c.createUserAgent(),
	}
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
func (c *CloudAvenue) GetOrgName() string {
	return c.Org
}
