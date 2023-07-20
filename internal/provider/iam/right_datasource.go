// Package iam provides a Terraform datasource.
package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
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
	resp.Schema = iamRightSchema()
}

type ImpliedRightsModel struct {
	Name types.String `tfsdk:"name"`
	ID   types.String `tfsdk:"id"`
}

func iamRightDataImpliedRightsAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"id":   types.StringType,
	}
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
	var data iamRightDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read right name from the provider
	rightName := data.Name.ValueString()

	var diag diag.Diagnostics

	right, err := d.client.Vmware.Client.GetRightByName(rightName)
	if err != nil {
		resp.Diagnostics.AddError("This right does not exist", rightName)
		return
	}

	// Set the data to be returned
	data.ID = types.StringValue(right.ID)
	data.Name = types.StringValue(right.Name)
	data.Description = types.StringValue(right.Description)
	data.CategoryID = types.StringValue(right.Category)
	data.BundleKey = types.StringValue(right.BundleKey)
	data.RightType = types.StringValue(right.RightType)

	var impliedRights []ImpliedRightsModel
	for _, ir := range right.ImpliedRights {
		p := ImpliedRightsModel{
			Name: types.StringValue(ir.Name),
			ID:   types.StringValue(ir.ID),
		}
		impliedRights = append(impliedRights, p)
	}

	data.ImpliedRights, diag = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: iamRightDataImpliedRightsAttrType()}, impliedRights)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
