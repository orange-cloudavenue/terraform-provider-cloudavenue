package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider"
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cloudavenue": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// You can add code here to run prior to any test case execution, for example assertions
// about the appropriate environment variables being set are common to see in a pre-check
// function.
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CLOUDAVENUE_URL"); v == "" {
		t.Fatal("CLOUDAVENUE_URL must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_USER"); v == "" {
		t.Fatal("CLOUDAVENUE_USER must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_PASSWORD"); v == "" {
		t.Fatal("CLOUDAVENUE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_ORG"); v == "" {
		t.Fatal("CLOUDAVENUE_ORG must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUDAVENUE_VDC"); v == "" {
		t.Fatal("CLOUDAVENUE_VDC must be set for acceptance tests")
	}
}
