package edgegw //nolint:dupl // This is a datasource, it is normal to have similar code to the other datasource.

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

var (
	_ datasource.DataSource              = &appPortProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &appPortProfileDataSource{}
)

func NewAppPortProfileDataSource() datasource.DataSource {
	return &appPortProfileDataSource{}
}

type appPortProfileDataSource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the data source.
func (d *appPortProfileDataSource) Init(ctx context.Context, dm *AppPortProfileModelDatasource) (diags diag.Diagnostics) {
	var err error

	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return
	}

	// Retrieve VDC from edge gateway
	d.edgegw, err = d.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(dm.EdgeGatewayID.Get()),
		Name: types.StringValue(dm.EdgeGatewayName.Get()),
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	return
}

func (d *appPortProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_app_port_profile"
}

func (d *appPortProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = appPortProfilesSchema(ctx).GetDataSource(ctx)
}

func (d *appPortProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *appPortProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_edgegateway_app_port_profile", d.client.GetOrgName(), metrics.Read)()

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

	appPortProfiles, err := d.edgegw.FindFirewallAppPortProfile(nameOrID)
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
			return
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
	stateRefreshed.EdgeGatewayID.Set(d.edgegw.EdgeGateway.ID)
	stateRefreshed.EdgeGatewayName.Set(d.edgegw.EdgeGateway.Name)
	stateRefreshed.Scope.Set(string(appP.Scope))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}
