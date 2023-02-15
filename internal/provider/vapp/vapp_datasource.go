// Package vapp provides a Terraform resource to manage vApps.
package vapp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &vappDataSource{}
	_ datasource.DataSourceWithConfigure = &vappDataSource{}
)

const vAppUnknownStatus = "-unknown-status-"

func NewVappDataSource() datasource.DataSource {
	return &vappDataSource{}
}

type vappDataSource struct {
	client *client.CloudAvenue
}

type vappDataSourceModel struct {
	ID              types.String                  `tfsdk:"id"`
	VappName        types.String                  `tfsdk:"vapp_name"`
	VappID          types.String                  `tfsdk:"vapp_id"`
	VDC             types.String                  `tfsdk:"vdc"`
	Description     types.String                  `tfsdk:"description"`
	Href            types.String                  `tfsdk:"href"`
	GuestProperties map[types.String]types.String `tfsdk:"guest_properties"`
	StatusCode      types.Int64                   `tfsdk:"status_code"`
	StatusText      types.String                  `tfsdk:"status_text"`
	Lease           types.Object                  `tfsdk:"lease"`
}

func (d *vappDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vapp"
}

func (d *vappDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The public IP data source show the list of public IP addresses.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				CustomType:          types.StringType,
				Computed:            true,
				MarkdownDescription: "The ID is a `vapp_id`.",
			},
			"vapp_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A name for the vApp, unique within the VDC. Required if `vapp_id` is not set.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_id")),
				},
			},
			"vapp_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The ID of vApp to use, unique within the VDC. Required if `vapp_name` is not set.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name")),
				},
			},
			"vdc": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name of VDC to use, optional if defined at provider level",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional description of the vApp",
			},
			"href": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "vApp Hyper Reference",
			},
			"guest_properties": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Key/value settings for guest properties",
			},
			"status_code": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Shows the status code of the vApp",
			},
			"status_text": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Shows the status of the vApp",
			},
		},
		Blocks: map[string]schema.Block{
			"lease": schema.SingleNestedBlock{
				// MarkdownDescription: "Defines lease parameters for this vApp",
				Attributes: map[string]schema.Attribute{
					"runtime_lease_in_sec": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "How long any of the VMs in the vApp can run before the vApp is automatically powered off or suspended. 0 means never expires",
					},
					"storage_lease_in_sec": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "How long the vApp is available before being automatically deleted or marked as expired. 0 means never expires",
					},
				},
			},
		},
	}
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
	var data vappDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If VDC is not defined at data source level, use the one defined at provider level
	if data.VDC.IsNull() {
		if d.client.DefaultVdcExist() {
			data.VDC = types.StringValue(d.client.GetDefaultVdc())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	_, vdc, err := d.client.GetOrgAndVdc(d.client.GetOrg(), data.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find VDC", err.Error())
		return
	}

	var v string

	if !data.VappID.IsNull() {
		v = data.VappID.ValueString()
	} else {
		v = data.VappName.ValueString()
	}

	// Request vApp
	vapp, err := vdc.GetVAppByNameOrId(v, false)
	if err != nil {
		resp.Diagnostics.AddError("Unable to find vApp", err.Error())
		return
	}

	// update guest properties
	guestProperties, err := vapp.GetProductSectionList()
	if err != nil {
		resp.Diagnostics.AddError("Unable to get guest properties", err.Error())
		return
	}

	if guestProperties != nil && guestProperties.ProductSection != nil && guestProperties.ProductSection.Property != nil {
		for _, guestProperty := range guestProperties.ProductSection.Property {
			if guestProperty.Value != nil {
				data.GuestProperties[types.StringValue(guestProperty.Key)] = types.StringValue(guestProperty.Value.Value)
			}
		}
	}

	leaseInfo, err := vapp.GetLease()
	if err != nil {
		resp.Diagnostics.AddError("Unable to get lease info", err.Error())
		return
	}

	if leaseInfo != nil {
		data.Lease.As(ctx, &vappLeaseModel{
			RuntimeLeaseInSec: types.Int64Value(int64(leaseInfo.DeploymentLeaseInSeconds)),
			StorageLeaseInSec: types.Int64Value(int64(leaseInfo.StorageLeaseInSeconds)),
		}, basetypes.ObjectAsOptions{})
	}

	statusText, err := vapp.GetStatus()
	if err != nil {
		statusText = vAppUnknownStatus
	}

	data.StatusCode = types.Int64Value(int64(vapp.VApp.Status))
	data.StatusText = types.StringValue(statusText)
	data.Href = types.StringValue(vapp.VApp.HREF)
	data.Description = types.StringValue(vapp.VApp.Description)

	data.ID = types.StringValue(vapp.VApp.ID)
	data.VappID = types.StringValue(vapp.VApp.ID)
	data.VappName = types.StringValue(vapp.VApp.Name)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
