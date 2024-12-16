package vdcg

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &NetworkIsolatedDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkIsolatedDataSource{}
)

func NewNetworkIsolatedDataSource() datasource.DataSource {
	return &NetworkIsolatedDataSource{}
}

type NetworkIsolatedDataSource struct {
	client *client.CloudAvenue
	vdcg   *v1.VDCGroup
}

func (d *NetworkIsolatedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_network_isolated"
}

// Init Initializes the resource.
func (d *NetworkIsolatedDataSource) Init(ctx context.Context, rm *networkIsolatedModel) (diags diag.Diagnostics) {
	var err error

	// Get the VDC Group
	if rm.VDCGroupID.IsKnown() {
		d.vdcg, err = d.client.CAVSDK.V1.VDC().GetVDCGroup(rm.VDCGroupID.Get())
	} else {
		d.vdcg, err = d.client.CAVSDK.V1.VDC().GetVDCGroup(rm.VDCGroupName.Get())
	}
	if err != nil {
		diags.AddError("Error retrieving VDC Group", err.Error())
		return
	}
	return
}

func (d *NetworkIsolatedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = networkIsolatedSchema(ctx).GetDataSource(ctx)
}

func (d *NetworkIsolatedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkIsolatedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdcg_network_isolated", d.client.GetOrgName(), metrics.Read)()

	config := &networkIsolatedModel{}

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

	s := &NetworkIsolatedResource{
		client: d.client,
		vdcg:   d.vdcg,
	}

	// Read data from the API
	data, found, diags := s.read(ctx, config)
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The isolated network '%s' was not found.", config.Name.Get()))
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
