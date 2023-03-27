// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &edgeGatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &edgeGatewaysDataSource{}
)

// NewEdgeGatewaysDataSource returns a new resource implementing the edge_gateways data source.
func NewEdgeGatewaysDataSource() datasource.DataSource {
	return &edgeGatewaysDataSource{}
}

type edgeGatewaysDataSource struct {
	client *client.CloudAvenue
}

type edgeGatewaysDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	EdgeGateways types.List   `tfsdk:"edge_gateways"`
}

var edgeGatewayDataSourceModelAttrTypes = map[string]attr.Type{
	"tier0_vrf_name": types.StringType,
	"name":           types.StringType,
	"id":             types.StringType,
	"owner_type":     types.StringType,
	"owner_name":     types.StringType,
	"description":    types.StringType,
	"lb_enabled":     types.BoolType,
}

func (d *edgeGatewaysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *edgeGatewaysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The edge gateways data source show the list of edge gateways of an organization.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"edge_gateways": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of Edge Gateways.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tier0_vrf_name": edgegwSchema().GetDataSource(ctx).Attributes["tier0_vrf_name"],
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the Edge Gateway.",
							Computed:            true,
						},
						"id":          edgegwSchema().GetDataSource(ctx).Attributes["id"],
						"owner_type":  edgegwSchema().GetDataSource(ctx).Attributes["owner_type"],
						"owner_name":  edgegwSchema().GetDataSource(ctx).Attributes["owner_name"],
						"description": edgegwSchema().GetDataSource(ctx).Attributes["description"],
						"lb_enabled":  edgegwSchema().GetDataSource(ctx).Attributes["lb_enabled"],
					},
				},
			},
		},
	}
}

func (d *edgeGatewaysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *edgeGatewaysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data  edgeGatewaysDataSourceModel
		names []string
	)
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	gateways, httpR, err := d.client.APIClient.EdgeGatewaysApi.GetEdges(d.client.Auth)
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}
	if x := helpers.CheckAPIError(err, httpR); x != nil {
		if !x.IsNotFound() {
			resp.Diagnostics.Append(x.GetTerraformDiagnostic())
			return
		}
		// Is Not Found
		data.EdgeGateways = types.ListNull(types.ObjectType{AttrTypes: edgeGatewayDataSourceModelAttrTypes})
		data.ID = types.StringNull()
	} else {
		var diag diag.Diagnostics
		gws := make([]edgeGatewayDataSourceModel, 0)
		for _, gw := range gateways {
			// Get LoadBalancing state.
			gatewaysLoadBalancing, httpR, err := d.client.APIClient.EdgeGatewaysApi.GetEdgeLoadBalancing(d.client.Auth, gw.EdgeId)
			if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
				defer httpR.Body.Close()
				resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
				if resp.Diagnostics.HasError() {
					return
				}
			}
			gws = append(gws, edgeGatewayDataSourceModel{
				Tier0VrfID:          types.StringValue(gw.Tier0VrfId),
				Name:                types.StringValue(gw.EdgeName),
				ID:                  types.StringValue(common.NormalizeID("urn:vcloud:gateway:", gw.EdgeId)),
				OwnerType:           types.StringValue(gw.OwnerType),
				OwnerName:           types.StringValue(gw.OwnerName),
				Description:         types.StringValue(gw.Description),
				EnableLoadBalancing: types.BoolValue((gatewaysLoadBalancing.Enabled)),
			})

			names = append(names, gw.EdgeName)
		}

		data.EdgeGateways, diag = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: edgeGatewayDataSourceModelAttrTypes}, gws)
		resp.Diagnostics.Append(diag...)
		data.ID = utils.GenerateUUID(names)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
