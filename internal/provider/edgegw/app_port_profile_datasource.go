// Package edgegw provides a Terraform datasource.
package edgegw

import (
	"context"
	"fmt"
	"net/url"

	"github.com/k0kubun/pp"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &portProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &portProfileDataSource{}
)

func NewPortProfileDataSource() datasource.DataSource {
	return &portProfileDataSource{}
}

type portProfileDataSource struct {
	client *client.CloudAvenue

	// Uncomment the following lines if you need to access the resource's.
	org org.Org

	// vapp   vapp.VAPP
}

// objectType.
func (p *portProfileModelAppPorts) ObjectType(ctx context.Context) types.ObjectType {
	return types.ObjectType{
		AttrTypes: p.AttrTypes(ctx),
	}
}

// attrTypes().
func (p *portProfileModelAppPorts) AttrTypes(_ context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"protocol": types.StringType,
		"ports":    types.SetType{ElemType: types.StringType},
	}
}

// If the data source don't have same schema/structure as the resource, you can use the following code:
// type appPortProfileDataSourceModel struct {
// 	ID types.String `tfsdk:"id"`
// }

// Init Initializes the data source.
func (d *portProfileDataSource) Init(ctx context.Context, dm *portProfileModel) (diags diag.Diagnostics) {
	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return
	}

	// d.vdc, diags = vdc.Init(d.client, dm.VDC)
	// if diags.HasError() {
	// 	return
	// }

	// Uncomment the following lines if you need to access to the VAPP
	// d.vapp, diags = vapp.Init(d.client, d.vdc, dm.VAppID, dm.VAppName)

	return
}

func (d *portProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_app_port_profile"
}

func (d *portProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = portProfilesSchema(ctx).GetDataSource(ctx)
}

func (d *portProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *portProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	plan := &portProfileModel{}

	// Get current plan
	resp.Diagnostics.Append(req.Config.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the data source read logic here.
	*/

	s := &portProfilesResource{
		client: d.client,
		org:    d.org,
	}

	// If error len >1 Got the first
	var (
		queryParams  url.Values
		portProfile  *govcd.NsxtAppPortProfile
		portProfiles []*govcd.NsxtAppPortProfile
		err          error
	)

	// TODO - waiting answer from UGO Vincent about AppPortProfile VDC or ORG propagation
	queryParams = queryParameterFilterAnd("name=="+plan.Name.ValueString(), queryParams)
	// queryParams.Set("pageSize", "1024")
	// TODO - Change resource to Add EdgeGateway
	// TODO - Use UUID pakage for new uuid appPortProfile
	// TODO - Use function GetAppPortProfileByNameOrID in method EdgeGateway (common)
	queryParams = queryParameterFilterAnd("_context==urn:vcloud:vdc:889318ed-e5ea-43a0-a03f-e99d9c06d3e3", queryParams)
	portProfiles, err = s.org.GetAllNsxtAppPortProfiles(queryParams, "")
	for _, v := range portProfiles {
		if v.NsxtAppPortProfile.Name == plan.Name.ValueString() {
			// portProfile, err = s.org.GetNsxtAppPortProfileById(v.NsxtAppPortProfile.ID)
			portProfile, err = s.org.GetNsxtAppPortProfileByName(plan.Name.ValueString(), "")
			tflog.Info(ctx, pp.Sprint(v.NsxtAppPortProfile))
			// tflog.Info(ctx, pp.Sprint(v.NsxtAppPortProfile.Name))
		}
	}

	if err != nil {
		if govcd.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading NSX-T App Port Profile",
			fmt.Sprintf("Error reading NSX-T App Port Profile: %s", err),
		)
		return
	}

	appPortsState, dia := s.AppPortRead(ctx, portProfile)
	resp.Diagnostics.Append(dia...)
	if resp.Diagnostics.HasError() {
		return
	}

	data := &portProfileModel{
		ID:          types.StringValue(portProfile.NsxtAppPortProfile.ID),
		Name:        types.StringValue(portProfile.NsxtAppPortProfile.Name),
		Description: utils.StringValueOrNull(portProfile.NsxtAppPortProfile.Description),
		VDC:         plan.VDC,
	}
	data.AppPorts, dia = appPortsState.ToPlan(ctx)
	resp.Diagnostics.Append(dia...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func queryParameterFilterAnd(filter string, parameters url.Values) url.Values {
	newParameters := copyOrNewUrlValues(parameters)

	existingFilter := newParameters.Get("filter")
	if existingFilter == "" {
		newParameters.Set("filter", filter)
		return newParameters
	}

	newParameters.Set("filter", existingFilter+";"+filter)
	return newParameters
}

func copyOrNewUrlValues(parameters url.Values) url.Values {
	parameterCopy := make(map[string][]string)

	// if supplied parameters are nil - we just return new initialized
	if parameters == nil {
		return parameterCopy
	}

	// Copy URL values
	for key, value := range parameters {
		parameterCopy[key] = value
	}

	return parameterCopy
}
