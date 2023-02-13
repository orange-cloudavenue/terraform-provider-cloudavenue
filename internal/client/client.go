package client

import (
	"context"
	"errors"
	"fmt"

	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

var (
	ErrAuthFailed = errors.New("authentication error")
	ErrTokenEmpty = errors.New("token is empty")
)

type CloudAvenue struct {
	User               string
	Org                string
	Password           string
	URL                string
	TerraformVersion   string
	CloudAvenueVersion string
	APIClient          *apiclient.APIClient
	Auth               context.Context
}

// New creates a new CloudAvenue client.
func (c *CloudAvenue) New() (*CloudAvenue, error) {
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
		UserAgent:     fmt.Sprintf("Terraform/%s CloudAvenue/%s", c.TerraformVersion, c.CloudAvenueVersion),
	}

	return cfg
}

// createTokenInContext creates a new context with the token value.
func createTokenInContext(token string) context.Context {
	return context.WithValue(context.Background(), apiclient.ContextAccessToken, token)
}
