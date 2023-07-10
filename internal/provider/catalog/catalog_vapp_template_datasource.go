// Package catalog provides a Terraform datasource.
package catalog

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

var (
	_ datasource.DataSource              = &vAppTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &vAppTemplateDataSource{}
	_ catalog                            = &vAppTemplateDataSource{}
)

func NewVAppTemplateDataSource() datasource.DataSource {
	return &vAppTemplateDataSource{}
}

type vAppTemplateDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	catalog  base
}

func (d *vAppTemplateDataSource) Init(ctx context.Context, rm *vAppTemplateDataSourceModel) (diags diag.Diagnostics) {
	d.catalog = base{
		name: rm.CatalogName.ValueString(),
		id:   rm.CatalogID.ValueString(),
	}

	d.adminOrg, diags = adminorg.Init(d.client)

	return
}

func (d *vAppTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "vapp_template"
}

func (d *vAppTemplateDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The `catalog_vapp_template` datasource provides information about a vApp Template in a catalog.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the vApp Template",
				Computed:            true,
			},
			"template_name": schema.StringAttribute{
				MarkdownDescription: "Name of the vApp Template. Required if `template_id` is not set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("template_name"), path.MatchRoot("template_id")),
				},
			},
			"template_id": schema.StringAttribute{
				MarkdownDescription: "ID of the vApp Template. Required if `template_name` is not set.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("template_name"), path.MatchRoot("template_id")),
				},
			},
			catalogID:   mediaSchema().GetDataSource(ctx).Attributes[catalogID],
			catalogName: mediaSchema().GetDataSource(ctx).Attributes[catalogName],
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

func (d *vAppTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vAppTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := &vAppTemplateDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalog, err := d.GetCatalog()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving catalog", err.Error())
		return
	}

	if state.TemplateID.IsNull() || state.TemplateID.IsUnknown() {
		vAppTemplates, err := catalog.QueryVappTemplateList()
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving vApp Templates", err.Error())
			return
		}

		var href string

		for _, vAppTemplate := range vAppTemplates {
			if vAppTemplate.Name == state.TemplateName.ValueString() {
				state.TemplateID = types.StringValue(vAppTemplate.ID)
				href = vAppTemplate.HREF
				break
			}
		}

		if state.TemplateID.ValueString() == "" {
			// govcd.GetUuidFromHref not working here because the href contains vappTemplate- before the uuid

			// get last 36 characters of href
			uuid := href[len(href)-36:]

			if uuid != "" {
				state.TemplateID = types.StringValue(uuid)
			} else {
				resp.Diagnostics.AddError("Error retrieving vApp Template", fmt.Sprintf("vApp Template '%s' not found", state.TemplateName.ValueString()))
				return
			}
		}

		state.ID = state.TemplateID
	}

	vappTemplate, err := d.client.Vmware.GetVAppTemplateById(state.TemplateID.ValueString())
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

	updatedState := state
	updatedState.Description = types.StringValue(vappTemplate.VAppTemplate.Description)
	updatedState.CreatedAt = types.StringValue(vappTemplate.VAppTemplate.DateCreated)
	updatedState.VMNames = vmS

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *vAppTemplateDataSource) GetID() string {
	return d.catalog.id
}

// GetName returns the name of the catalog.
func (d *vAppTemplateDataSource) GetName() string {
	return d.catalog.name
}

// GetIDOrName returns the ID if it is set, otherwise it returns the name.
func (d *vAppTemplateDataSource) GetIDOrName() string {
	if d.GetID() != "" {
		return d.GetID()
	}
	return d.GetName()
}

// GetCatalog returns the govcd.Catalog.
func (d *vAppTemplateDataSource) GetCatalog() (*govcd.AdminCatalog, error) {
	return d.adminOrg.GetAdminCatalogByNameOrId(d.GetIDOrName(), true)
}
