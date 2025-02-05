// Package alb provides a Terraform datasource.
package alb

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &serviceEngineGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceEngineGroupDataSource{}
)

func NewServiceEngineGroupDataSource() datasource.DataSource {
	return &serviceEngineGroupDataSource{}
}

type serviceEngineGroupDataSource struct {
	client   *client.CloudAvenue
	edgegwlb edgeloadbalancer.Client
}

// Init Initializes the data source.
func (d *serviceEngineGroupDataSource) Init(ctx context.Context, dm *serviceEngineGroupModel) (diags diag.Diagnostics) {
	var err error

	d.edgegwlb, err = edgeloadbalancer.NewClient()
	if err != nil {
		diags.AddError("Error creating edge load balancer client", err.Error())
	}

	return
}

func (d *serviceEngineGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_service_engine_group"
}

func (d *serviceEngineGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = serviceEngineGroupSchema(ctx).GetDataSource(ctx)
}

func (d *serviceEngineGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serviceEngineGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_alb_service_engine_group", d.client.GetOrgName(), metrics.Read)()

	config := &serviceEngineGroupModel{}

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

	// Read data from the API
	data, found, diags := d.read(ctx, config)
	if !found {
		diags.AddError("Error Not Found", "The Service Engine Group was not found")
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *serviceEngineGroupDataSource) read(ctx context.Context, dm *serviceEngineGroupModel) (data *serviceEngineGroupModel, found bool, diags diag.Diagnostics) {
	data = &serviceEngineGroupModel{}

	// Get ServiceEngineGroup
	var (
		err    error
		albSEG *edgeloadbalancer.ServiceEngineGroupModel
	)

	if !dm.EdgeGatewayID.IsKnown() {
		edge, err := d.client.CAVSDK.V1.EdgeGateway.Get(dm.EdgeGatewayName.Get())
		if err != nil {
			diags.AddError("Error retrieving Edge Gateway", err.Error())
			// True because there was an error on fetch the edge gateway and not the service engine group
			return nil, true, diags
		}

		dm.EdgeGatewayID.Set(urn.Normalize(urn.Gateway, edge.GetID()).String())
	}

	if dm.ID.IsKnown() {
		albSEG, err = d.edgegwlb.GetServiceEngineGroup(ctx, dm.EdgeGatewayID.Get(), dm.ID.Get())
	} else {
		albSEG, err = d.edgegwlb.GetServiceEngineGroup(ctx, dm.EdgeGatewayID.Get(), dm.Name.Get())
	}
	if err != nil {
		if commoncloudavenue.IsNotFound(err) || govcd.IsNotFound(err) {
			return nil, false, diags
		}
		diags.AddError("Error retrieving Service Engine Group", err.Error())
		return nil, true, diags
	}

	data.ID.Set(albSEG.ID)
	data.Name.Set(albSEG.Name)
	data.EdgeGatewayID.Set(albSEG.GatewayRef.ID)
	data.EdgeGatewayName.Set(albSEG.GatewayRef.Name)
	data.MaxVirtualServices.SetIntPtr(albSEG.MaxVirtualServices)
	data.ReservedVirtualServices.SetIntPtr(albSEG.MinVirtualServices)
	data.DeployedVirtualServices.SetInt(albSEG.NumDeployedVirtualServices)

	return data, true, diags
}
