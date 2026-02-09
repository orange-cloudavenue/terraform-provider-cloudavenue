/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &AppPortProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &AppPortProfileDataSource{}
)

func NewAppPortProfileDataSource() datasource.DataSource {
	return &AppPortProfileDataSource{}
}

type AppPortProfileDataSource struct {
	client   *client.CloudAvenue
	vdcGroup *v1.VDCGroup
}

func (d *AppPortProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_app_port_profile"
}

// Init Initializes the resource.
func (d *AppPortProfileDataSource) Init(_ context.Context, rm *AppPortProfileModelDatasource) (diags diag.Diagnostics) {
	var err error

	idOrName := rm.VDCGroupName.Get()
	if rm.VDCGroupID.IsKnown() && urn.IsVDCGroup(rm.VDCGroupID.Get()) {
		// Use the ID
		idOrName = rm.VDCGroupID.Get()
	}

	d.vdcGroup, err = d.client.CAVSDK.V1.VDC().GetVDCGroup(idOrName)
	if err != nil {
		diags.AddError("Error retrieving VDC Group", err.Error())
		return diags
	}
	return diags
}

func (d *AppPortProfileDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = appPortProfileSchema(ctx).GetDataSource(ctx)
}

func (d *AppPortProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AppPortProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdcg_app_port_profile", d.client.GetOrgName(), metrics.Read)()

	config := &AppPortProfileModelDatasource{}

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

	nameOrID := config.Name.Get()
	if config.ID.IsKnown() {
		nameOrID = config.ID.Get()
	}

	var appP *v1.FirewallGroupAppPortProfileModelResponse

	appPortProfiles, err := d.vdcGroup.FindFirewallAppPortProfile(nameOrID)
	if err != nil {
		if errors.IsNotFound(err) {
			resp.Diagnostics.AddError("App Port Profile not found", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error reading App Port Profile", err.Error())
		return
	}

	if len(appPortProfiles.AppPortProfiles) > 1 && !config.Scope.IsKnown() {
		resp.Diagnostics.AddError("Multiple App Port Profiles found", "Multiple App Port Profiles found with the same name.")
		resp.Diagnostics.AddError("Details of the App Port Profiles", func() (msg string) {
			for _, appPortProfile := range appPortProfiles.AppPortProfiles {
				msg += fmt.Sprintf("ID: %s\nName: %s\nScope: %s\n\n", appPortProfile.ID, appPortProfile.Name, appPortProfile.Scope)
			}

			msg += "Please provide the ID of the App Port Profile to uniquely identify it or add the scope."
			return msg
		}())
		return
	}

	if len(appPortProfiles.AppPortProfiles) >= 1 && config.Scope.IsKnown() {
		// Find the App Port Profile with the correct scope
		for _, appPortProfile := range appPortProfiles.AppPortProfiles {
			if appPortProfile.Scope == v1.FirewallGroupAppPortProfileModelScope(config.Scope.Get()) {
				appP = appPortProfile
				break
			}
		}
		if appP == nil {
			resp.Diagnostics.AddError("App Port Profile not found", "App Port Profile not found with the specified scope.")
			return
		}
	}

	if len(appPortProfiles.AppPortProfiles) == 1 && !config.Scope.IsKnown() {
		appP = appPortProfiles.AppPortProfiles[0]
	}

	appPorts := make([]*AppPortProfileModelAppPort, len(appP.ApplicationPorts))
	for index, singlePort := range appP.ApplicationPorts {
		ap := &AppPortProfileModelAppPort{
			Protocol: supertypes.NewStringNull(),
			Ports:    supertypes.NewSetValueOfNull[string](ctx),
		}

		ap.Protocol.Set(string(singlePort.Protocol))
		if singlePort.Protocol == v1.FirewallGroupAppPortProfileModelPortProtocolTCP || singlePort.Protocol == v1.FirewallGroupAppPortProfileModelPortProtocolUDP {
			// DestinationPorts is optional
			if len(singlePort.DestinationPorts) > 0 {
				resp.Diagnostics.Append(ap.Ports.Set(ctx, singlePort.DestinationPorts)...)
				if resp.Diagnostics.HasError() {
					return
				}
			}
		}
		appPorts[index] = ap
	}

	stateRefreshed := config.Copy()

	stateRefreshed.ID.Set(appP.ID)
	stateRefreshed.Name.Set(appP.Name)
	stateRefreshed.Description.Set(appP.Description)
	stateRefreshed.AppPorts.Set(ctx, appPorts)
	stateRefreshed.VDCGroupID.Set(d.vdcGroup.GetID())
	stateRefreshed.VDCGroupName.Set(d.vdcGroup.GetName())
	stateRefreshed.Scope.Set(string(appP.Scope))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}
