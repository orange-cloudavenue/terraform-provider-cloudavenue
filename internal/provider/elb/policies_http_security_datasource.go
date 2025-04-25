/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package elb provides a Terraform datasource.
package elb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &PoliciesHTTPSecurityDataSource{}
	_ datasource.DataSourceWithConfigure = &PoliciesHTTPSecurityDataSource{}
)

func NewPoliciesHTTPSecurityDataSource() datasource.DataSource {
	return &PoliciesHTTPSecurityDataSource{}
}

type PoliciesHTTPSecurityDataSource struct {
	client *client.CloudAvenue
	elb    edgeloadbalancer.Client
}

// Init Initializes the data source.
func (d *PoliciesHTTPSecurityDataSource) Init(ctx context.Context, dm *PoliciesHTTPSecurityModel) (diags diag.Diagnostics) {
	var err error

	d.elb, err = edgeloadbalancer.NewClient()
	if err != nil {
		diags.AddError("Error creating elb client", err.Error())
	}

	return
}

func (d *PoliciesHTTPSecurityDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_policies_http_security"
}

func (d *PoliciesHTTPSecurityDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = policiesHTTPSecuritySchema(ctx).GetDataSource(ctx)
}

func (d *PoliciesHTTPSecurityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoliciesHTTPSecurityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_elb_policies_http_security", d.client.GetOrgName(), metrics.Read)()

	config := &PoliciesHTTPSecurityModel{}

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

	s := &PoliciesHTTPSecurityResource{
		client: d.client,
		elb:    d.elb,
	}

	// Read data from the API
	data, found, diags := s.read(ctx, config)
	if !found {
		resp.Diagnostics.AddError("Resource not found", "The resource was not found")
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
