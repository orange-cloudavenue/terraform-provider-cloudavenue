// Package client is the main client for the CloudAvenue provider.
package client

import (
	"errors"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

var (
	ErrEmptyOrgNameProvided = errors.New("empty Org name provided")
	ErrEmptyVDCNameProvided = errors.New("empty VDC name provided")
	ErrRetrievingOrg        = errors.New("error retrieving Org")
	ErrRetrievingOrgAdmin   = errors.New("error retrieving Org admin")
	ErrEmptyOrgFound        = errors.New("empty Org found")
	ErrRetrievingVDC        = errors.New("error retrieving VDC")
	ErrRetrievingAdminVDC   = errors.New("error retrieving AdminVDC")
	ErrRetrievingVDCGroup   = errors.New("error retrieving VDC Group")
	ErrEmptyVDCFound        = errors.New("empty VDC found")
)

func (c *CloudAvenue) getTemplate(iD string) (vAppTemplate *govcd.VAppTemplate, err error) {
	return c.Vmware.GetVAppTemplateById(iD)
}

// GetTemplate retrieves a vApp template by name or ID.
func (c *CloudAvenue) GetTemplate(iD string) (vAppTemplate *govcd.VAppTemplate, err error) {
	vAppTemplate = govcd.NewVAppTemplate(&c.Vmware.Client)
	template, err := c.getTemplate(iD)
	if err != nil || len(template.VAppTemplate.Children.VM) == 0 {
		return nil, fmt.Errorf("error retrieving vApp template %s: %w", iD, err)
	}

	vAppTemplate.VAppTemplate = template.VAppTemplate.Children.VM[0]
	return
}

// GetTemplateWithVMName retrieves a vApp template with a VM name.
func (c *CloudAvenue) GetTemplateWithVMName(iD, vmName string) (vAppTemplate *govcd.VAppTemplate, err error) {
	vAppTemplate = govcd.NewVAppTemplate(&c.Vmware.Client)
	template, err := c.getTemplate(iD)
	if err != nil {
		return nil, fmt.Errorf("error retrieving vApp template %s: %w", iD, err)
	}

	for i, vm := range template.VAppTemplate.Children.VM {
		if vm.Name == vmName {
			vAppTemplate.VAppTemplate = template.VAppTemplate.Children.VM[i]
			return
		}
	}

	return nil, fmt.Errorf("error retrieving vApp template %s: %w", iD, err)
}

// getAffinityRule retrieves an affinity rule by name.
func (c *CloudAvenue) GetAffinityRule(affinityRuleID string) (affinityRule *govcd.VdcComputePolicyV2, err error) {
	return c.Vmware.GetVdcComputePolicyV2ById(affinityRuleID)
}

// GetBootImage retrieves a boot image by ID.
func (c *CloudAvenue) GetBootImage(bootImageID string) (bootImage *govcdtypes.Media, err error) {
	bi, err := c.Vmware.QueryMediaById(bootImageID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving boot image %s: %w", bootImageID, err)
	}

	return &govcdtypes.Media{HREF: bi.MediaRecord.HREF}, nil
}
