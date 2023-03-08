// Package iam provides a Terraform datasource.
package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &iamRightDataSource{}
	_ datasource.DataSourceWithConfigure = &iamRightDataSource{}
)

func NewIamRightDataSource() datasource.DataSource {
	return &iamRightDataSource{}
}

type iamRightDataSource struct {
	client *client.CloudAvenue
}

type iamRightDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CategoryID    types.String `tfsdk:"category_id"`
	BundleKey     types.String `tfsdk:"bundle_key"`
	RightType     types.String `tfsdk:"right_type"`
	ImpliedRights types.Set    `tfsdk:"implied_rights"`
}

func (d *iamRightDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "right"
}

func (d *iamRightDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a data source for available rights in Cloud Avenue.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The id of the right.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the right.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description for the right.",
			},
			"category_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The category id for the right.",
			},
			"bundle_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The bundle key for the right.",
			},
			"right_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The right type for the right.",
			},
			"implied_rights": schema.SetNestedAttribute{
				MarkdownDescription: "The list of rights that are implied with this one.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the implied right.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "ID of the implied right.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
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
