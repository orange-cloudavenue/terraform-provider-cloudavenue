// Package vdc provides a resource to manage VDCs.
package vdc

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &vdcsDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcsDataSource{}
)

// NewVDCsDataSource returns a new resource implementing the vdcs data source.
func NewVDCsDataSource() datasource.DataSource {
	return &vdcsDataSource{}
}

type vdcsDataSource struct {
	client *client.CloudAvenue
}

func (d *vdcsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *vdcsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vdcsSchema()
}

func (d *vdcsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vdcsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdcs", d.client.GetOrgName(), metrics.Read)()

	var (
		data  vdcsDataSourceModel
		names []string
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vdcs, httpR, err := d.client.APIClient.VDCApi.GetOrgVdcs(d.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vdcs detail, got error: %s", err))
		return
	}
	defer func() {
		err = errors.Join(err, httpR.Body.Close())
	}()

	data = vdcsDataSourceModel{}

	for _, v := range vdcs {
		data.VDCs = append(data.VDCs, vdcRef{
			VDCName: types.StringValue(v.VdcName),
			VDCUuid: types.StringValue(v.VdcUuid),
		})
		names = append(names, v.VdcName)
	}

	data.ID = utils.GenerateUUID(names)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
