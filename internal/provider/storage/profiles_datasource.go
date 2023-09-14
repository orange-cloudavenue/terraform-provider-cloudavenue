// Package storage provides a Terraform datasource.
package storage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &profilesDataSource{}
	_ datasource.DataSourceWithConfigure = &profilesDataSource{}
)

func NewProfilesDataSource() datasource.DataSource {
	return &profilesDataSource{}
}

type profilesDataSource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
}

type profilesDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	VDC             types.String `tfsdk:"vdc"`
	StorageProfiles types.List   `tfsdk:"storage_profiles"`
}

func (s *storageProfiles) objectType(ctx context.Context) types.ObjectType {
	return types.ObjectType{AttrTypes: s.attrTypes(ctx)}
}

// Init Initializes the data source.
func (d *profilesDataSource) Init(ctx context.Context, dm *profilesDataSourceModel) (diags diag.Diagnostics) {
	d.vdc, diags = vdc.Init(d.client, dm.VDC)
	if diags.HasError() {
		return
	}

	return
}

func (d *profilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_profiles"
}

func (d *profilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.superSchema(ctx).GetDataSource(ctx)
}

func (d *profilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *profilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_storage_profiles", d.client.GetOrgName(), metrics.Read)()

	config := &profilesDataSourceModel{}

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

	sPs := make(storageProfiles, 0)
	for _, sP := range d.vdc.Vdc.Vdc.VdcStorageProfiles.VdcStorageProfile {
		x := profileDataSourceModel{
			ID:  types.StringValue(sP.ID),
			VDC: types.StringValue(d.vdc.GetName()),
		}

		storageProfile, err := d.client.Vmware.GetStorageProfileByHref(sP.HREF)
		if err != nil {
			resp.Diagnostics.AddError(
				"Storage Profile (Reference) not found",
				fmt.Sprintf("Storage profile %s not found in VDC %s", x.Name, d.vdc.GetName()),
			)
			return
		}

		x.Limit = types.Int64Value(storageProfile.Limit)
		x.VDC = types.StringValue(d.vdc.GetName())
		x.Name = types.StringValue(storageProfile.Name)
		x.UsedStorage = types.Int64Value(storageProfile.StorageUsedMB)
		x.Default = types.BoolValue(storageProfile.Default)
		x.Enabled = types.BoolValue(*storageProfile.Enabled)
		x.IopsAllocated = types.Int64Value(storageProfile.IopsAllocated)
		x.Units = types.StringValue(storageProfile.Units)
		x.IopsLimitingEnabled = types.BoolNull()
		x.MaximumDiskIops = types.Int64Null()
		x.DefaultDiskIops = types.Int64Null()
		x.DiskIopsPerGbMax = types.Int64Null()
		x.IopsLimit = types.Int64Null()
		if storageProfile.IopsSettings != nil {
			x.IopsLimitingEnabled = types.BoolValue(storageProfile.IopsSettings.Enabled)
			x.MaximumDiskIops = types.Int64Value(storageProfile.IopsSettings.DiskIopsMax)
			x.DefaultDiskIops = types.Int64Value(storageProfile.IopsSettings.DiskIopsDefault)
			x.DiskIopsPerGbMax = types.Int64Value(storageProfile.IopsSettings.DiskIopsPerGbMax)
			x.IopsLimit = types.Int64Value(storageProfile.IopsSettings.StorageProfileIopsLimit)
		}

		sPs = append(sPs, x)
	}

	var di diag.Diagnostics
	config.StorageProfiles, di = types.ListValueFrom(ctx, sPs.objectType(ctx), sPs)
	resp.Diagnostics.Append(di...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.VDC = types.StringValue(d.vdc.GetName())
	config.ID = utils.GenerateUUID(d.vdc.GetName())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
