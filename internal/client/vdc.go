package client

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// * VDC

type VDC struct {
	name string
	*govcd.Vdc
}

// GetName give you the name of the vDC.
func (v VDC) GetName() string {
	return v.Vdc.Vdc.Name
}

// GetID give you the ID of the vDC.
func (v VDC) GetID() string {
	return v.Vdc.Vdc.ID
}

// IsVDCGroup return true if the object is a VDC Group.
func (v VDC) IsVDCGroup() bool {
	return govcd.OwnerIsVdcGroup(v.GetID())
}

// GetSecurityGroupByID return the NSX-T security group using the ID provided in the argument.
func (v VDC) GetSecurityGroupByID(nsxtFirewallGroupID string) (*govcd.NsxtFirewallGroup, error) {
	return v.Vdc.GetNsxtFirewallGroupById(nsxtFirewallGroupID)
}

// GetSecurityGroupByName return the NSX-T security group using the name provided in the argument.
func (v VDC) GetSecurityGroupByName(nsxtFirewallGroupName string) (*govcd.NsxtFirewallGroup, error) {
	return v.Vdc.GetNsxtFirewallGroupByName(nsxtFirewallGroupName, govcdtypes.FirewallGroupTypeSecurityGroup)
}

// GetSecurityGroupByNameOrID return the NSX-T security group using the name or ID provided in the argument.
func (v VDC) GetSecurityGroupByNameOrID(nsxtFirewallGroupNameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if uuid.IsValid(nsxtFirewallGroupNameOrID) {
		return v.GetSecurityGroupByID(nsxtFirewallGroupNameOrID)
	}

	return v.GetSecurityGroupByName(nsxtFirewallGroupNameOrID)
}

// GetDefaultPlacementPolicyID give you the ID of the default placement policy.
func (v VDC) GetDefaultPlacementPolicyID() string {
	return v.Vdc.Vdc.DefaultComputePolicy.ID
}

// GetVAPP give you the vApp using the name provided in the argument.
func (v VDC) GetVAPP(nameOrID string, refresh bool) (*VAPP, error) {
	vapp, err := v.GetVAppByNameOrId(nameOrID, refresh)
	if err != nil {
		return nil, err
	}

	return &VAPP{vapp}, nil
}

// CreateVAPP create new vApp.
func (v VDC) CreateVAPP(name, description string) (*VAPP, error) {
	vapp, err := v.CreateRawVApp(name, description)
	if err != nil {
		return nil, err
	}

	return &VAPP{vapp}, nil
}

// * VDCGroup

type VDCGroup struct {
	*govcd.VdcGroup
}

// GetName give you the name of the vDC Group.
func (g VDCGroup) GetName() string {
	return g.VdcGroup.VdcGroup.Name
}

// GetID give you the ID of the vDC Group.
func (g VDCGroup) GetID() string {
	return g.VdcGroup.VdcGroup.Id
}

// IsVDCGroup return true if the object is a VDC Group.
func (g VDCGroup) IsVDCGroup() bool {
	return govcd.OwnerIsVdcGroup(g.GetID())
}

// GetSecurityGroupByID return the NSX-T security group using the ID provided in the argument.
func (g VDCGroup) GetSecurityGroupByID(nsxtFirewallGroupID string) (*govcd.NsxtFirewallGroup, error) {
	return g.VdcGroup.GetNsxtFirewallGroupById(nsxtFirewallGroupID)
}

// GetSecurityGroupByName return the NSX-T security group using the name provided in the argument.
func (g VDCGroup) GetSecurityGroupByName(nsxtFirewallGroupName string) (*govcd.NsxtFirewallGroup, error) {
	return g.VdcGroup.GetNsxtFirewallGroupByName(nsxtFirewallGroupName, govcdtypes.FirewallGroupTypeSecurityGroup)
}

// GetSecurityGroupByNameOrID return the NSX-T security group using the name or ID provided in the argument.
func (g VDCGroup) GetSecurityGroupByNameOrID(nsxtFirewallGroupNameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if uuid.IsValid(nsxtFirewallGroupNameOrID) {
		return g.GetSecurityGroupByID(nsxtFirewallGroupNameOrID)
	}

	return g.GetSecurityGroupByName(nsxtFirewallGroupNameOrID)
}

// VDCOrVDCGroupHandler is an interface to access some common methods on VDC or VDC Group without
// explicitly handling exact types.
type VDCOrVDCGroupHandler interface {
	// * Global Get
	GetName() string
	GetID() string

	// * Global Is
	IsVDCGroup() bool

	// * Network
	GetOpenApiOrgVdcNetworkByName(string) (*govcd.OpenApiOrgVdcNetwork, error)

	// * Security Group
	GetSecurityGroupByID(nsxtFirewallGroupID string) (*govcd.NsxtFirewallGroup, error)
	GetSecurityGroupByName(nsxtFirewallGroupName string) (*govcd.NsxtFirewallGroup, error)
	GetSecurityGroupByNameOrID(nsxtFirewallGroupNameOrID string) (*govcd.NsxtFirewallGroup, error)
}

type GetVDCOpts func(*VDC)

func WithVDCName(name string) GetVDCOpts {
	return func(vdc *VDC) {
		vdc.name = name
	}
}

// GetVDC
// return the vdc using the name provided in the argument.
// If the name is empty, it will try to use the default vdc provided in the provider.
func (c *CloudAvenue) GetVDC(opts ...GetVDCOpts) (*VDC, error) {
	v := &VDC{}

	for _, opt := range opts {
		opt(v)
	}

	if v.name == "" {
		if c.DefaultVDCExist() {
			v.name = c.GetDefaultVDC()
		} else {
			return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
		}
	}

	org, err := c.GetOrg()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRetrievingOrg, err)
	}

	x, err := org.GetVDCByName(v.name, false)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDC, v.name, err)
	}

	v.Vdc = x

	return v, nil
}

// GetVDCGroup return the vdc group using the name provided in the argument.
func (c *CloudAvenue) GetVDCGroup(vdcGroupName string) (*VDCGroup, error) {
	if vdcGroupName == "" {
		return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
	}

	org, err := c.GetOrg()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRetrievingOrg, err)
	}

	adminOrg, err := c.Vmware.GetAdminOrgByName(org.GetName())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRetrievingOrgAdmin, err)
	}

	x, err := adminOrg.GetVdcGroupByName(vdcGroupName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDCGroup, vdcGroupName, err)
	}

	return &VDCGroup{x}, nil
}

// GetVDCOrVDCGroup return the vdc or vdc group using the name provided in the argument.
func (c *CloudAvenue) GetVDCOrVDCGroup(vdcOrVDCGroupName string) (VDCOrVDCGroupHandler, error) {
	x, err := c.GetVDC(
		WithVDCName(vdcOrVDCGroupName),
	)
	if err == nil {
		return x, nil
	}

	xGroup, err := c.GetVDCGroup(vdcOrVDCGroupName)
	if err == nil {
		return xGroup, nil
	}

	return nil, fmt.Errorf("error retrieving VDC or VDC Group %s: not found", vdcOrVDCGroupName)
}
