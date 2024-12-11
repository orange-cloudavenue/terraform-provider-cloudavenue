package alb

import (
	"context"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
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
func (rm *VirtualServiceModel) toALBVirtualService(ctx context.Context, r *VirtualServiceResource) (albConfig *v1.EdgeGatewayALBVirtualService, diags diag.Diagnostics) {
	// Create the resource ALB Configuration
	albConfig = &v1.EdgeGatewayALBVirtualService{
		VirtualService: &v1.EdgeGatewayALBVirtualServiceModel{
			Name:               rm.Name.Get(),
			Description:        rm.Description.Get(),
			Enabled:            rm.Enabled.GetPtr(),
			ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{Type: rm.ServiceType.Get()},
			VirtualIPAddress:   rm.VirtualIP.Get(),
		},
	}

	// Add Pool if provided
	if rm.PoolID.IsUnknown() {
		// Find the pool by name and use ID
		x, err := r.client.Vmware.GetAlbPoolByName(r.edgegw.GetID(), rm.PoolName.Get())
		if err != nil {
			diags.AddError("Error while fetching Pool", err.Error())
			return nil, diags
		}
		rm.PoolID.Set(x.NsxtAlbPool.ID)
	} else {
		// Use the provided pool ID
		x, err := r.client.Vmware.GetAlbPoolById(rm.PoolID.Get())
		if err != nil {
			diags.AddError("Error while fetching Pool", err.Error())
			return nil, diags
		}
		rm.PoolName.Set(x.NsxtAlbPool.Name)
	}
	albConfig.VirtualService.LoadBalancerPoolRef = govcdtypes.OpenApiReference{ID: rm.PoolID.Get(), Name: rm.PoolName.Get()}

	// Add Service Engine Group if provided
	if rm.ServiceEngineGroupName.IsKnown() {
		albConfig.VirtualService.ServiceEngineGroupRef = govcdtypes.OpenApiReference{Name: rm.ServiceEngineGroupName.Get()}
	}

	// Add certificate if provided
	if rm.CertificateID.IsKnown() {
		albConfig.VirtualService.CertificateRef = &govcdtypes.OpenApiReference{ID: rm.CertificateID.Get()}
	}

	// Add service ports
	sp, d := rm.ServicePorts.Get(ctx)
	if d.HasError() {
		diags = append(diags, d...)
		return nil, diags
	}
	for _, svcPort := range sp {
		if svcPort.PortEnd.IsNull() {
			svcPort.PortEnd = svcPort.PortStart
		}
		albConfig.VirtualService.ServicePorts = append(albConfig.VirtualService.ServicePorts, govcdtypes.NsxtAlbVirtualServicePort{
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
