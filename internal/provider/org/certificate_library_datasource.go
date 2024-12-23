// Package org provides a Terraform datasource.
package org

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

var (
	_ datasource.DataSource              = &CertificateLibraryDatasource{}
	_ datasource.DataSourceWithConfigure = &CertificateLibraryDatasource{}
)

func NewCertificateLibraryDatasource() datasource.DataSource {
	return &CertificateLibraryDatasource{}
}

type CertificateLibraryDatasource struct {
	client *client.CloudAvenue
	org    org.Org
}

// Init Initializes the data source.
func (d *CertificateLibraryDatasource) Init(ctx context.Context, dm *CertificateLibraryDatasourcesGoModel) (diags diag.Diagnostics) {
	// Uncomment the following lines if you need to access to the Org
	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return
	}

	return
}

func (d *CertificateLibraryDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_certificate_library"
}

func (d *CertificateLibraryDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = certificateLibraryDatasourceSchema(ctx).GetDataSource(ctx)
}

func (d *CertificateLibraryDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *CertificateLibraryDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_org_certificate_library", d.client.GetOrgName(), metrics.Read)()

	config := &CertificateLibraryDatasourcesGoModel{}

	// If the data source don't have same schema/structure as the resource, you can use the following code:
	// config := &CertificateLibraryDatasourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the data source read logic here.
	*/

	// If read function is identical to the resource, you can use the following code:
	/*
		s := &CertificateLibraryDatasourcesGoResource{
			client: d.client,
			// org:    d.org,
			// vdc:    d.vdc,
			// vapp:   d.vapp,
		}

		// Read data from the API
		data, _, diags := s.read(ctx, config)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	*/
}
