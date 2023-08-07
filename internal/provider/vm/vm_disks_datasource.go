// Package vm provides a Terraform datasource.
package vm

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var (
	_ datasource.DataSource              = &disksDataSource{}
	_ datasource.DataSourceWithConfigure = &disksDataSource{}
)

func NewDisksDataSource() datasource.DataSource {
	return &disksDataSource{}
}

type disksDataSource struct {
	client *client.CloudAvenue
	vapp   vapp.VAPP
	vdc    vdc.VDC
	org    org.Org
	vm     vm.VM
}

// Init Initializes the data source.
func (d *disksDataSource) Init(ctx context.Context, dm *DisksModel) (diags diag.Diagnostics) {
	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return
	}

	d.vdc, diags = vdc.Init(d.client, dm.VDC.StringValue)
	if diags.HasError() {
		return
	}

	d.vapp, diags = vapp.Init(d.client, d.vdc, dm.VAppID.StringValue, dm.VAppName.StringValue)
	if diags.HasError() {
		return
	}

	d.vm, diags = vm.Get(d.vapp, vm.GetVMOpts{
		ID:   dm.VMID.StringValue,
		Name: dm.VMName.StringValue,
	})
	return
}

func (d *disksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_disks"
}

func (d *disksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = disksSchema(ctx).GetDataSource(ctx)
}

func (d *disksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *disksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	config := &DisksModel{}

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

	disks := make(DisksModelDisks, 0)

	detachableDisks, err := d.vdc.QueryDisks("*")
	if err != nil {
		resp.Diagnostics.AddError("Unable to query disks", err.Error())
		return
	}

	for _, disk := range *detachableDisks {
		dDisk := DisksModelDisk{}
		id, err := govcd.GetUuidFromHref(disk.HREF, true)
		if err != nil {
			resp.Diagnostics.AddError("Unable to get disk ID", err.Error())
			return
		}
		dDisk.ID.Set(uuid.Normalize(uuid.Disk, id).String())
		dDisk.Name.Set(disk.Name)
		dDisk.SizeInMb.Set(disk.SizeMb)
		dDisk.StorageProfile.Set(disk.StorageProfileName)
		dDisk.IsDetachable.Set(true) // Because we are querying detachable disks

		disks = append(disks, dDisk)
	}

	if err := d.vm.Refresh(); err != nil {
		resp.Diagnostics.AddError("Unable to refresh VM", err.Error())
		return
	}

	if (d.vm.VM.VM.VM.VmSpecSection.DiskSection != nil && d.vm.VM.VM.VM.VmSpecSection.DiskSection.DiskSettings != nil) && len(d.vm.VM.VM.VM.VmSpecSection.DiskSection.DiskSettings) != 0 {
		for _, diskSetting := range d.vm.VM.VM.VM.VmSpecSection.DiskSection.DiskSettings {
			// if diskSetting.Disk != nil so it's an attached independent disk
			if diskSetting.Disk == nil {
				dDisk := DisksModelDisk{}
				dDisk.ID.Set(diskSetting.DiskId)
				dDisk.Name.SetNull() // Because we are querying attached disks
				dDisk.SizeInMb.Set(diskSetting.SizeMb)
				dDisk.StorageProfile.Set(diskSetting.StorageProfile.Name)
				dDisk.IsDetachable.Set(false) // Because we are querying attached disks

				disks = append(disks, dDisk)
			}
		}
	}

	idsDisks := make([]string, 0)
	for _, disk := range disks {
		idsDisks = append(idsDisks, disk.ID.Get())
	}

	if len(idsDisks) == 0 {
		idsDisks = append(idsDisks, utils.GenerateUUID("vm_disks").String())
	}

	config.ID.Set(utils.GenerateUUID(idsDisks).ValueString())
	config.VAppID.Set(d.vapp.GetID())
	config.VAppName.Set(d.vapp.GetName())
	config.VDC.Set(d.vdc.GetName())
	config.VMID.Set(d.vm.GetID())
	config.VMName.Set(d.vm.GetName())
	resp.Diagnostics.Append(config.Disks.Set(ctx, disks)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
