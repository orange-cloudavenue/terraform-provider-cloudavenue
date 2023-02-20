// Package catalog provides a Terraform datasource.
package catalog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &catalogVappTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogVappTemplateDataSource{}
)

func NewCatalogVappTemplateDataSource() datasource.DataSource {
	return &catalogVappTemplateDataSource{}
}

type catalogVappTemplateDataSource struct {
	client *client.CloudAvenue
}

type catalogVappTemplateDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	VappName    types.String `tfsdk:"vapp_name"`
	VappID      types.String `tfsdk:"vapp_id"`
	CatalogID   types.String `tfsdk:"catalog_id"`
	CatalogName types.String `tfsdk:"catalog_name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	VMNames     types.List   `tfsdk:"vm_names"`
}

func (d *catalogVappTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + "catalog_vapp_template"
}

func (d *catalogVappTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The catalog_vapp_template datasource provides information about a vApp Template in a catalog.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID is a `vapp_id` of the vApp Template",
				Computed:            true,
			},
			"vapp_name": schema.StringAttribute{
				MarkdownDescription: "Name of the vApp Template. Required if `vapp_id` is not set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
				},
			},
			"vapp_id": schema.StringAttribute{
				MarkdownDescription: "ID of the vApp Template. Required if `vapp_name` is not set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
				},
			},
			"catalog_id": schema.StringAttribute{
				MarkdownDescription: "ID of the catalog containing the vApp Template. Required if `catalog_name` is not set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_id"), path.MatchRoot("catalog_name")),
				},
			},
			"catalog_name": schema.StringAttribute{
				MarkdownDescription: "Name of the catalog containing the vApp Template. Required if `catalog_id` is not set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_id"), path.MatchRoot("catalog_name")),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the vApp Template",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation date of the vApp Template",
				Computed:            true,
			},
			"vm_names": schema.ListAttribute{
				MarkdownDescription: "Set of VM names within the vApp template",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *catalogVappTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *catalogVappTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data, state catalogVappTemplateDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	adminOrg, err := d.client.Vmware.GetAdminOrgByNameOrId(d.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	// Get VAPP Template by Catalog ID or Name
	var (
		catalog           *govcd.Catalog
		catalogByNameOrID string
		vappByNameOrID    string
	)

	if !data.CatalogID.IsNull() {
		tflog.Info(ctx, fmt.Sprintf("Catalog ID: %s", data.CatalogID.ValueString()))
		catalogByNameOrID = data.CatalogID.ValueString()
	} else {
		tflog.Info(ctx, fmt.Sprintf("Catalog Name: %s", data.CatalogName.ValueString()))
		catalogByNameOrID = data.CatalogName.ValueString()
	}

	catalog, err = adminOrg.GetCatalogByNameOrId(catalogByNameOrID, false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Catalog", err.Error())
		return
	}

	if !data.VappID.IsNull() {
		vappByNameOrID = data.VappID.ValueString()
	} else {
		vappByNameOrID = data.VappName.ValueString()
	}

	// catalog.QueryVappTemplateWithName()
	vappTemplate, err := catalog.GetVAppTemplateByNameOrId(vappByNameOrID, false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp Template", err.Error())
		return
	}

	// This checks that the vApp Template is synchronized in the catalog
	_, err = d.client.Vmware.QuerySynchronizedVAppTemplateById(vappTemplate.VAppTemplate.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error check vApp Template synchronization", err.Error())
		return
	}

	var vmNames []attr.Value
	if vappTemplate.VAppTemplate.Children != nil {
		for _, vm := range vappTemplate.VAppTemplate.Children.VM {
			vmNames = append(vmNames, types.StringValue(vm.Name))
		}
	}
	vmS, diag := basetypes.NewListValue(types.StringType, vmNames)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	state = catalogVappTemplateDataSourceModel{
		ID:          types.StringValue(vappTemplate.VAppTemplate.ID),
		VappName:    types.StringValue(vappTemplate.VAppTemplate.Name),
		VappID:      types.StringValue(vappTemplate.VAppTemplate.ID),
		CatalogID:   types.StringValue(catalog.Catalog.ID),
		CatalogName: types.StringValue(catalog.Catalog.Name),
		Description: types.StringValue(vappTemplate.VAppTemplate.Description),
		CreatedAt:   types.StringValue(vappTemplate.VAppTemplate.DateCreated),
		VMNames:     vmS,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
