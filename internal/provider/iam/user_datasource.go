/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/iam"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

// NewuserDataSource returns a new Org User data source.
func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

// userDataSource implements the DataSource interface.
type userDataSource struct {
	client    *client.CloudAvenue
	iamClient *iam.Client
}

func (d *userDataSource) Init(_ context.Context, _ *userDataSourceModel) (diags diag.Diagnostics) {
	var err error

	d.iamClient, err = d.client.CAVSDK.V1.IAM()
	if err != nil {
		diags.AddError("Error initializing IAM client", err.Error())
	}

	return diags
}

// Metadata returns the resource type name.
func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_user"
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = userSchema().GetDataSource(ctx)
}

// Configure configures the data source.
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_iam_user", d.client.GetOrgName(), metrics.Read)()

	config := &userDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		user *iam.UserClient
		err  error
	)

	// Get the user by name or ID and return an error if it doesn't exist or there is another error
	if config.ID.IsKnown() {
		user, err = d.iamClient.GetUser(config.ID.Get())
	} else {
		user, err = d.iamClient.GetUser(config.Name.Get())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.Diagnostics.AddError("User not found", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	config.ID.Set(user.User.ID)
	config.Name.Set(user.User.Name)
	config.RoleName.Set(user.User.RoleName)
	config.FullName.Set(user.User.FullName)
	config.Email.Set(user.User.Email)
	config.Telephone.Set(user.User.Telephone)
	config.Enabled.Set(user.User.Enabled)
	config.ProviderType.Set(string(user.User.Type))
	config.DeployedVMQuota.Set(int64(user.User.DeployedVMQuota))
	config.StoredVMQuota.Set(int64(user.User.StoredVMQuota))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
