package main

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	CAProvider "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider"
	VCDProvider "github.com/vmware/terraform-provider-vcd/v3/vcd"
)

var (
	numberCAResources = 1
	vcdEquivalentCA   = map[string]string{
		"vcd_catalog_access_control ":    "_catalog_acl",
		"vcd_independent_disk":           "_vm_disk",
		"vcd_inserted_media":             "_vm_inserted_media",
		"vcd_nsxt_network_dhcp_binding ": "_network_dhcp_binding",
		"vcd_network_isolated_v2":        "_network_isolated",
		"vcd_network_routed_v2":          "_network_routed",
		"vcd_nsxt_app_port_profile":      "_edgegateway_app_port_profile",
		"vcd_nsxt_ip_set":                "_edgegateway_ip_set",
		"vcd_nsxt_ipsec_vpn_tunnel":      "_edgegateway_ipsec_vpn_tunnel",
		"vcd_nsxt_nat_rule":              "_edgegateway_nat_rule",
		"vcd_nsxt_security_group":        "_edgegateway_security_group",
		"vcd_vapp_network":               "_vapp_isolated_network",
		"vcd_vapp_vm":                    "_vm",
		"vcd_vapp_access_control":        "_vapp_acl",
		"vcd_vm_internal_disk":           "_vm_disk",
		"vcd_org_group ":                 "_iam_group",
		"vcd_org_user":                   "_iam_user",
		"vcd_role":                       "_iam_role",
		"vcd_security_tag":               "_vm_security_tag",
	}

	vcdNotApplicableCA = []string{
		"vcd_catalog_item",
		"vcd_certificate_library",
		"vcd_edgegateway",
		"vcd_edgegateway_settings",
		"vcd_external_network",
		"vcd_external_network_v2",
		"vcd_global_role",
		"vcd_network_direct",
		"vcd_network_isolated",
		"vcd_network_routed",
		"vcd_nsxt_alb_cloud",
		"vcd_nsxt_alb_controller",
		"vcd_nsxt_alb_edgegateway_service_engine_group",
		"vcd_nsxt_alb_pool",
		"vcd_nsxt_alb_service_engine_group",
		"vcd_nsxt_alb_settings",
		"vcd_nsxt_alb_virtual_service",
		"vcd_nsxt_edgegateway",
		"vcd_nsxt_edgegateway_bgp_configuration",
		"vcd_nsxt_edgegateway_bgp_ip_prefix_list",
		"vcd_nsxt_edgegateway_bgp_neighbor",
		"vcd_nsxt_network_imported",
		"vcd_nsxt_route_advertisement",
		"vcd_nsxv_dhcp_relay",
		"vcd_nsxv_dnat",
		"vcd_nsxv_firewall_rule",
		"vcd_nsxv_ip_set",
		"vcd_nsxv_snat",
		"vcd_nsxv_distributed_firewall",
		"vcd_org",
		"vcd_org_ldap",
		"vcd_org_vdc",
		"vcd_rights_bundle",
		"vcd_subscribed_catalog",
		"vcd_vdc_group",
		"vcd_vm",
		"vcd_vm_placement_policy",
		"vcd_vm_sizing_policy",
	}
)

//nolint:all
func main() {
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	// Init provider Cloud Avenue
	ppCA := CAProvider.New(CAProvider.VCDVersion)
	// Init provider VMware Cloud Director
	ppVCD := VCDProvider.Provider()

	// Print Resources List in Orange Cloud Avenue Provider
	blue.Printf("Checking resources and datasources\n")
	fmt.Printf("=====================================\n\n")
	CAtfSchemaR := exportCAResources(ppCA)
	blue.Printf("* Found %v resources in terraform in Orange Cloud Avenue provider\n", len(CAtfSchemaR))
	// Print DataSources List in Orange Cloud Avenue Provider
	CAtfSchemaD := exportCADataSources(ppCA)
	blue.Printf("* Found %v datasources in terraform Orange Cloud Avenue provider\n", len(CAtfSchemaD))

	// Print Resources List in VMware Cloud Director Provider
	blue.Printf("\nChecking resources and datasources\n")
	fmt.Printf("=====================================\n\n")

	// Sort Resources List in VMware Cloud Director Provider
	// VCDtfSchemaR := ppVCD.ResourcesMap
	VCDtfSchemaR := sortVCDResources(ppVCD.ResourcesMap)
	blue.Printf("* Found %v resources in terraform in VMware Cloud Director provider\n", len(VCDtfSchemaR))

	VCDtfSchemaD := ppVCD.DataSourcesMap
	blue.Printf("* Found %v datasources in terraform VMware Cloud Director provider\n\n", len(VCDtfSchemaD))

Search:
	for k, v := range VCDtfSchemaR {
		// Print if the Resource Name in VMWARE Cloud Provider is applicable for Orange Cloud Avenue Provider
		for _, j := range vcdNotApplicableCA {
			if k == j {
				// red.Printf("* Found %v resources ==> Not Applicable\n", k)
				continue Search
			}

		}
		// Print if Resource is deprecated in VMware Cloud Provider
		if v.DeprecationMessage == "" {
			blue.Printf("* Found %v resources", k)
		} else {
			red.Printf("* Found %v resources deprecated", k)
		}

		// Print if the Resource is implemented in Orange Cloud Avenue Provider
		for _, v := range CAtfSchemaD {
			if k == "vcd"+v {
				green.Printf(" ==> * Found resources implemented: %v (%v)", "cloudavenue"+v, numberCAResources)
				numberCAResources++
			}
		}
		// Print if the Resource is renamed in Orange Cloud Avenue Provider
		for i, j := range vcdEquivalentCA {
			if i == k {
				x := "cloudavenue" + j
				// if renamed, find the name in Orange Cloud Avenue Provider
				if findCAResourceName(ppCA, x) {
					green.Printf(" ==> * Found rename resources implemented: %v (%v)", "cloudavenue"+j, numberCAResources)
					numberCAResources++
				} else {
					yellow.Printf(" ==> * Found ERROR: Not yet implemented: %v", "cloudavenue"+j)
				}
			}
		}
		green.Printf("\n")
	}

	// blue.Printf("* Found %v datasources in terraform VMware Cloud Director provider\n\n", len(VCDtfSchemaD))
	// for k, v := range VCDtfSchemaD {
	// 	if v.DeprecationMessage == "" {
	// 		blue.Printf("* Found %v datasources implemented\n", k)
	// 	} else {
	// 		red.Printf("* Found %v datasources deprecated\n", k)
	// 	}
	// }

}

// Export Resources List.
func exportCAResources(pp func() provider.Provider) []string {
	var export []string
	rResp := &resource.MetadataResponse{}

	// Export Resource List
	for _, i := range pp().Resources(nil) { //nolint: staticcheck
		i().Metadata(nil, resource.MetadataRequest{}, rResp) //nolint: staticcheck
		export = append(export, rResp.TypeName)
	}
	return export
}

// Export DataSources List.
func exportCADataSources(pp func() provider.Provider) []string {
	var export []string
	dResp := &datasource.MetadataResponse{}

	// Export DataSource List
	for _, i := range pp().DataSources(nil) { //nolint: staticcheck
		i().Metadata(nil, datasource.MetadataRequest{}, dResp) //nolint: staticcheck
		export = append(export, dResp.TypeName)
	}
	return export
}

// Find Resource Name in Orange Cloud Avenue Provider.
func findCAResourceName(pp func() provider.Provider, name string) bool {
	rResp := &resource.MetadataResponse{}
	for _, i := range pp().Resources(nil) { //nolint: staticcheck
		i().Metadata(nil, resource.MetadataRequest{}, rResp) //nolint: staticcheck
		if "cloudavenue"+rResp.TypeName == name {
			return true
		}
	}
	return false
}

// Sort Resources List.
func sortVCDResources(m map[string]*schema.Resource) map[string]*schema.Resource {
	sortedKeys := make([]string, 0, len(m))
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	// return sliceSorted
	ms := make(map[string]*schema.Resource)
	for _, k := range sortedKeys {
		ms[k] = m[k]
	}
	return ms
}
