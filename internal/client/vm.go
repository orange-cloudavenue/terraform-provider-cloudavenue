package client

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type VM struct {
	name string
	id   string
	*govcd.VM
}

// * Guest properties
// GetGuestProperties returns the guest properties of a VM.
func (v VM) GetGuestProperties() (guestProperties []*govcdtypes.Property, err error) {
	x, err := v.GetProductSectionList()
	if err != nil {
		return nil, err
	}

	guestProperties = append(guestProperties, x.ProductSection.Property...)

	return
}

// SetGuestProperties sets the guest properties of a VM
// If the guest property already exists, it will be updated.
func (v VM) SetGuestProperties(guestProperties map[string]string) (err error) {
	listGuestProperties := make([]*govcdtypes.Property, 0)

	for key, value := range guestProperties {
		listGuestProperties = append(listGuestProperties, &govcdtypes.Property{
			UserConfigurable: true,
			Type:             "string",
			Key:              key,
			Label:            key,
			Value:            &govcdtypes.Value{Value: value},
		})
	}

	_, err = v.SetProductSectionList(&govcdtypes.ProductSectionList{
		ProductSection: &govcdtypes.ProductSection{
			Info:     "Custom properties",
			Property: listGuestProperties,
		},
	})

	return
}

// * Customization

// GetCustomization returns the customization of a VM.
func (v VM) GetCustomization() (guestCustomization *govcdtypes.GuestCustomizationSection, err error) {
	return v.GetGuestCustomizationSection()
}

// SetCustomization sets the customization of a VM.
func (v VM) SetCustomization(guestCustomization *govcdtypes.GuestCustomizationSection) (err error) {
	_, err = v.SetGuestCustomizationSection(guestCustomization)
	return
}

// * OS type

// SetOSType sets the OS type of a VM.
func (v VM) SetOSType(osType string) (err error) {
	updateOsType := v.VM.VM.VmSpecSection

	updateOsType.OsType = osType

	_, err = v.UpdateVmSpecSection(updateOsType, v.VM.VM.Description)
	return
}

// GetOSType returns the OS type of a VM.
func (v VM) GetOSType() string {
	return v.VM.VM.VmSpecSection.OsType
}
