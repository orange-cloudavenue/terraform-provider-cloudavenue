package s3_test

import (
	"context"
	"testing"

	// The fwresource import alias is so there is no collistion
	// with the more typical acceptance testing import:
	// "github.com/hashicorp/terraform-plugin-testing/helper/resource".
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/s3"
)

// TODO : Comment or uncomment the following imports if you are using resources or/and datasources

// Unit test for the schema of the resource cloudavenue_s3_3BucketAcl.
func Test3BucketAclResourceSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaResponse := &fwresource.SchemaResponse{}

	// Instantiate the resource.Resource and call its Schema method
	s3.NewBucketACLResource().Schema(ctx, fwresource.SchemaRequest{}, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

// Unit test for the schema of the datasource cloudavenue_s3_3BucketAcl
/*
func Test3BucketAclDataSourceSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaResponse := &fwresource.SchemaResponse{}

	// Instantiate the datasource.Datasource and call its Schema method
	s3.New3BucketAclDataSource().Schema(ctx, fwresource.SchemaRequest{}, schemaResponse)

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
