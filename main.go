// Package main is the main package for the CloudAvenue Terraform Provider.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider"

	_ "github.com/rs/zerolog"
	_ "gopkg.in/yaml.v3"
)

// Example version string that can be overwritten by a release process.
var version = "dev"

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name cloudavenue

func main() {
	var debug bool

	flag.BoolVar(
		&debug,
		"debug",
		false,
		"set to true to run the provider with support for debuggers like delve",
	)
	flag.Parse()

	// Generate a new execution ID for this run.
	// Not error checking here because it's not critical.
	x, _ := uuid.NewUUID()
	metrics.GlobalExecutionID = x.String()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/orange-cloudavenue/cloudavenue",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
