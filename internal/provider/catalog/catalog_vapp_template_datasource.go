/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package catalog provides a Terraform datasource.
package catalog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

var (
	_ datasource.DataSource              = &vAppTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &vAppTemplateDataSource{}
)

func NewVAppTemplateDataSource() datasource.DataSource {
	return &vAppTemplateDataSource{}
}

type vAppTemplateDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	catalog  base
}

func (d *vAppTemplateDataSource) Init(ctx context.Context, rm *VAPPTemplateModel) (diags diag.Diagnostics) {
	d.catalog = base{
		name: rm.CatalogName.ValueString(),
		id:   rm.CatalogID.ValueString(),
	}

	d.adminOrg, diags = adminorg.Init(d.client)

	return
}

func (d *vAppTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_vapp_template"
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

func (d *vAppTemplateDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vappTemplateSuperSchema(ctx).GetDataSource(ctx)
}

func (d *vAppTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_catalog_vapp_template", d.client.GetOrgName(), metrics.Read)()

	state := &VAPPTemplateModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalog, err := d.adminOrg.GetAdminCatalogByNameOrId(d.GetIDOrName(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving catalog", err.Error())
		return
	}

	stateUpdated := state.Copy()
	stateUpdated.CatalogID.Set(catalog.AdminCatalog.ID)
	stateUpdated.CatalogName.Set(catalog.AdminCatalog.Name)

	vAppTemplates, err := catalog.QueryVappTemplateList()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp Templates", err.Error())
		return
	}

	for _, vAppTemplate := range vAppTemplates {
		if (state.TemplateID.IsKnown() && vAppTemplate.ID == state.TemplateID.Get()) || (state.TemplateName.IsKnown() && vAppTemplate.Name == state.TemplateName.Get()) {
			// govcd.GetUuidFromHref not working here because the href contains vappTemplate- before the uuid
			// field ID in vAppTemplate attribute is always empty. ID exist in HREF attribute.
			// get last 36 characters of href
			// href ex : http://url.com/xx/xx/xx/vappTemplate-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
			vappTemplateID := urn.Normalize(urn.VAPPTemplate, vAppTemplate.HREF[len(vAppTemplate.HREF)-36:])
			stateUpdated.TemplateID.Set(vappTemplateID.String())
			stateUpdated.ID.Set(vappTemplateID.String())

			stateUpdated.TemplateName.Set(vAppTemplate.Name)

			vappTemplate, err := d.client.Vmware.GetVAppTemplateById(stateUpdated.TemplateID.Get())
			if err != nil {
				resp.Diagnostics.AddError("Error retrieving vApp Template", err.Error())
				return
			}

			// This checks that the vApp Template is synchronized in the catalog
			if _, err = d.client.Vmware.QuerySynchronizedVAppTemplateById(stateUpdated.TemplateID.Get()); err != nil {
				resp.Diagnostics.AddError("Error check vApp Template synchronization", err.Error())
				return
			}

			vmNames := make([]string, 0)
			if vappTemplate.VAppTemplate.Children != nil {
				for _, vm := range vappTemplate.VAppTemplate.Children.VM {
					vmNames = append(vmNames, vm.Name)
				}
			}

			stateUpdated.CreatedAt.Set(vappTemplate.VAppTemplate.DateCreated)
			stateUpdated.Description.Set(vappTemplate.VAppTemplate.Description)

			if len(vmNames) > 0 {
				resp.Diagnostics.Append(stateUpdated.VMNames.Set(ctx, vmNames)...)
				if resp.Diagnostics.HasError() {
					return
				}
			} else {
				stateUpdated.VMNames.SetNull(ctx)
			}

			break
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateUpdated)...)
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
