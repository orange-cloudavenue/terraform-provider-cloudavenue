package org_test

import (
	"context"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
)

// TODO : Comment or uncomment the following imports if you are using resources or/and datasources

// Unit test for the schema of the resource cloudavenue_org_CertificateLibraryDatasourcesGo
// func TestCertificateLibraryDatasourcesGoResourceSchema(t *testing.T) {
// 	t.Parallel()

// 	ctx := context.Background()
// 	schemaResponse := &fwresource.SchemaResponse{}

// 	// Instantiate the resource.Resource and call its Schema method
// 	org.NewCertificateLibraryDatasource().Schema(ctx, fwresource.SchemaRequest{}, schemaResponse)

// 	if schemaResponse.Diagnostics.HasError() {
// 		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
// 	}

// 	// Validate the schema
// 	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

// 	if diagnostics.HasError() {
// 		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
// 	}
// }

// Unit test for the schema of the datasource cloudavenue_org_CertificateLibraryDatasourcesGo

func TestCertificateLibraryDatasourceSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaResponse := &fwdatasource.SchemaResponse{}

	// Instantiate the datasource.Datasource and call its Schema method

	// org.NewCertificateLibraryDatasource().Schema(ctx, fwdatasource.SchemaRequest{}, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}
