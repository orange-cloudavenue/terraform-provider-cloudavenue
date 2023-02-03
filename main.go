package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/cloudavenue"
)

var (
	// Example version string that can be overwritten by a release process
	version string = "dev"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/orange-cloudavenue/cloudavenue",
	}

	err := providerserver.Serve(context.Background(), cloudavenue.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
