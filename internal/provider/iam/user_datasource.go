package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamUserDataSource{}
	_ datasource.DataSourceWithConfigure = &iamUserDataSource{}
	_ user                               = &iamUserDataSource{}
)

// NewiamUserDataSource returns a new Org User data source.
func NewIAMUserDataSource() datasource.DataSource {
	return &iamUserDataSource{}
}

// iamUserDataSource implements the DataSource interface.
type iamUserDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	user     commonUser
}

func (r *iamUserDataSource) Init(_ context.Context, rm *iamUserDataSourceModel) (diags diag.Diagnostics) {
	r.user = commonUser{
		ID:   rm.ID,
		Name: rm.Name,
	}

	r.adminOrg, diags = adminorg.Init(r.client)

	return
}

// Metadata returns the resource type name.
func (d *iamUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_user"
}

// Schema defines the schema for the data source.
func (d *iamUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = userSchema(withDataSource()).GetDataSource()
}

// Configure configures the data source.
func (d *iamUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the data source.
func (d *iamUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	config := &iamUserDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the user by name or ID and return an error if it doesn't exist or there is another error
	user, err := d.GetUser(false)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	// Populate the data source model with the user data

	state := &iamUserDataSourceModel{
		ID:              types.StringValue(user.User.ID),
		Name:            types.StringValue(user.User.Name),
		FullName:        types.StringValue(user.User.FullName),
		RoleName:        types.StringValue(user.User.Role.Name),
		Email:           types.StringValue(user.User.EmailAddress),
		Telephone:       types.StringValue(user.User.Telephone),
		Enabled:         types.BoolValue(user.User.IsEnabled),
		ProviderType:    types.StringValue(user.User.ProviderType),
		DeployedVMQuota: types.Int64Value(int64(user.User.DeployedVmQuota)),
		StoredVMQuota:   types.Int64Value(int64(user.User.StoredVmQuota)),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamUserDataSource) GetUser(refresh bool) (*govcd.OrgUser, error) {
	return r.user.GetUser(r.adminOrg, refresh)
}
