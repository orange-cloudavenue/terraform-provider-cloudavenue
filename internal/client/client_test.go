package client

import (
	"reflect"
	"testing"

	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

func TestCloudAvenueClient(t *testing.T) {
	t.Parallel()
	t.Run("CreateBasicAuthContext", func(t *testing.T) {
		t.Parallel()

		ca := CloudAvenue{
			User:     "dasilva",
			Password: "dasilva",
			Org:      "acme",
		}

		authCtx := ca.createBasicAuthContext()
		auth, isBasicAuth := authCtx.Value(apiclient.ContextBasicAuth).(apiclient.BasicAuth)
		if !isBasicAuth {
			t.Fatalf("expected context with cloudavenue.BasicAuth value, got %v", reflect.TypeOf(authCtx.Value(apiclient.ContextBasicAuth)))
		}

		if auth.UserName != "dasilva@acme" {
			t.Fatalf("expected username to be %q, got %q", "dasilva@acme", auth.UserName)
		}

		if auth.Password != "dasilva" {
			t.Fatalf("expected password to be %q, got %q", "dasilva", auth.Password)
		}
	})

	t.Run("CreateConfiguration", func(t *testing.T) {
		t.Parallel()

		ca := CloudAvenue{
			TerraformVersion:   "0.13.0",
			CloudAvenueVersion: "0.1.0",
			URL:                "https://console1.cloudavenue.orange-business.com",
		}

		cfg := ca.createConfiguration()
		emptyMap := make(map[string]string)

		if cfg.BasePath != "https://console1.cloudavenue.orange-business.com" {
			t.Fatalf("expected base path to be %q, got %q", "https://console1.cloudavenue.orange-business.com", cfg.BasePath)
		}

		if cfg.UserAgent != "Terraform/0.13.0 CloudAvenue/0.1.0" {
			t.Fatalf("expected user agent to be %q, got %q", "Terraform/0.13.0 CloudAvenue/0.1.0", cfg.UserAgent)
		}

		if !reflect.DeepEqual(cfg.DefaultHeader, emptyMap) {
			t.Fatalf("expected default header to be %v, got %v", emptyMap, cfg.DefaultHeader)
		}
	})

	t.Run("CreateTokenContext", func(t *testing.T) {
		t.Parallel()
		authCtx := createTokenInContext("t0k3n")
		token, isString := authCtx.Value(apiclient.ContextAccessToken).(string)

		if !isString {
			t.Fatalf("expected token with string value, got %v", reflect.TypeOf(authCtx.Value(apiclient.ContextAccessToken)))
		}

		if token != "t0k3n" {
			t.Fatalf("expected token to be %s, got %s", "t0k3n", token)
		}
	})

	t.Run("CreateUserAgent", func(t *testing.T) {
		t.Parallel()

		ca := CloudAvenue{
			TerraformVersion:   "0.13.0",
			CloudAvenueVersion: "0.1.0",
		}

		ua := ca.createUserAgent()

		if ua != "Terraform/0.13.0 CloudAvenue/0.1.0" {
			t.Fatalf("expected user agent to be %q, got %q", "Terraform/0.13.0 CloudAvenue/0.1.0", ua)
		}
	})

	t.Run("ConfigureVmware", func(t *testing.T) {
		t.Parallel()

		ca := CloudAvenue{
			TerraformVersion:   "0.13.0",
			CloudAvenueVersion: "0.1.0",
			URL:                "https://console1.cloudavenue.orange-business.com",
		}

		err := ca.configureVmware()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if ca.urlVmware.String() != "https://console1.cloudavenue.orange-business.com/api" {
			t.Fatalf("expected urlVmware to be %q, got %q", "https://console1.cloudavenue.orange-business.com/vmware", ca.urlVmware)
		}
	})

	t.Run("GetOrg", func(t *testing.T) {
		t.Parallel()

		ca := CloudAvenue{
			Org: "acme",
		}

		org := ca.GetOrg()

		if org != "acme" {
			t.Fatalf("expected org to be %q, got %q", "acme", org)
		}
	})

	t.Run("GetURL", func(t *testing.T) {
		t.Parallel()

		ca := CloudAvenue{
			URL: "https://console1.cloudavenue.orange-business.com",
		}

		url := ca.GetURL()

		if url != "https://console1.cloudavenue.orange-business.com" {
			t.Fatalf("expected url to be %q, got %q", "https://console1.cloudavenue.orange-business.com", url)
		}
	})

	t.Run("GetDefaultVdc", func(t *testing.T) {
		t.Parallel()

		ca := CloudAvenue{
			Vdc: "acme-vdc",
		}

		vdc := ca.GetDefaultVdc()

		if vdc != "acme-vdc" {
			t.Fatalf("expected default vdc to be %q, got %q", "acme-vdc", vdc)
		}
	})

	t.Run("DefaultVdcExist", func(t *testing.T) {
		t.Parallel()

		ca := CloudAvenue{
			Vdc: "acme-vdc",
		}

		exist := ca.DefaultVdcExist()

		if !exist {
			t.Fatalf("expected default vdc to exist")
		}
	})
}
