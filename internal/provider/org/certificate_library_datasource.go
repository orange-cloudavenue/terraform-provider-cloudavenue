// Package org provides a Terraform datasource.
package org

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/org"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &CertificateLibraryDatasource{}
	_ datasource.DataSourceWithConfigure = &CertificateLibraryDatasource{}
)

func NewCertificateLibraryDatasource() datasource.DataSource {
	return &CertificateLibraryDatasource{}
}

type CertificateLibraryDatasource struct {
	client    *client.CloudAvenue
	orgClient *org.Client
}

// Init Initializes the data source.
func (d *CertificateLibraryDatasource) Init(ctx context.Context, dm *CertificateLibraryDatasourceModel) (diags diag.Diagnostics) {
	var err error

	org, err := d.client.CAVSDK.V1.Org()
	if err != nil {
		diags.AddError("Error initializing ORG client", err.Error())
	}

	d.orgClient = org.Client

	return
}

func (d *CertificateLibraryDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_certificate_library"
}

func (d *CertificateLibraryDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = certificateLibrarySchema(ctx).GetDataSource(ctx)
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

	config := &CertificateLibraryDatasourceModel{}

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

	s := &CertificateLibraryResource{
		client:    d.client,
		orgClient: d.orgClient,
	}

	configResource := &CertificateLibraryModel{
		ID:   config.ID,
		Name: config.Name,
	}

	// Read data from the API
	data, found, diags := s.read(ctx, configResource)
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The Certificate %s(%s) was not found", config.Name.Get(), config.ID.Get()))
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.ID.Set(data.ID.Get())
	config.Name.Set(data.Name.Get())
	config.Description.Set(data.Description.Get())
	config.Certificate.Set(data.Certificate.Get())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
