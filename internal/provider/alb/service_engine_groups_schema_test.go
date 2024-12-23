package alb_test

import (
	"context"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/alb"
)

// Unit test for the schema of the datasource cloudavenue_alb_ServiceEngineGroups.
func TestServiceEngineGroupsDataSourceSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaResponse := &fwdatasource.SchemaResponse{}

	// Instantiate the datasource.Datasource and call its Schema method
	alb.NewServiceEngineGroupsDataSource().Schema(ctx, fwdatasource.SchemaRequest{}, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}
