// Package alb provides a Terraform datasource.
package alb

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

var (
	_ datasource.DataSource              = &albPoolDataSource{}
	_ datasource.DataSourceWithConfigure = &albPoolDataSource{}
	_ albPool                            = &albPoolDataSource{}
)

func NewAlbPoolDataSource() datasource.DataSource {
	return &albPoolDataSource{}
}

type albPoolDataSource struct {
	client  *client.CloudAvenue
	org     org.Org
	vdc     vdc.VDC
	albPool base
}

func (d *albPoolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_pool"
}

func (d *albPoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = albPoolSchema().GetDataSource(ctx)
}

func (d *albPoolDataSource) Init(ctx context.Context, dm *albPoolModel) (diags diag.Diagnostics) {
	d.albPool = base{
		name: dm.Name.ValueString(),
		id:   dm.ID.ValueString(),
	}

	d.vdc, diags = vdc.Init(d.client, dm.VDC)
	if diags.HasError() {
		return
	}

	d.org, diags = org.Init(d.client)
	return
}

func (d *albPoolDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *albPoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data  *albPoolModel
		diags diag.Diagnostics
	)
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get albPool.
	albPool, err := d.GetAlbPool(data.EdgeGatewayID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find ALB Pool", err.Error())
		return
	}

	// Set data
	data.ID = types.StringValue(albPool.NsxtAlbPool.ID)
	data.Description = types.StringValue(albPool.NsxtAlbPool.Description)
	data.Enabled = types.BoolValue(*albPool.NsxtAlbPool.Enabled)
	data.Algorithm = types.StringValue(albPool.NsxtAlbPool.Algorithm)
	data.DefaultPort = types.Int64Value(int64(*albPool.NsxtAlbPool.DefaultPort))
	data.GracefulTimeoutPeriod = types.Int64Value(int64(*albPool.NsxtAlbPool.GracefulTimeoutPeriod))
	data.PassiveMonitoringEnabled = types.BoolValue(*albPool.NsxtAlbPool.PassiveMonitoringEnabled)

	// Init Set and List type
	data.Member = types.SetNull(types.ObjectType{AttrTypes: memberAttrTypes})
	data.HealthMonitor = types.SetNull(types.ObjectType{AttrTypes: healthMonitorAttrTypes})
	data.PersistenceProfile = types.ListNull(types.ObjectType{AttrTypes: persistenceProfileAttrTypes})

	// Set members
	members := []member{}
	if len(albPool.NsxtAlbPool.Members) > 0 {
		for _, albMember := range albPool.NsxtAlbPool.Members {
			members = append(members, member{
				Enabled:   types.BoolValue(albMember.Enabled),
				IPAddress: types.StringValue(albMember.IpAddress),
				Port:      types.Int64Value(int64(albMember.Port)),
				Ratio:     types.Int64Value(int64(*albMember.Ratio)),
			})
		}
	}

	data.Member, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: memberAttrTypes}, members)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set health monitors.
	healtMonitors := []healthMonitor{}
	if len(albPool.NsxtAlbPool.HealthMonitors) > 0 {
		for _, albHealthMonitor := range albPool.NsxtAlbPool.HealthMonitors {
			healtMonitors = append(healtMonitors, healthMonitor{
				Type: types.StringValue(albHealthMonitor.Type),
				Name: types.StringValue(albHealthMonitor.Name),
			})
		}
	}

	data.HealthMonitor, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: healthMonitorAttrTypes}, healtMonitors)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set persistence profile
	persistenceProfiles := []persistenceProfile{}
	if albPool.NsxtAlbPool.PersistenceProfile != nil {
		persistenceProfiles = append(persistenceProfiles, persistenceProfile{
			Type:  types.StringValue(albPool.NsxtAlbPool.PersistenceProfile.Type),
			Value: types.StringValue(albPool.NsxtAlbPool.PersistenceProfile.Value),
		})
	}

	data.PersistenceProfile, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: persistenceProfileAttrTypes}, persistenceProfiles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// GetID returns the ID of the albPool.
func (d *albPoolDataSource) GetID() string {
	return d.albPool.id
}

// GetName returns the name of the albPool.
func (d *albPoolDataSource) GetName() string {
	return d.albPool.name
}

// GetAlbPool returns the govcd.NsxtAlbPool.
func (d *albPoolDataSource) GetAlbPool(edgegwID string) (albPool *govcd.NsxtAlbPool, err error) {
	if d.GetID() != "" {
		albPool, err = d.client.Vmware.GetAlbPoolById(d.GetID())
	} else {
		nsxtEdge, err := d.org.GetNsxtEdgeGatewayById(edgegwID)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve Edge gateway with ID '%s'", edgegwID)
		}
		albPool, err = d.client.Vmware.GetAlbPoolByName(nsxtEdge.EdgeGateway.ID, d.GetName())
	}
	return
}
