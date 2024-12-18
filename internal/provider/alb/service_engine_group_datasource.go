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
)

var (
	_ datasource.DataSource              = &albServiceEngineGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &albServiceEngineGroupDataSource{}
)

func NewALBServiceEngineGroupDataSource() datasource.DataSource {
	return &albServiceEngineGroupDataSource{}
}

type albServiceEngineGroupDataSource struct {
	client *client.CloudAvenue
	edgegw edgegw.EdgeGateway
	org    org.Org
}

// Init Initializes the data source.
func (d *albServiceEngineGroupDataSource) Init(ctx context.Context, dm *albServiceEngineGroupModel) (diags diag.Diagnostics) {
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

func (d *albServiceEngineGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_service_engine_group"
}

func (d *albServiceEngineGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = albServiceEngineGroupSchema(ctx).GetDataSource(ctx)
}

func (d *albServiceEngineGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *albServiceEngineGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_alb_service_engine_group", d.client.GetOrgName(), metrics.Read)()

	config := &albServiceEngineGroupModel{}

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
		err    error
		albSEG *v1.EdgeGatewayALBServiceEngineGroupModel
	)
	if config.ID.IsKnown() {
		albSEG, err = d.edgegw.GetALBServiceEngineGroup(config.ID.Get())
	} else {
		albSEG, err = d.edgegw.GetALBServiceEngineGroup(config.Name.Get())
	}
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Service Engine Group", err.Error())
		return
	}

	config.ID.Set(albSEG.ID)
	config.Name.Set(albSEG.Name)
	config.EdgeGatewayID.Set(albSEG.GatewayRef.ID)
	config.EdgeGatewayName.Set(albSEG.GatewayRef.Name)
	config.MaxVirtualServices.SetIntPtr(albSEG.MaxVirtualServices)
	config.ReservedVirtualServices.SetIntPtr(albSEG.MinVirtualServices)
	config.DeployedVirtualServices.SetInt(albSEG.NumDeployedVirtualServices)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
