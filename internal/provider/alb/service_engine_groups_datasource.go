// Package alb provides a Terraform datasource.
package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &serviceEngineGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceEngineGroupsDataSource{}
)

func NewServiceEngineGroupsDataSource() datasource.DataSource {
	return &serviceEngineGroupsDataSource{}
}

type serviceEngineGroupsDataSource struct {
	client   *client.CloudAvenue
	edgegwlb edgeloadbalancer.Client
}

// Init Initializes the data source.
func (d *serviceEngineGroupsDataSource) Init(ctx context.Context, dm *serviceEngineGroupsModel) (diags diag.Diagnostics) {
	var err error

	d.edgegwlb, err = edgeloadbalancer.NewClient()
	if err != nil {
		diags.AddError("Error creating edge load balancer client", err.Error())
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
	edge, err := d.client.CAVSDK.V1.EdgeGateway.Get(func() string {
		if config.EdgeGatewayID.IsKnown() {
			return config.EdgeGatewayID.Get()
		}

		return config.EdgeGatewayName.Get()
	}())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	config.EdgeGatewayID.Set(urn.Normalize(urn.Gateway, edge.GetID()).String())
	config.EdgeGatewayName.Set(edge.GetName())

	// Get Service Engine Groups
	albSEGs, err := d.edgegwlb.ListServiceEngineGroups(ctx, config.EdgeGatewayID.Get())
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

	config.ID.Set(urn.Normalize(urn.Gateway, config.EdgeGatewayID.Get()).String())
	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
