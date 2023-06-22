// Package storage provides a Terraform datasource.
package storage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

var (
	_ datasource.DataSource              = &profileDataSource{}
	_ datasource.DataSourceWithConfigure = &profileDataSource{}
)

func NewProfileDataSource() datasource.DataSource {
	return &profileDataSource{}
}

type profileDataSource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
}

type profileDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	VDC                 types.String `tfsdk:"vdc"`
	Name                types.String `tfsdk:"name"`
	Limit               types.Int64  `tfsdk:"limit"`
	UsedStorage         types.Int64  `tfsdk:"used_storage"`
	Default             types.Bool   `tfsdk:"default"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	IopsAllocated       types.Int64  `tfsdk:"iops_allocated"`
	Units               types.String `tfsdk:"units"`
	IopsLimitingEnabled types.Bool   `tfsdk:"iops_limiting_enabled"`
	MaximumDiskIops     types.Int64  `tfsdk:"maximum_disk_iops"`
	DefaultDiskIops     types.Int64  `tfsdk:"default_disk_iops"`
	DiskIopsPerGbMax    types.Int64  `tfsdk:"disk_iops_per_gb_max"`
	IopsLimit           types.Int64  `tfsdk:"iops_limit"`
}

// Init Initializes the data source.
func (d *profileDataSource) Init(ctx context.Context, dm *profileDataSourceModel) (diags diag.Diagnostics) {
	d.vdc, diags = vdc.Init(d.client, dm.VDC)
	if diags.HasError() {
		return
	}

	return
}

func (d *profileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_profile"
}

func (d *profileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.superSchema(ctx).GetDataSource(ctx)
}

func (d *profileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *profileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	config := &profileDataSourceModel{}

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

	storageProfileID, err := d.vdc.FindStorageProfileID(config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Storage profile not found",
			fmt.Sprintf("Storage profile %s not found in VDC %s", config.Name.ValueString(), d.vdc.GetName()),
		)
		return
	}

	storageProfileRef, err := d.vdc.GetStorageProfileReference(storageProfileID, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Storage profile not found",
			fmt.Sprintf("Storage profile %s not found in VDC %s", config.Name.ValueString(), d.vdc.GetName()),
		)
		return
	}

	storageProfile, err := d.client.Vmware.GetStorageProfileByHref(storageProfileRef.HREF)
	if err != nil {
		resp.Diagnostics.AddError(
			"Storage profile not found",
			fmt.Sprintf("Storage profile %s not found in VDC %s", config.Name.ValueString(), d.vdc.GetName()),
		)
		return
	}

	config.ID = types.StringValue(storageProfileID)
	config.VDC = types.StringValue(d.vdc.GetName())
	config.Limit = types.Int64Value(storageProfile.Limit)
	config.UsedStorage = types.Int64Value(storageProfile.StorageUsedMB)
	config.Default = types.BoolValue(storageProfile.Default)
	config.Enabled = types.BoolValue(*storageProfile.Enabled)
	config.IopsAllocated = types.Int64Value(storageProfile.IopsAllocated)
	config.Units = types.StringValue(storageProfile.Units)
	config.IopsLimitingEnabled = types.BoolNull()
	config.MaximumDiskIops = types.Int64Null()
	config.DefaultDiskIops = types.Int64Null()
	config.DiskIopsPerGbMax = types.Int64Null()
	config.IopsLimit = types.Int64Null()
	if storageProfile.IopsSettings != nil {
		config.IopsLimitingEnabled = types.BoolValue(storageProfile.IopsSettings.Enabled)
		config.MaximumDiskIops = types.Int64Value(storageProfile.IopsSettings.DiskIopsMax)
		config.DefaultDiskIops = types.Int64Value(storageProfile.IopsSettings.DiskIopsDefault)
		config.DiskIopsPerGbMax = types.Int64Value(storageProfile.IopsSettings.DiskIopsPerGbMax)
		config.IopsLimit = types.Int64Value(storageProfile.IopsSettings.StorageProfileIopsLimit)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

func (d *profileDataSource) superSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_storage_profile` data source can be used to access information about a storage profile in a VDC.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "ID of storage profile.",
					Computed:            true,
				},
			},
			"name": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Name of storage profile.",
					Required:            true,
				},
			},
			"vdc": vdc.SuperSchema(),
			"limit": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Maximum number of storage bytes (scaled by 'units' field) allocated for this profile. `0` means `maximum possible`",
					Computed:            true,
				},
			},
			// used_storage
			"used_storage": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Storage used, in Megabytes, by the storage profile.",
					Computed:            true,
				},
			},
			// default
			"default": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Indicates whether this is the default storage profile for the VDC.",
					Computed:            true,
				},
			},
			// enabled
			"enabled": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "Indicates whether this storage profile is enabled for the VDC.",
					Computed:            true,
				},
			},
			// iops_allocated
			"iops_allocated": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Total IOPS currently allocated to this storage profile.",
					Computed:            true,
				},
			},
			// units
			"units": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Scale used to define Limit.",
					Computed:            true,
				},
			},
			// iops_limiting_enabled
			"iops_limiting_enabled": superschema.BoolAttribute{
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "True if this storage profile is IOPS-based placement enabled.",
					Computed:            true,
				},
			},
			// maximum_disk_iops
			"maximum_disk_iops": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "The maximum IOPS value that this storage profile is permitted to deliver. Value of 0 means this max setting is disabled and there is no max disk IOPS restriction.",
					Computed:            true,
				},
			},
			// default_disk_iops
			"default_disk_iops": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Value of 0 for disk IOPS means that no IOPS would be reserved or provisioned for that virtual disk.",
					Computed:            true,
				},
			},
			// disk_iops_per_gb_max
			"disk_iops_per_gb_max": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "The maximum IOPS per GB value that this storage profile is permitted to deliver. Value of 0 means this max setting is disabled and there is no max disk IOPS per GB restriction.",
					Computed:            true,
				},
			},
			// iops_limit
			"iops_limit": superschema.Int64Attribute{
				DataSource: &schemaD.Int64Attribute{
					MarkdownDescription: "Maximum number of IOPs that can be allocated for this profile. `0` means `maximum possible`.",
					Computed:            true,
				},
			},
		},
	}
}
