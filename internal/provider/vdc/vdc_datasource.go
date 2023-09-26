// Package vdc provides a resource to manage VDCs.
package vdc

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var (
	_ datasource.DataSource              = &vdcDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcDataSource{}
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
	defer metrics.New("data.cloudavenue_vdc", d.client.GetOrgName(), metrics.Read)()

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
			ID = uuid.Normalize(uuid.VDC, v.VdcUuid).String()
			break
		}
	}

	// Get storageProfile
	profiles := make(vdcResourceModelVDCStorageProfiles, 0)
	for _, profile := range vdc.Vdc.VdcStorageProfiles {
		p := vdcResourceModelVDCStorageProfile{}
		p.Class.Set(profile.Class)
		p.Limit.SetInt32(profile.Limit)
		p.Default.Set(profile.Default_)
		profiles = append(profiles, p)
	}

	data.ID.Set(ID)
	data.VDCGroup.Set(vdc.VdcGroup)
	data.Name.Set(vdc.Vdc.Name)
	data.Description.Set(vdc.Vdc.Description)
	data.VDCServiceClass.Set(vdc.Vdc.VdcServiceClass)
	data.VDCDisponibilityClass.Set(vdc.Vdc.VdcDisponibilityClass)
	data.VDCBillingModel.Set(vdc.Vdc.VdcBillingModel)
	data.VcpuInMhz2.Set(int64(vdc.Vdc.VcpuInMhz2))
	data.CPUAllocated.Set(int64(vdc.Vdc.CpuAllocated))
	data.MemoryAllocated.Set(int64(vdc.Vdc.MemoryAllocated))
	data.VDCStorageBillingModel.Set(vdc.Vdc.VdcStorageBillingModel)
	resp.Diagnostics.Append(data.VDCStorageProfiles.Set(ctx, profiles)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
