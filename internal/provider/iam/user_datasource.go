package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &iamUserDataSource{}
	_ datasource.DataSourceWithConfigure = &iamUserDataSource{}
)

// NewiamUserDataSource returns a new Org User data source.
func NewIAMUserDataSource() datasource.DataSource {
	return &iamUserDataSource{}
}

// iamUserDataSource implements the DataSource interface.
type iamUserDataSource struct {
	client *client.CloudAvenue
}

// iamUserDataSourceModel is the data source schema.
type iamUserDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	UserName        types.String `tfsdk:"user_name"`
	UserID          types.String `tfsdk:"user_id"`
	Role            types.String `tfsdk:"role"`
	Description     types.String `tfsdk:"description"`
	ProviderType    types.String `tfsdk:"provider_type"`
	FullName        types.String `tfsdk:"full_name"`
	Email           types.String `tfsdk:"email"`
	Telephone       types.String `tfsdk:"telephone"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	IsGroupRole     types.Bool   `tfsdk:"is_group_role"`
	IsLocked        types.Bool   `tfsdk:"is_locked"`
	DeployedVMQuota types.Int64  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   types.Int64  `tfsdk:"stored_vm_quota"`
	GroupNames      types.List   `tfsdk:"group_names"`
}

// Metadata returns the resource type name.
func (d *iamUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "user"
}

// Schema defines the schema for the data source.
func (d *iamUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Avenue iam User data source. This can be used to read users.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID is a `user_id`.",
				Computed:            true,
			},

			// Optional attributes
			"user_name": schema.StringAttribute{
				MarkdownDescription: "The name of the user. Required if `user_id` is not set.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("user_id")),
				},
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user. Required if `user_name` is not set.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("user_name")),
				},
			},

			// Computed attributes
			"role": schema.StringAttribute{
				MarkdownDescription: "Role within the organization",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The user's description",
				Computed:            true,
			},
			"provider_type": schema.StringAttribute{
				MarkdownDescription: "Identity provider type for this this user. One of: 'INTEGRATED', 'SAML', 'OAUTH'. ",
				Computed:            true,
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "The user's full name",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The user's email address",
				Computed:            true,
			},
			"telephone": schema.StringAttribute{
				MarkdownDescription: "The user's telephone number",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "True if the user is enabled and can log in.",
				Computed:            true,
			},
			"is_group_role": schema.BoolAttribute{
				MarkdownDescription: "True if this user has a group role.",
				Computed:            true,
			},
			"is_locked": schema.BoolAttribute{
				MarkdownDescription: "True if the user account has been locked due to too many invalid login attempts.",
				Computed:            true,
			},
			"deployed_vm_quota": schema.Int64Attribute{
				MarkdownDescription: "Quota of vApps that this user can deploy. A value of 0 specifies an unlimited quota.",
				Computed:            true,
			},
			"stored_vm_quota": schema.Int64Attribute{
				MarkdownDescription: "Quota of vApps that this user can store. A value of 0 specifies an unlimited quota.",
				Computed:            true,
			},
			"group_names": schema.ListAttribute{
				MarkdownDescription: "List of group names that this user is a member of.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
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
	var data *iamUserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminOrg, err := d.client.Vmware.GetAdminOrgByNameOrId(d.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	// If user_id is not set, use user_name
	var userNameID string
	if !data.UserID.IsNull() && !data.UserID.IsUnknown() {
		userNameID = data.UserID.ValueString()
	} else {
		userNameID = data.UserName.ValueString()
	}

	// If neither user_id nor user_name is set, return an error
	if userNameID == "" {
		resp.Diagnostics.AddError("Error retrieving user", "user_name or user_id must be set")
		return
	}

	// Get the user by name or ID and return an error if it doesn't exist or there is another error
	user, err := adminOrg.GetUserByNameOrId(userNameID, false)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	// Populate the data source model with the user data
	data = &iamUserDataSourceModel{
		ID:              types.StringValue(user.User.ID),
		UserName:        types.StringValue(user.User.Name),
		UserID:          types.StringValue(user.User.ID),
		FullName:        types.StringValue(user.User.FullName),
		Role:            types.StringValue(user.User.Role.Name),
		Email:           types.StringValue(user.User.EmailAddress),
		Telephone:       types.StringValue(user.User.Telephone),
		Enabled:         types.BoolValue(user.User.IsEnabled),
		Description:     types.StringValue(user.User.Description),
		DeployedVMQuota: types.Int64Value(int64(user.User.DeployedVmQuota)),
		StoredVMQuota:   types.Int64Value(int64(user.User.StoredVmQuota)),
	}

	listGroupNames := make([]string, 0)
	for _, group := range user.User.GroupReferences.GroupReference {
		listGroupNames = append(listGroupNames, group.Name)
	}

	l, errListVal := types.ListValueFrom(ctx, types.StringType, listGroupNames)
	resp.Diagnostics.Append(errListVal...)

	if resp.Diagnostics.HasError() {
		data.GroupNames = types.ListNull(types.StringType)
	} else {
		data.GroupNames = l
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
