package vm

import (
	"sort"

	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type osType struct {
	name        string
	description string
}

type osAllTypes struct {
	windows osTypes
	other   osTypes
	linux   osTypes
}

type osTypes map[string]osType

func getOsTypeLinux() osTypes {
	return utils.SortMapStringByKeys(map[string]osType{
		"sles15_64Guest":       {"sles15_64Guest", "SUSE Linux Enterprise Server 15 (64-bit)"},
		"rhel8_64Guest":        {"rhel8_64Guest", "Red Hat Enterprise Linux 8 (64-bit)"},
		"other4xLinux64Guest":  {"other4xLinux64Guest", "Other Linux 4.x (64-bit)"},
		"other4xLinuxGuest":    {"other4xLinuxGuest", "Other Linux 4.x (32-bit)"},
		"oracleLinux8_64Guest": {"oracleLinux8_64Guest", "Oracle Linux 8 (64-bit)"},
		"centos8_64Guest":      {"centos8_64Guest", "CentOS Linux 8 (64-bit)"},
		"asianux8_64Guest":     {"asianux8_64Guest", "ASIANUX 8 (64-bit)"},
		"amazonlinux2_64Guest": {"amazonlinux2_64Guest", "Amazon Linux 2 (64-bit)"},
		"vmwarePhoton64Guest":  {"vmwarePhoton64Guest", "VMware Photon OS 64-bit"},
		"oracleLinux7_64Guest": {"oracleLinux7_64Guest", "Oracle Linux 7 (64-bit)"},
		"oracleLinux6_64Guest": {"oracleLinux6_64Guest", "Oracle Linux 6 (64-bit)"},
		"oracleLinux6Guest":    {"oracleLinux6Guest", "Oracle Linux 6 (32-bit)"},
		"debian9_64Guest":      {"debian9_64Guest", "Debian Linux 9 (64-bit)"},
		"debian9Guest":         {"debian9Guest", "Debian Linux 9 (32-bit)"},
		"debian10_64Guest":     {"debian10_64Guest", "Debian Linux 10 (64-bit)"},
		"debian10Guest":        {"debian10Guest", "Debian Linux 10 (32-bit)"},
		"centos7_64Guest":      {"centos7_64Guest", "CentOS Linux 7 (64-bit)"},
		"centos6_64Guest":      {"centos6_64Guest", "CentOS Linux 6 (64-bit)"},
		"centos6Guest":         {"centos6Guest", "CentOS Linux 6 (32-bit)"},
		"asianux7_64Guest":     {"asianux7_64Guest", "ASIANUX 7 (64-bit)"},
		"debian8_64Guest":      {"debian8_64Guest", "Debian Linux 8 (64-bit)"},
		"debian8Guest":         {"debian8Guest", "Debian Linux 8 (32-bit)"},
		"coreos64Guest":        {"coreos64Guest", "CoreOS Linux (64-bit)"},
		"other3xLinux64Guest":  {"other3xLinux64Guest", "Other Linux 3.x (64-bit)"},
		"other3xLinuxGuest":    {"other3xLinuxGuest", "Other Linux 3.x (32-bit)"},
		"debian7_64Guest":      {"debian7_64Guest", "Debian Linux 7 (64-bit)"},
		"debian7Guest":         {"debian7Guest", "Debian Linux 7 (32-bit)"},
		"sles12_64Guest":       {"sles12_64Guest", "SUSE Linux Enterprise Server 12 (64-bit)"},
		"rhel7_64Guest":        {"rhel7_64Guest", "Red Hat Enterprise Linux 7 (64-bit)"},
		"asianux4_64Guest":     {"asianux4_64Guest", "ASIANUX 4 (64-bit)"},
		"asianux4Guest":        {"asianux4Guest", "ASIANUX 4 (32-bit)"},
		"rhel6_64Guest":        {"rhel6_64Guest", "Red Hat Enterprise Linux 6 (64-bit)"},
		"rhel6Guest":           {"rhel6Guest", "Red Hat Enterprise Linux 6 (32-bit)"},
		"other26xLinux64Guest": {"other26xLinux64Guest", "Other Linux 2.6.x (64-bit)"},
		"other26xLinuxGuest":   {"other26xLinuxGuest", "Other Linux 2.6.x (32-bit)"},
		"other24xLinux64Guest": {"other24xLinux64Guest", "Other Linux 2.4.x (64-bit)"},
		"other24xLinuxGuest":   {"other24xLinuxGuest", "Other Linux 2.4.x (32-bit)"},
		"oracleLinux64Guest":   {"oracleLinux64Guest", "Oracle Linux 5 (64-bit)"},
		"oracleLinuxGuest":     {"oracleLinuxGuest", "Oracle Linux 5 (32-bit)"},
		"debian6_64Guest":      {"debian6_64Guest", "Debian Linux 6 (64-bit)"},
		"debian6Guest":         {"debian6Guest", "Debian Linux 6 (32-bit)"},
		"debian5_64Guest":      {"debian5_64Guest", "Debian Linux 5 (64-bit)"},
		"debian5Guest":         {"debian5Guest", "Debian Linux 5 (32-bit)"},
		"debian4_64Guest":      {"debian4_64Guest", "Debian Linux 4 (64-bit)"},
		"debian4Guest":         {"debian4Guest", "Debian Linux 4 (32-bit)"},
		"centos64Guest":        {"centos64Guest", "CentOS Linux 5 (64-bit)"},
		"centosGuest":          {"centosGuest", "CentOS Linux 5 (32-bit)"},
		"asianux3_64Guest":     {"asianux3_64Guest", "ASIANUX 3 (64-bit)"},
		"asianux3Guest":        {"asianux3Guest", "ASIANUX 3 (32-bit)"},
		"ubuntu64Guest":        {"ubuntu64Guest", "Ubuntu Linux 64-bit"},
		"ubuntuGuest":          {"ubuntuGuest", "Ubuntu Linux 32-bit"},
		"sles64Guest":          {"sles64Guest", "SUSE Linux Enterprise Server 11 (64-bit)"},
		"slesGuest":            {"slesGuest", "SUSE Linux Enterprise Server 11 (32-bit)"},
		"sles11_64Guest":       {"sles11_64Guest", "SUSE Linux Enterprise Server 11 (64-bit)"},
		"sles11Guest":          {"sles11Guest", "SUSE Linux Enterprise Server 11 (32-bit)"},
		"sles10_64Guest":       {"sles10_64Guest", "SUSE Linux Enterprise Server 10 (64-bit)"},
		"sles10Guest":          {"sles10Guest", "SUSE Linux Enterprise Server 10 (32-bit)"},
		"rhel5_64Guest":        {"rhel5_64Guest", "Red Hat Enterprise Linux 5 (64-bit)"},
		"rhel5Guest":           {"rhel5Guest", "Red Hat Enterprise Linux 5 (32-bit)"},
		"rhel4_64Guest":        {"rhel4_64Guest", "Red Hat Enterprise Linux 4 (64-bit)"},
		"rhel4Guest":           {"rhel4Guest", "Red Hat Enterprise Linux 4 (32-bit)"},
		"rhel3_64Guest":        {"rhel3_64Guest", "Red Hat Enterprise Linux 3 (64-bit)"},
		"rhel3Guest":           {"rhel3Guest", "Red Hat Enterprise Linux 3 (32-bit)"},
		"rhel2Guest":           {"rhel2Guest", "Red Hat Enterprise Linux 2 (32-bit)"},
		"otherLinux64Guest":    {"otherLinux64Guest", "Other Linux (64-bit)"},
		"otherLinuxGuest":      {"otherLinuxGuest", "Other Linux (32-bit)"},
		"oesGuest":             {"oesGuest", "Novell Open Enterprise Server (32-bit)"},
	})
}

// //nolint:dupl
// getOsTypeWindows returns the osType for the given name.
func getOsTypeWindows() osTypes {
	return utils.SortMapStringByKeys(map[string]osType{
		"windows9_64Guest":        {"windows9_64Guest", "Microsoft Windows 10 (64-bit)"},
		"windows9Guest":           {"windows9Guest", "Microsoft Windows 10 (32-bit)"},
		"windows8Server64Guest":   {"windows8Server64Guest", "Microsoft Windows Server 2012 (64-bit)"},
		"windows8_64Guest":        {"windows8_64Guest", "Microsoft Windows 8.x (64-bit)"},
		"windows8Guest":           {"windows8Guest", "Microsoft Windows 8.x (32-bit)"},
		"win98Guest":              {"win98Guest", "Microsoft Windows 98"},
		"win95Guest":              {"win95Guest", "Microsoft Windows 95"},
		"win31Guest":              {"win31Guest", "Microsoft Windows 3.1"},
		"dosGuest":                {"dosGuest", "Microsoft MS-DOS"},
		"winXPPro64Guest":         {"winXPPro64Guest", "Microsoft Windows XP Professional (64-bit)"},
		"winXPProGuest":           {"winXPProGuest", "Microsoft Windows XP Professional (32-bit)"},
		"winVista64Guest":         {"winVista64Guest", "Microsoft Windows Vista (64-bit)"},
		"winVistaGuest":           {"winVistaGuest", "Microsoft Windows Vista (32-bit)"},
		"windows7Server64Guest":   {"windows7Server64Guest", "Microsoft Windows Server 2008 R2 (64-bit)"},
		"winLonghorn64Guest":      {"winLonghorn64Guest", "Microsoft Windows Server 2008 (64-bit)"},
		"winLonghornGuest":        {"winLonghornGuest", "Microsoft Windows Server 2008 (32-bit)"},
		"winNetWebGuest":          {"winNetWebGuest", "Microsoft Windows Server 2003 Web Edition (32-bit)"},
		"winNetStandard64Guest":   {"winNetStandard64Guest", "Microsoft Windows Server 2003 Standard Edition (64-bit)"},
		"winNetStandardGuest":     {"winNetStandardGuest", "Microsoft Windows Server 2003 Standard Edition (32-bit)"},
		"winNetDatacenter64Guest": {"winNetDatacenter64Guest", "Microsoft Windows Server 2003 Datacenter Edition (64-bit)"},
		"winNetDatacenterGuest":   {"winNetDatacenterGuest", "Microsoft Windows Server 2003 Datacenter Edition (32-bit)"},
		"winNetEnterprise64Guest": {"winNetEnterprise64Guest", "Microsoft Windows Server 2003 (64-bit)"},
		"winNetEnterpriseGuest":   {"winNetEnterpriseGuest", "Microsoft Windows Server 2003 (32-bit)"},
		"winNTGuest":              {"winNTGuest", "Microsoft Windows NT"},
		"windows7_64Guest":        {"windows7_64Guest", "Microsoft Windows 7 (64-bit)"},
		"windows7Guest":           {"windows7Guest", "Microsoft Windows 7 (32-bit)"},
		"win2000ServGuest":        {"win2000ServGuest", "Microsoft Windows 2000 Server"},
		"win2000ProGuest":         {"win2000ProGuest", "Microsoft Windows 2000 Professional"},
		"win2000AdvServGuest":     {"win2000AdvServGuest", "Microsoft Windows 2000"},
		"winNetBusinessGuest":     {"winNetBusinessGuest", "Microsoft Windows Small Business Server 2003"},
	})
}

// //nolint:dupl
// getOsTypeOther returns the osType for the given name.
func getOsTypeOther() osTypes {
	// allowVMReboot := true
	// if allowVMReboot {
	// 	return map[string]osType{}
	// }
	return utils.SortMapStringByKeys(map[string]osType{
		"freebsd12_64Guest": {"freebsd12_64Guest", "FreeBSD 12 or later versions (64-bit)"},
		"freebsd12Guest":    {"freebsd12Guest", "FreeBSD 12 or later versions (32-bit)"},
		"freebsd11_64Guest": {"freebsd11_64Guest", "FreeBSD 11 (64-bit)"},
		"freebsd11Guest":    {"freebsd11Guest", "FreeBSD 11 (32-bit)"},
		"darwin18_64Guest":  {"darwin18_64Guest", "Apple macOS 10.14 (64-bit)"},
		"darwin17_64Guest":  {"darwin17_64Guest", "Apple macOS 10.13 (64-bit)"},
		"darwin16_64Guest":  {"darwin16_64Guest", "Apple macOS 10.12 (64-bit)"},
		"darwin15_64Guest":  {"darwin15_64Guest", "Apple macOS 10.11 (64-bit)"},
		"darwin14_64Guest":  {"darwin14_64Guest", "Apple macOS 10.10 (64-bit)"},
		"darwin13_64Guest":  {"darwin13_64Guest", "Apple macOS 10.9 (64-bit)"},
		"darwin12_64Guest":  {"darwin12_64Guest", "Apple macOS 10.8 (64-bit)"},
		"eComStation2Guest": {"eComStation2Guest", "Serenity Systems eComStation 2"},
		"openServer6Guest":  {"openServer6Guest", "SCO OpenServer 6"},
		"solaris11_64Guest": {"solaris11_64Guest", "Oracle Solaris 11 (64-bit)"},
		"darwin11_64Guest":  {"darwin11_64Guest", "Apple macOS 10.7 (64-bit)"},
		"darwin11Guest":     {"darwin11Guest", "Apple macOS 10.7 (32-bit)"},
		"eComStationGuest":  {"eComStationGuest", "Serenity Systems eComStation 1"},
		"unixWare7Guest":    {"unixWare7Guest", "SCO UnixWare 7"},
		"openServer5Guest":  {"openServer5Guest", "SCO OpenServer 5"},
		"os2Guest":          {"os2Guest", "IBM OS/2"},
		"freebsd64Guest":    {"freebsd64Guest", "FreeBSD Pre-11 versions (64-bit)"},
		"freebsdGuest":      {"freebsdGuest", "FreeBSD Pre-11 versions (32-bit)"},
		"darwin10_64Guest":  {"darwin10_64Guest", "Apple macOS 10.6 (64-bit)"},
		"darwin10Guest":     {"darwin10Guest", "Apple macOS 10.6 (32-bit)"},
		"otherGuest64":      {"otherGuest64", "Other (64-bit)"},
		"otherGuest":        {"otherGuest", "Other (32-bit)"},
		"solaris10_64Guest": {"solaris10_64Guest", "Oracle Solaris 10 (64-bit)"},
		"solaris10Guest":    {"solaris10Guest", "Oracle Solaris 10 (32-bit)"},
		"netware6Guest":     {"netware6Guest", "Novell NetWare 6.x"},
		"netware5Guest":     {"netware5Guest", "Novell NetWare 5.x"},
	})
}

// sortByName sorts the list of osTypes by name.
// func (o *osTypes) sortByName() {
// 	x := []string{}
// 	for k := range *o {
// 		x = append(x, k)
// 	}
//
// 	sort.Strings(x)
//
// 	y := map[string]osType{}
// 	for _, k := range x {
// 		y[k] = (*o)[k]
// 	}
//
// 	*o = y
// }

// getAllOsTypes returns all osTypes name.
func GetAllOsTypes() []string {
	all := osAllTypes{
		linux:   getOsTypeLinux(),
		windows: getOsTypeWindows(),
		other:   getOsTypeOther(),
	}

	x := []string{}
	for k := range all.linux {
		x = append(x, k)
	}

	for k := range all.windows {
		x = append(x, k)
	}

	for k := range all.other {
		x = append(x, k)
	}

	fstringvalidator.OneOfWithDescription()

	sort.Strings(x)

	return x
}

func GetAllOsTypesWithDescription() []fstringvalidator.OneOfWithDescriptionValues {
	all := osAllTypes{
		linux:   getOsTypeLinux(),
		windows: getOsTypeWindows(),
		other:   getOsTypeOther(),
	}

	x := []fstringvalidator.OneOfWithDescriptionValues{}
	for k, v := range all.linux {
		x = append(x, fstringvalidator.OneOfWithDescriptionValues{
			Value:       k,
			Description: v.description,
		})
	}

	for k, v := range all.windows {
		x = append(x, fstringvalidator.OneOfWithDescriptionValues{
			Value:       k,
			Description: v.description,
		})
	}

	for k, v := range all.other {
		x = append(x, fstringvalidator.OneOfWithDescriptionValues{
			Value:       k,
			Description: v.description,
		})
	}

	sort.Slice(x, func(i, j int) bool {
		return x[i].Value < x[j].Value
	})

	return x
}
