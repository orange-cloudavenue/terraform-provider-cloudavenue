// Package vapp provides a Terraform resource to manage vApps.
package vapp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

var (
	_ datasource.DataSource              = &vappDataSource{}
	_ datasource.DataSourceWithConfigure = &vappDataSource{}
)

// const vAppUnknownStatus = "-unknown-status-"

// NewVappDataSource returns a new vApp data source.
func NewVappDataSource() datasource.DataSource {
	return &vappDataSource{}
}

type vappDataSource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	vapp   vapp.VAPP
}

func (d *vappDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *vappDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vappSchema().GetDataSource(ctx)
}

func (d *vappDataSource) Init(ctx context.Context, dm *vappResourceModel) (diags diag.Diagnostics) {
	d.vdc, diags = vdc.Init(d.client, dm.VDC)
	if diags.HasError() {
		return
	}

	d.vapp, diags = vapp.Init(d.client, d.vdc, dm.VAppID, dm.VAppName)
	return
}

func (d *vappDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vappDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data  vappResourceModel
		diags diag.Diagnostics
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(d.Init(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set data
	data.Description = types.StringValue(d.vapp.GetDescription())
	data.VAppID = types.StringValue(d.vapp.GetID())
	data.VAppName = types.StringValue(d.vapp.GetName())
	data.VDC = types.StringValue(d.vdc.GetName())

	if d.vapp.GetStatusCode() == 4 {
		data.PowerON = types.BoolValue(true)
	} else {
		data.PowerON = types.BoolValue(false)
	}

	// Get guest properties
	guestProperties, diags := processGuestProperties(d.vapp)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.GuestProperties, diags = types.MapValue(types.StringType, guestProperties)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	leaseInfo, err := d.vapp.GetLease()
	if err != nil {
		resp.Diagnostics.AddError("Unable to get lease info", err.Error())
		return
	}

	if leaseInfo != nil {
		data.Lease, diags = types.ObjectValueFrom(ctx, vappLeaseAttrTypes, vappLeaseModel{
			RuntimeLeaseInSec: types.Int64Value(int64(leaseInfo.DeploymentLeaseInSeconds)),
			StorageLeaseInSec: types.Int64Value(int64(leaseInfo.StorageLeaseInSeconds)),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
