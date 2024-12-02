package alb

import (
	"context"
	"fmt"
	"net/url"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	VirtualServiceModel struct {
		ID                     supertypes.StringValue                                             `tfsdk:"id"`
		Name                   supertypes.StringValue                                             `tfsdk:"name"`
		EdgeGatewayName        supertypes.StringValue                                             `tfsdk:"edge_gateway_name"`
		EdgeGatewayID          supertypes.StringValue                                             `tfsdk:"edge_gateway_id"`
		Description            supertypes.StringValue                                             `tfsdk:"description"`
		Enabled                supertypes.BoolValue                                               `tfsdk:"enabled"`
		PoolName               supertypes.StringValue                                             `tfsdk:"pool_name"`
		PoolID                 supertypes.StringValue                                             `tfsdk:"pool_id"`
		ServiceEngineGroupName supertypes.StringValue                                             `tfsdk:"service_engine_group_name"`
		VirtualIP              supertypes.StringValue                                             `tfsdk:"virtual_ip"`
		ServiceType            supertypes.StringValue                                             `tfsdk:"service_type"`
		CertificateID          supertypes.StringValue                                             `tfsdk:"certificate_id"`
		ServicePorts           supertypes.ListNestedObjectValueOf[VirtualServiceModelServicePort] `tfsdk:"service_ports"`
	}

	VirtualServiceModelServicePort struct {
		PortStart supertypes.Int64Value  `tfsdk:"port_start"`
		PortEnd   supertypes.Int64Value  `tfsdk:"port_end"`
		PortType  supertypes.StringValue `tfsdk:"port_type"`
		PortSSL   supertypes.BoolValue   `tfsdk:"port_ssl"`
	}
)

func (rm *VirtualServiceModel) Copy() *VirtualServiceModel {
	x := &VirtualServiceModel{}
	utils.ModelCopy(rm, x)
	return x
}

// toNSXTAlbVirtualService converts VirtualServiceModel to NSXTAlbVirtualService.
func (rm *VirtualServiceModel) toNSXTALBVirtualService(ctx context.Context, r *VirtualServiceResource) (albConfig *govcdtypes.NsxtAlbVirtualService, diags diag.Diagnostics) {
	// Create the resource ALB Configuration
	albConfig = &govcdtypes.NsxtAlbVirtualService{
		Name:                rm.Name.Get(),
		Description:         rm.Description.Get(),
		ApplicationProfile:  govcdtypes.NsxtAlbVirtualServiceApplicationProfile{Type: rm.ServiceType.Get()},
		Enabled:             rm.Enabled.GetPtr(),
		VirtualIpAddress:    rm.VirtualIP.Get(),
		GatewayRef:          govcdtypes.OpenApiReference{ID: r.edgegw.GetID(), Name: r.edgegw.GetName()},
		LoadBalancerPoolRef: govcdtypes.OpenApiReference{ID: rm.PoolID.Get(), Name: rm.PoolName.Get()},
	}

	// Add Service Engine Group if provided
	if rm.ServiceEngineGroupName.IsKnown() {
		albConfig.ServiceEngineGroupRef = govcdtypes.OpenApiReference{Name: rm.ServiceEngineGroupName.Get()}
	} else {
		// Find the first service engine group
		queryParams := url.Values{}
		queryParams.Add("filter", fmt.Sprintf("gatewayRef.id==%s", r.edgegw.GetID()))
		x, err := r.client.Vmware.GetAllAlbServiceEngineGroupAssignments(queryParams)
		if err != nil {
			diags.AddError("Error while fetching Service Engine Group", err.Error())
			return nil, diags
		}
		if len(x) != 1 {
			diags.AddError("No Service Engine Group found for Edge Gateway %s", r.edgegw.GetName())
			return nil, diags
		}
		albConfig.ServiceEngineGroupRef = govcdtypes.OpenApiReference{Name: x[0].NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name, ID: x[0].NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.ID}
	}

	// Add certificate if provided
	if rm.CertificateID.IsKnown() {
		albConfig.CertificateRef = &govcdtypes.OpenApiReference{ID: rm.CertificateID.Get()}
	}

	// Add service ports
	sp, d := rm.ServicePorts.Get(ctx)
	if d.HasError() {
		diags = append(diags, d...)
		return nil, diags
	}
	for _, svcPort := range sp {
		if svcPort.PortEnd.IsNull() || svcPort.PortEnd.IsUnknown() {
			svcPort.PortEnd = svcPort.PortStart
		}
		albConfig.ServicePorts = append(albConfig.ServicePorts, govcdtypes.NsxtAlbVirtualServicePort{
			PortStart:  svcPort.PortStart.GetIntPtr(),
			PortEnd:    svcPort.PortEnd.GetIntPtr(),
			SslEnabled: svcPort.PortSSL.GetPtr(),
			TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
				SystemDefined: true,
				Type:          svcPort.PortType.Get(),
			},
		})
	}

	return albConfig, diags
}
