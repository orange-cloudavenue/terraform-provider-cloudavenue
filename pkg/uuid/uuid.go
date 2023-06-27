package uuid

import (
	"regexp"
	"strings"
)

const (

	// VcloudUUIDPrefix is the prefix for all vCloud UUIDs.
	VcloudUUIDPrefix      = "urn:vcloud:"
	CloudAvenueUUIDPrefix = "urn:cloudavenue:"

	// * VCD.
	VM                = VcloudUUID(VcloudUUIDPrefix + "vm:")
	User              = VcloudUUID(VcloudUUIDPrefix + "user:")
	Group             = VcloudUUID(VcloudUUIDPrefix + "group:")
	Gateway           = VcloudUUID(VcloudUUIDPrefix + "gateway:")
	VDC               = VcloudUUID(VcloudUUIDPrefix + "vdc:")
	Network           = VcloudUUID(VcloudUUIDPrefix + "network:")
	LoadBalancerPool  = VcloudUUID(VcloudUUIDPrefix + "loadBalancerPool:")
	VDCStorageProfile = VcloudUUID(VcloudUUIDPrefix + "vdcstorageProfile:")
	VAPP              = VcloudUUID(VcloudUUIDPrefix + "vapp:")
	Disk              = VcloudUUID(VcloudUUIDPrefix + "disk:")

	// * CLOUDAVENUE.
	VCDA = VcloudUUID(CloudAvenueUUIDPrefix + "vcda:")
)

var vcloudUUIDs = []VcloudUUID{
	VM,
	User,
	Group,
	Gateway,
	VDC,
	Network,
	LoadBalancerPool,
	VDCStorageProfile,
	VAPP,
}

type (
	VcloudUUID string
)

// String returns the string representation of the UUID.
func (uuid VcloudUUID) String() string {
	return string(uuid)
}

// IsType returns true if the UUID is of the specified type.
func (uuid VcloudUUID) IsType(prefix VcloudUUID) bool {
	// remove prefix
	uuidv4 := uuid[len(prefix):]
	return strings.HasPrefix(string(uuid), string(prefix)) && isUUIDV4(string(uuidv4))
}

func isUUIDV4(uuid string) bool {
	return regexp.MustCompile(`(?m)^\w{8}-\w{4}-\w{4}-\w{4}-\w{12}$`).MatchString(uuid)
}

// ContainsPrefix returns true if the UUID contains any prefix.
func (uuid VcloudUUID) ContainsPrefix() bool {
	return strings.Contains(string(uuid), string(VcloudUUIDPrefix))
}

func extractUUIDv4(uuid string, prefix VcloudUUID) string {
	return uuid[len(prefix):]
}

func IsValid(uuid string) bool {
	u := VcloudUUID(uuid)

	for _, prefix := range vcloudUUIDs {
		if u.IsType(prefix) {
			return isUUIDV4(extractUUIDv4(uuid, prefix))
		}
	}
	return false
}

// Normalize returns the UUID with the prefix if prefix is missing.
func Normalize(prefix VcloudUUID, uuid string) VcloudUUID {
	u := VcloudUUID(uuid)

	if u.ContainsPrefix() {
		return u
	}

	return prefix + u
}

// IsVM returns true if the UUID is a VM UUID.
func (uuid VcloudUUID) IsVM() bool {
	return uuid.IsType(VM)
}

// IsUser returns true if the UUID is a User UUID.
func (uuid VcloudUUID) IsUser() bool {
	return uuid.IsType(User)
}

// IsGroup returns true if the UUID is a Group UUID.
func (uuid VcloudUUID) IsGroup() bool {
	return uuid.IsType(Group)
}

// IsGateway returns true if the UUID is a Gateway UUID.
func (uuid VcloudUUID) IsGateway() bool {
	return uuid.IsType(Gateway)
}

// IsVDC returns true if the UUID is a VDC UUID.
func (uuid VcloudUUID) IsVDC() bool {
	return uuid.IsType(VDC)
}

// IsNetwork returns true if the UUID is a Network UUID.
func (uuid VcloudUUID) IsNetwork() bool {
	return uuid.IsType(Network)
}

// IsLoadBalancerPool returns true if the UUID is a LoadBalancerPool UUID.
func (uuid VcloudUUID) IsLoadBalancerPool() bool {
	return uuid.IsType(LoadBalancerPool)
}

// IsVDCStorageProfile returns true if the UUID is a VDCStorageProfile UUID.
func (uuid VcloudUUID) IsVDCStorageProfile() bool {
	return uuid.IsType(VDCStorageProfile)
}

// IsVAPP returns true if the UUID is a VAPP UUID.
func (uuid VcloudUUID) IsVAPP() bool {
	return uuid.IsType(VAPP)
}
