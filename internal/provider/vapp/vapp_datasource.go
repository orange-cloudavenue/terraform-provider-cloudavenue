// Package vapp provides a Terraform resource to manage vApps.
package vapp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

var (
	_ datasource.DataSource              = &vappDataSource{}
	_ datasource.DataSourceWithConfigure = &vappDataSource{}
)

const vAppUnknownStatus = "-unknown-status-"

// NewVappDataSource returns a new vApp data source.
func NewVappDataSource() datasource.DataSource {
	return &vappDataSource{}
}

type vappDataSource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	vapp   vapp.VAPP
}

type vappDataSourceModel struct {
	VAppName        types.String                  `tfsdk:"name"`
	VAppID          types.String                  `tfsdk:"id"`
	VDC             types.String                  `tfsdk:"vdc"`
	Description     types.String                  `tfsdk:"description"`
	Href            types.String                  `tfsdk:"href"`
	GuestProperties map[types.String]types.String `tfsdk:"guest_properties"`
	StatusCode      types.Int64                   `tfsdk:"status_code"`
	StatusText      types.String                  `tfsdk:"status_text"`
	Lease           types.Object                  `tfsdk:"lease"`
}

func (d *vappDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *vappDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Avenue vApp data source. This can be used to reference vApps.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A name for the vApp, unique within the VDC. Required if `id` is not set.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
				},
			},
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of vApp to use, unique within the VDC. Required if `name` is not set.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
				},
			},
			"vdc": vdc.Schema(),
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

func (d *vappDataSource) Init(ctx context.Context, dm *vappDataSourceModel) (diags diag.Diagnostics) {
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
	var data vappDataSourceModel

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

	// Get guest properties
	guestProperties, err := d.vapp.GetProductSectionList()
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

	leaseInfo, err := d.vapp.GetLease()
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

	statusText, err := d.vapp.GetStatus()
	if err != nil {
		statusText = vAppUnknownStatus
	}

	data.StatusCode = types.Int64Value(int64(d.vapp.GetStatusCode()))
	data.StatusText = types.StringValue(statusText)
	data.Href = types.StringValue(d.vapp.GetHREF())
	data.Description = types.StringValue(d.vapp.GetDescription())

	data.VAppID = types.StringValue(d.vapp.GetID())
	data.VAppName = types.StringValue(d.vapp.GetName())
	data.VDC = types.StringValue(d.vdc.GetName())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
