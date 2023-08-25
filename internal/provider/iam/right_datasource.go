// Package iam provides a Terraform datasource.
package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &iamRightDataSource{}
	_ datasource.DataSourceWithConfigure = &iamRightDataSource{}
)

func NewIAMRightDataSource() datasource.DataSource {
	return &iamRightDataSource{}
}

type iamRightDataSource struct {
	client *client.CloudAvenue
}

func (d *iamRightDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "right"
}

// ! Convert to iam_rightS

func (d *iamRightDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = iamRightSuperSchema(ctx).GetDataSource(ctx)
}

func (d *iamRightDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *iamRightDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_iam_right", d.client.GetOrgName(), metrics.Read)()

	var data RightModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	right, err := d.client.Vmware.Client.GetRightByName(data.Name.Get())
	if err != nil {
		resp.Diagnostics.AddError("This right does not exist", data.Name.Get())
		return
	}

	// Set the data to be returned
	data.ID.Set(right.ID)
	data.Name.Set(right.Name)
	data.Description.Set(right.Description)
	data.CategoryID.Set(right.Category)
	data.BundleKey.Set(right.BundleKey)
	data.RightType.Set(right.RightType)

	impliedRights := make(RightModelImpliedRights, 0)
	for _, ir := range right.ImpliedRights {
		p := RightModelImpliedRight{}
		p.ID.Set(ir.ID)
		p.Name.Set(ir.Name)
		impliedRights = append(impliedRights, p)
	}

	resp.Diagnostics.Append(data.ImpliedRights.Set(ctx, impliedRights)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
