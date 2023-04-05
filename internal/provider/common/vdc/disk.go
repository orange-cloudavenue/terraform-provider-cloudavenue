package vdc

import "github.com/vmware/go-vcloud-director/v2/govcd"

// DiskExist checks if a disk exists in a VDC.
func (v VDC) DiskExist(diskName string) (bool, error) {
	existingDisk, err := v.QueryDisk(diskName)
	return existingDisk != (govcd.DiskRecord{}), err
}
