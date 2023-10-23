package uuid

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	VDCGroup          = VcloudUUID(VcloudUUIDPrefix + "vdcGroup:")
	VDCComputePolicy  = VcloudUUID(VcloudUUIDPrefix + "vdcComputePolicy:")
	Network           = VcloudUUID(VcloudUUIDPrefix + "network:")
	LoadBalancerPool  = VcloudUUID(VcloudUUIDPrefix + "loadBalancerPool:")
	VDCStorageProfile = VcloudUUID(VcloudUUIDPrefix + "vdcstorageProfile:")
	VAPP              = VcloudUUID(VcloudUUIDPrefix + "vapp:")
	VAPPTemplate      = VcloudUUID(VcloudUUIDPrefix + "vappTemplate:")
	Disk              = VcloudUUID(VcloudUUIDPrefix + "disk:")
	SecurityGroup     = VcloudUUID(VcloudUUIDPrefix + "firewallGroup:")
	Catalog           = VcloudUUID(VcloudUUIDPrefix + "catalog:")
	Token             = VcloudUUID(VcloudUUIDPrefix + "token:")

	// * CLOUDAVENUE.
	VCDA = VcloudUUID(CloudAvenueUUIDPrefix + "vcda:")
)

var vcloudUUIDs = []VcloudUUID{
	VM,
	User,
	Group,
	Gateway,
	VDC,
	VDCGroup,
	VDCComputePolicy,
	Network,
	LoadBalancerPool,
	VDCStorageProfile,
	VAPP,
	VAPPTemplate,
	Disk,
	SecurityGroup,
	Catalog,
	Token,
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

	return strings.HasPrefix(string(uuid), prefix.String()) && isUUIDV4(uuid.extractUUIDv4(prefix))
}

// isNotEmpty returns true if the UUID is not empty.
func (uuid VcloudUUID) isEmpty() bool {
	return len(uuid) == 0
}

func isUUIDV4(uuid string) bool {
	return regexp.MustCompile(`(?m)^\w{8}-\w{4}-\w{4}-\w{4}-\w{12}$`).MatchString(uuid)
}

func IsUUIDV4(uuid string) bool {
	return isUUIDV4(uuid)
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

// IsVDCGroup returns true if the UUID is a VDCGroup UUID.
func (uuid VcloudUUID) IsVDCGroup() bool {
	return uuid.IsType(VDCGroup)
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

// IsVAPPTemplate returns true if the UUID is a VAPPTemplate UUID.
func (uuid VcloudUUID) IsVAPPTemplate() bool {
	return uuid.IsType(VAPPTemplate)
}

// IsDisk returns true if the UUID is a Disk UUID.
func (uuid VcloudUUID) IsDisk() bool {
	return uuid.IsType(Disk)
}

// IsSecurityGroup returns true if the UUID is a SecurityGroup UUID.
func (uuid VcloudUUID) IsSecurityGroup() bool {
	return uuid.IsType(SecurityGroup)
}

// IsCatalog returns true if the UUID is a Catalog UUID.
func (uuid VcloudUUID) IsCatalog() bool {
	return uuid.IsType(Catalog)
}

// IsToken returns true if the UUID is a Token UUID.
func (uuid VcloudUUID) IsToken() bool {
	return uuid.IsType(Token)
}

// IsVDCComputePolicy returns true if the UUID is a VDCComputePolicy UUID.
func (uuid VcloudUUID) IsVDCComputePolicy() bool {
	return uuid.IsType(VDCComputePolicy)
}

// * End Methods

// IsEdgeGateway returns true if the UUID is a EdgeGateway UUID.
func IsEdgeGateway(uuid string) bool {
	return VcloudUUID(uuid).IsType(Gateway)
}

// IsVDC returns true if the UUID is a VDC UUID.
func IsVDC(uuid string) bool {
	return VcloudUUID(uuid).IsType(VDC)
}

// IsVDCGroup returns true if the UUID is a VDCGroup UUID.
func IsVDCGroup(uuid string) bool {
	return VcloudUUID(uuid).IsType(VDCGroup)
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

// IsVAPPTemplate returns true if the UUID is a VAPPTemplate UUID.
func IsVAPPTemplate(uuid string) bool {
	return VcloudUUID(uuid).IsType(VAPPTemplate)
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

// IsCatalog returns true if the UUID is a Catalog UUID.
func IsCatalog(uuid string) bool {
	return VcloudUUID(uuid).IsType(Catalog)
}

// IsToken returns true if the UUID is a Token UUID.
func IsToken(uuid string) bool {
	return VcloudUUID(uuid).IsType(Token)
}

// IsVDCComputePolicy returns true if the UUID is a VDCComputePolicy UUID.
func IsVDCComputePolicy(uuid string) bool {
	return VcloudUUID(uuid).IsType(VDCComputePolicy)
}

// * End Functions

// TestIsType returns true if the UUID is of the specified type.
func TestIsType(uuidType VcloudUUID) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		if value == "" {
			return nil
		}

		ok := VcloudUUID(value).IsType(uuidType)
		if !ok {
			return fmt.Errorf("uuid %s is not of type %s", value, uuidType)
		}
		return nil
	}
}
