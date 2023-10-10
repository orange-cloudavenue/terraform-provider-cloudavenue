// Package backup provides a Terraform datasource.
package backup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &backupDataSource{}
	_ datasource.DataSourceWithConfigure = &backupDataSource{}
)

func NewBackupDataSource() datasource.DataSource {
	return &backupDataSource{}
}

type backupDataSource struct {
	client *client.CloudAvenue
}

func (d *backupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *backupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = backupSchema(ctx).GetDataSource(ctx)
}

func (d *backupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *backupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("cloudavenue_backup", d.client.GetOrgName(), metrics.Read)()

	config := &backupModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	s := &backupResource{
		client: d.client,
	}

	// Refresh data NetBackup from the API
	job, err := d.client.BackupClient.V1.Netbackup.Inventory.Refresh()
	if err != nil {
		resp.Diagnostics.AddError("Error refreshing NetBackup inventory", err.Error())
		return
	}
	if err := job.Wait(1, 45); err != nil {
		resp.Diagnostics.AddError("Error waiting for NetBackup inventory refresh", err.Error())
		return
	}

	// Read data from the API
	data, _, diags := s.read(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
