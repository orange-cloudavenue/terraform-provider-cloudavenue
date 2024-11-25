package edgegw //nolint:dupl // This is a datasource, it is normal to have similar code to the other datasource.

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
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
}

// Init Initializes the data source.
func (d *appPortProfileDataSource) Init(ctx context.Context, dm *AppPortProfileModelADatasource) (diags diag.Diagnostics) {
	d.org, diags = org.Init(d.client)
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

	config := &AppPortProfileModelADatasource{}

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

	// If read function is identical to the resource, you can use the following code:
	s := &appPortProfileResource{
		client: d.client,
		org:    d.org,
	}

	// Read data from the API
	data, found, diags := s.read(ctx, &AppPortProfileModel{
		ID:              config.ID,
		Name:            config.Name,
		Description:     config.Description,
		EdgeGatewayID:   supertypes.NewStringNull(),
		EdgeGatewayName: supertypes.NewStringNull(),
		AppPorts:        config.AppPorts,
	})
	if !found {
		if config.ID.IsKnown() {
			resp.Diagnostics.AddError("Not found", fmt.Sprintf("App Port Profile ID %q not found", config.ID))
		} else {
			resp.Diagnostics.AddError("Not found", fmt.Sprintf("App Port Profile name %q not found", config.Name))
		}
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write data into the model
	config.ID = data.ID
	config.Name = data.Name
	config.Description = data.Description
	config.AppPorts = data.AppPorts

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
