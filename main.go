package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/cloudavenue"
)

// Example version string that can be overwritten by a release process
var version string = "dev"

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name cloudavenue

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/orange-cloudavenue/cloudavenue",
	}

	err := providerserver.Serve(context.Background(), cloudavenue.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
