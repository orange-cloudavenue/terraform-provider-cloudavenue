// Package vdc provides a resource to manage VDCs.
package vdc

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
)

var (
	_ datasource.DataSource              = &vdcsDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcsDataSource{}
)

// NewVDCDataSource returns a new resource implementing the vdcs data source.
func NewVDCDataSource() datasource.DataSource {
	return &vdcDataSource{}
}

type vdcDataSource struct {
	client *client.CloudAvenue
}

func (d *vdcDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *vdcDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vdcSchema().GetDataSource(ctx)
}

func (d *vdcDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vdcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vdcDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get vDC info
	vdc, httpR, err := d.client.APIClient.VDCApi.GetOrgVdcByName(d.client.Auth, data.Name.ValueString())
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	// Get vDC UUID by parsing vDCs list and set URN ID
	var ID string
	vdcs, httpR, err := d.client.APIClient.VDCApi.GetOrgVdcs(d.client.Auth)

	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vdcs detail, got error: %s", err))
		return
	}

	for _, v := range vdcs {
		if data.Name.ValueString() == v.VdcName {
			ID = common.NormalizeID("urn:vcloud:vdc:", v.VdcUuid)
			break
		}
	}

	// Get storageProfile
	var profiles []vdcStorageProfileModel
	for _, profile := range vdc.Vdc.VdcStorageProfiles {
		p := vdcStorageProfileModel{
			Class:   types.StringValue(profile.Class),
			Limit:   types.Int64Value(int64(profile.Limit)),
			Default: types.BoolValue(profile.Default_),
		}
		profiles = append(profiles, p)
	}

	data = vdcDataSourceModel{
		ID:                     types.StringValue(ID),
		VDCGroup:               types.StringValue(vdc.VdcGroup),
		Name:                   types.StringValue(vdc.Vdc.Name),
		Description:            types.StringValue(vdc.Vdc.Description),
		VDCServiceClass:        types.StringValue(vdc.Vdc.VdcServiceClass),
		VDCDisponibilityClass:  types.StringValue(vdc.Vdc.VdcDisponibilityClass),
		VDCBillingModel:        types.StringValue(vdc.Vdc.VdcBillingModel),
		VcpuInMhz2:             types.Float64Value(vdc.Vdc.VcpuInMhz2),
		CPUAllocated:           types.Float64Value(vdc.Vdc.CpuAllocated),
		MemoryAllocated:        types.Float64Value(vdc.Vdc.MemoryAllocated),
		VDCStorageBillingModel: types.StringValue(vdc.Vdc.VdcStorageBillingModel),
		VDCStorageProfiles:     profiles,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
