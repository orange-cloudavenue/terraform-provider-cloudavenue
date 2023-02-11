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
}
