// Package client is the main client for the CloudAvenue provider.
package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

var (
	ErrAuthFailed      = errors.New("authentication error")
	ErrTokenEmpty      = errors.New("token is empty")
	ErrConfigureVmware = errors.New("error configuring vmware")
)

type CloudAvenue struct {
	Org                string
	Vdc                string
	User               string
	Password           string
	URL                string
	TerraformVersion   string
	CloudAvenueVersion string

	// API CLOUDAVENUE
	APIClient *apiclient.APIClient
	Auth      context.Context

	// API VMWARE
	Vmware    *govcd.VCDClient
	urlVmware *url.URL
}

// New creates a new CloudAvenue client.
func (c *CloudAvenue) New() (*CloudAvenue, error) {
	// API CLOUDAVENUE
	auth := c.createBasicAuthContext()
	cfg := c.createConfiguration()

	c.APIClient = apiclient.NewAPIClient(cfg)
	_, ret, err := c.APIClient.AuthenticationApi.GetToken(auth)
	if err != nil {
		return nil, fmt.Errorf("%w : %v", ErrAuthFailed, err)
	}
	token := ret.Header.Get("x-vmware-vcloud-access-token")
	if token == "" {
		return nil, ErrTokenEmpty
	}

	c.Auth = createTokenInContext(token)

	// API VMWARE
	err = c.configureVmware()
	if err != nil {
		return nil, fmt.Errorf("%w : %v", ErrConfigureVmware, err)
	}

	c.Vmware = govcd.NewVCDClient(*c.urlVmware, false)
	err = c.Vmware.SetToken(c.GetOrg(), govcd.AuthorizationHeader, token)
	if err != nil {
		return nil, fmt.Errorf("%w : %v", ErrConfigureVmware, err)
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

// GetOrg returns the default Org
func (c *CloudAvenue) GetOrg() string {
	return c.Org
}

// DefaultVdcExists returns true if the default VDC exists
func (c *CloudAvenue) DefaultVdcExist() bool {
	return c.Vdc != ""
}

// GetDefaultVdc returns the default VDC
func (c *CloudAvenue) GetDefaultVdc() string {
	return c.Vdc
}

// GetBasePath returns the base path of the API
func (c *CloudAvenue) GetURL() string {
	return c.URL
}
