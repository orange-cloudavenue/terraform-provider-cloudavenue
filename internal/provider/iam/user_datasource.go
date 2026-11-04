// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	iam "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/iam/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
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

	var v2 cav.Client
	if v2, err = d.client.V2(); err == nil {
		d.iamClient, err = iam.New(v2)
	}
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
		user *iam.ModelUser
		err  error
	)

	// Get the user by name or ID and return an error if it doesn't exist or there is another error
	if config.ID.IsKnown() {
		user, err = d.iamClient.GetUser(ctx, iam.ParamsGetUser{ID: config.ID.Get()})
	} else {
		user, err = d.iamClient.GetUser(ctx, iam.ParamsGetUser{Name: config.Name.Get()})
	}
	if err != nil {
		if cerrs.IsNotFound(err) {
			resp.Diagnostics.AddError("User not found", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	config.ID.Set(user.ID)
	config.Name.Set(user.Name)
	config.RoleName.Set(user.RoleName)
	config.FullName.Set(user.FullName)
	config.Email.Set(user.EmailAddress)
	config.Telephone.Set(user.Telephone)
	config.Enabled.Set(user.IsEnabled)
	config.ProviderType.Set(mapLegacyProviderType(user.ProviderType))

	// CloudAPI never returns quotas on the user resource; read them back
	// from the legacy admin user record (the XML API is deprecated but is
	// still the only surface that echoes back the exact assigned values).
	// TODO(entityQuotas): switch to the CloudAPI quotas read once validated
	// against quota policies — see getUserQuotas.
	deployedQuota, storedQuota, err := getUserQuotas(ctx, d.iamClient, user.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user quotas", err.Error())
		return
	}
	config.DeployedVMQuota.Set(int64(deployedQuota))
	config.StoredVMQuota.Set(int64(storedQuota))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
