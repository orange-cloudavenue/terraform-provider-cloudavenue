package iam_test

import (
	"context"
	"testing"

	// fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource".
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/iam"
)

// TODO : Comment or uncomment the following imports if you are using resources or/and datasources

// Unit test for the schema of the resource cloudavenue_iam_UserSaml.
func TestUserSamlResourceSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaResponse := &fwresource.SchemaResponse{}

	// Instantiate the resource.Resource and call its Schema method
	iam.NewUserSAMLResource().Schema(ctx, fwresource.SchemaRequest{}, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

// Unit test for the schema of the datasource cloudavenue_iam_UserSaml
/*
func TestUserSamlDataSourceSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaResponse := &fwdatasource.SchemaResponse{}

	// Instantiate the datasource.Datasource and call its Schema method
	iam.NewUserSamlDataSource().Schema(ctx, fwdatasource.SchemaRequest{}, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}
*/