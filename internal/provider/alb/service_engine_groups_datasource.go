// Package alb provides a Terraform datasource.
package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &serviceEngineGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceEngineGroupsDataSource{}
)

func NewServiceEngineGroupsDataSource() datasource.DataSource {
	return &serviceEngineGroupsDataSource{}
}

type serviceEngineGroupsDataSource struct {
	client *client.CloudAvenue
	edgegw edgegw.EdgeGateway
	org    org.Org
}

// Init Initializes the data source.
func (d *serviceEngineGroupsDataSource) Init(ctx context.Context, dm *serviceEngineGroupsModel) (diags diag.Diagnostics) {
	var err error

	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return
	}

	// Retrieve VDC from edge gateway
	d.edgegw, err = d.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   dm.EdgeGatewayID.StringValue,
		Name: dm.EdgeGatewayName.StringValue,
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	return
}

func (d *serviceEngineGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_service_engine_groups"
}

func (d *serviceEngineGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = serviceEngineGroupsSchema(ctx).GetDataSource(ctx)
}

func (d *serviceEngineGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serviceEngineGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_alb_service_engine_groups", d.client.GetOrgName(), metrics.Read)()

	config := &serviceEngineGroupsModel{}

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

	// Get ServiceEngineGroup
	var (
		err     error
		albSEGs []*v1.EdgeGatewayALBServiceEngineGroupModel
	)

	// Get Service Engine Groups
	albSEGs, err = d.edgegw.ListALBServiceEngineGroups()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Service Engine Groups", err.Error())
		return
	}

	// Set SDK model to Terraform model
	seg := make([]*serviceEngineGroupModel, 0)
	for _, albSEG := range albSEGs {
		x := &serviceEngineGroupModel{}
		x.ID.Set(albSEG.ID)
		x.Name.Set(albSEG.Name)
		x.EdgeGatewayID.Set(albSEG.GatewayRef.ID)
		x.EdgeGatewayName.Set(albSEG.GatewayRef.Name)
		x.MaxVirtualServices.SetIntPtr(albSEG.MaxVirtualServices)
		x.ReservedVirtualServices.SetIntPtr(albSEG.MinVirtualServices)
		x.DeployedVirtualServices.SetInt(albSEG.NumDeployedVirtualServices)
		seg = append(seg, x)
	}

	// Set config
	resp.Diagnostics.Append(config.ServiceEngineGroups.Set(ctx, seg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.ID.Set(utils.GenerateUUID(fmt.Sprint(d.edgegw.EdgeName + d.edgegw.EdgeID)).ValueString())
	config.EdgeGatewayID.Set(d.edgegw.EdgeGateway.ID)
	config.EdgeGatewayName.Set(d.edgegw.EdgeGateway.Name)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
