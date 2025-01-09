// Package org provides a Terraform datasource.
package org

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
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
func (d *CertificateLibraryDatasource) Init(ctx context.Context, dm *CertificateLibraryDatasourceModel) (diags diag.Diagnostics) {
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

	// Read data from the API
	data, _, diags := d.read(config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// read reads the data from the API and returns the data to be saved in the Terraform state.
func (d *CertificateLibraryDatasource) read(config *CertificateLibraryDatasourceModel) (data *CertificateLibraryDatasourceModel, found bool, diags diag.Diagnostics) {
	var (
		certificate *v1.CertificateLibraryModel
		err         error
	)

	// Get CertificateLibrary
	if config.ID.IsKnown() {
		certificate, err = d.org.GetOrgCertificateLibrary(config.ID.Get())
	} else {
		certificate, err = d.org.GetOrgCertificateLibrary(config.Name.Get())
	}
	if err != nil {
		if commoncloudavenue.IsNotFound(err) || govcd.IsNotFound(err) {
			return nil, false, diags
		}
		diags.AddError("error while fetching certificate library: %s", err.Error())
		return nil, true, diags
	}

	data = &CertificateLibraryDatasourceModel{
		ID:          supertypes.NewStringNull(),
		Name:        supertypes.NewStringNull(),
		Description: supertypes.NewStringNull(),
		Certificate: supertypes.NewStringNull(),
	}

	data.ID.Set(certificate.ID)
	data.Name.Set(certificate.Name)
	data.Description.Set(certificate.Description)
	data.Certificate.Set(certificate.Certificate)

	return data, true, nil
}
