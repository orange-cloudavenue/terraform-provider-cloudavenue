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
	SecurityGroup     = VcloudUUID(VcloudUUIDPrefix + "firewallGroup:")

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
	Disk,
	SecurityGroup,
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
	if uuid.isEmpty() || prefix.isEmpty() {
		return false
	}

	return strings.HasPrefix(string(uuid), string(prefix)) && isUUIDV4(uuid.extractUUIDv4(prefix))
}

// isNotEmpty returns true if the UUID is not empty.
func (uuid VcloudUUID) isEmpty() bool {
	return len(uuid) == 0
}

func isUUIDV4(uuid string) bool {
	return regexp.MustCompile(`(?m)^\w{8}-\w{4}-\w{4}-\w{4}-\w{12}$`).MatchString(uuid)
}

// ContainsPrefix returns true if the UUID contains any prefix.
func (uuid VcloudUUID) ContainsPrefix() bool {
	return strings.Contains(string(uuid), string(VcloudUUIDPrefix))
}

// extractUUIDv4 returns the UUIDv4 from the UUID.
func (uuid VcloudUUID) extractUUIDv4(prefix VcloudUUID) string {
	return extractUUIDv4(uuid.String(), prefix)
}

func extractUUIDv4(uuid string, prefix VcloudUUID) string {
	if len(uuid) == 0 || prefix.isEmpty() {
		return ""
	}

	return uuid[len(prefix):]
}

func IsValid(uuid string) bool {
	if len(uuid) == 0 {
		return false
	}

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
	if len(uuid) == 0 || prefix.isEmpty() {
		return ""
	}

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

// IsDisk returns true if the UUID is a Disk UUID.
func (uuid VcloudUUID) IsDisk() bool {
	return uuid.IsType(Disk)
}

// IsSecurityGroup returns true if the UUID is a SecurityGroup UUID.
func (uuid VcloudUUID) IsSecurityGroup() bool {
	return uuid.IsType(SecurityGroup)
}

// IsEdgeGateway returns true if the UUID is a EdgeGateway UUID.
func IsEdgeGateway(uuid string) bool {
	return VcloudUUID(uuid).IsType(Gateway)
}

// IsVDC returns true if the UUID is a VDC UUID.
func IsVDC(uuid string) bool {
	return VcloudUUID(uuid).IsType(VDC)
}

// IsNetwork returns true if the UUID is a Network UUID.
func IsNetwork(uuid string) bool {
	return VcloudUUID(uuid).IsType(Network)
}

// IsLoadBalancerPool returns true if the UUID is a LoadBalancerPool UUID.
func IsLoadBalancerPool(uuid string) bool {
	return VcloudUUID(uuid).IsType(LoadBalancerPool)
}

// IsVDCStorageProfile returns true if the UUID is a VDCStorageProfile UUID.
func IsVDCStorageProfile(uuid string) bool {
	return VcloudUUID(uuid).IsType(VDCStorageProfile)
}

// IsVAPP returns true if the UUID is a VAPP UUID.
func IsVAPP(uuid string) bool {
	return VcloudUUID(uuid).IsType(VAPP)
}

// IsDisk returns true if the UUID is a Disk UUID.
func IsDisk(uuid string) bool {
	return VcloudUUID(uuid).IsType(Disk)
}

// IsSecurityGroup returns true if the UUID is a SecurityGroup UUID.
func IsSecurityGroup(uuid string) bool {
	return VcloudUUID(uuid).IsType(SecurityGroup)
}

// IsVCDA returns true if the UUID is a VCDA UUID.
func IsVCDA(uuid string) bool {
	return VcloudUUID(uuid).IsType(VCDA)
}

// IsVM returns true if the UUID is a VM UUID.
func IsVM(uuid string) bool {
	return VcloudUUID(uuid).IsType(VM)
}

// IsUser returns true if the UUID is a User UUID.
func IsUser(uuid string) bool {
	return VcloudUUID(uuid).IsType(User)
}

// IsGroup returns true if the UUID is a Group UUID.
func IsGroup(uuid string) bool {
	return VcloudUUID(uuid).IsType(Group)
}
