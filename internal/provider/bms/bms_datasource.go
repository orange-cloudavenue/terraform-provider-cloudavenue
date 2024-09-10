// Package bms provides a Terraform datasource.
package bms

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &bmsDataSource{}
	_ datasource.DataSourceWithConfigure = &bmsDataSource{}
)

func NewBMSDataSource() datasource.DataSource {
	return &bmsDataSource{}
}

type bmsDataSource struct {
	client *client.CloudAvenue

	// Uncomment the following lines if you need to access the resource's.
	// org    org.Org
	// vdc    vdc.VDC
	// vapp   vapp.VAPP
}

// If the data source don't have same schema/structure as the resource, you can use the following code:
// type DatasourceDataSourceModel struct {
// 	ID types.String `tfsdk:"id"`
// }

// Init Initializes the data source.
func (d *bmsDataSource) Init(ctx context.Context, dm *bmsModelDatasource) (diags diag.Diagnostics) {
	// Uncomment the following lines if you need to access to the Org
	// d.org, diags = org.Init(d.client)
	// if diags.HasError() {
	// 	return
	// }

	// Uncomment the following lines if you need to access to the VDC
	// d.vdc, diags = vdc.Init(d.client, dm.VDC)
	// if diags.HasError() {
	// 	return
	// }

	// Uncomment the following lines if you need to access to the VAPP
	// d.vapp, diags = vapp.Init(d.client, d.vdc, dm.VAppID, dm.VAppName)

	return
}

func (d *bmsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *bmsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = bmsSchema(ctx).GetDataSource(ctx)
}

func (d *bmsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *bmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_bms", d.client.GetOrgName(), metrics.Read)()

	config := &bmsModelDatasource{}

	// If the data source don't have same schema/structure as the resource, you can use the following code:
	// config := &DatasourceDataSourceModel{}

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
	// Set default timeouts
	readTimeout, diags := config.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	// Read data from the API - need to be implemented
	// ex: data, _, diags := d.bms.Get(ctx, config)

	// If read function is identical to the resource, you can use the following code:
	/*
		s := &DatasourceResource{
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
