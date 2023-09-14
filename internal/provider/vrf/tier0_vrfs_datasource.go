// Package tier0 provides a Terraform resource to manage Tier-0 VRFs.
package vrf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &tier0VrfsDataSource{}
	_ datasource.DataSourceWithConfigure = &tier0VrfsDataSource{}
)

// NewTier0VrfsDataSource returns a new resource implementing the Tier-0 VRFs data source.
func NewTier0VrfsDataSource() datasource.DataSource {
	return &tier0VrfsDataSource{}
}

type tier0VrfsDataSource struct {
	client *client.CloudAvenue
}

func (d *tier0VrfsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "vrfs"
}

func (d *tier0VrfsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tier0VrfsSchema().GetDataSource(ctx)
}

func (d *tier0VrfsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *tier0VrfsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_tier0_vrfs", d.client.GetOrgName(), metrics.Read)()

	var data tier0VrfsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tier0vrfs, _, err := d.client.APIClient.Tier0Api.GetTier0Vrfs(d.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	for _, name := range tier0vrfs {
		data.Names = append(data.Names, types.StringValue(name))
	}

	// Generate a UUID from the list of names
	data.ID = utils.GenerateUUID(tier0vrfs)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
