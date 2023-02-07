package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &jobsDataSource{}
	_ datasource.DataSourceWithConfigure = &jobsDataSource{}
)

func NewJobsDataSource() datasource.DataSource {
	return &jobsDataSource{}
}

type jobsDataSource struct {
	client *CloudAvenueClient
}

type jobsDataSourceModel struct {
	ID          types.String  `tfsdk:"id"`
	Actions     []actionModel `tfsdk:"actions"`
	Name        types.String  `tfsdk:"name"`
	Description types.String  `tfsdk:"description"`
	Status      types.String  `tfsdk:"status"`
}

type actionModel struct {
	Name    types.String `tfsdk:"name"`
	Details types.String `tfsdk:"details"`
	Status  types.String `tfsdk:"status"`
}

func (d *jobsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobs"
}

func (d *jobsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The jobs data source show the details of the jobs.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the job.",
				Required:    true,
			},
			"actions": schema.ListNestedAttribute{
				MarkdownDescription: "The actions of the job.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the action.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "The status of the action.",
							Computed:            true,
						},
						"details": schema.StringAttribute{
							MarkdownDescription: "The details of the action.",
							Computed:            true,
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the job.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the job.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the job.",
				Computed:    true,
			},
		},
	}
}

func (d *jobsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*CloudAvenueClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *CloudAvenueClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

func (d *jobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data jobsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ParentJobMetadata, _, err := d.client.JobsApi.ApiCustomersV10JobsJobIdGet(d.client.auth, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Job detail, got error: %s", err))
		return
	}
	if len(ParentJobMetadata) == 0 || len(ParentJobMetadata) > 1 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Job detail, got error: %s", err))
		return
	}

	// Read data from the API into the model
	tab := ParentJobMetadata[0]
	data.Name = types.StringValue(tab.Name)
	data.Status = types.StringValue(tab.Status)
	data.Description = types.StringValue(tab.Description)
	data.Actions = make([]actionModel, len(tab.Actions))
	for j, action := range tab.Actions {
		data.Actions[j].Name = types.StringValue(action.Name)
		data.Actions[j].Status = types.StringValue(action.Status)
		data.Actions[j].Details = types.StringValue(action.Details)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Write the model back to the Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
