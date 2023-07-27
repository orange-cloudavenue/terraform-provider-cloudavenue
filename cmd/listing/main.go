package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"

	vcdProvider "github.com/vmware/terraform-provider-vcd/v3/vcd"

	"github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	caProvider "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider"

	_ "github.com/vmware/terraform-provider-vcd/v3/vcd"
)

var (
	red    = color.New(color.FgRed)
	green  = color.New(color.FgGreen)
	blue   = color.New(color.FgBlue)
	yellow = color.New(color.FgYellow)

	vcdEquivalentCA = map[string]string{
		"vcd_catalog_access_control ":          "cloudavenue_catalog_acl",
		"vcd_independent_disk":                 "cloudavenue_vm_disk",
		"vcd_inserted_media":                   "cloudavenue_vm_inserted_media",
		"vcd_network_isolated_v2":              "cloudavenue_network_isolated",
		"vcd_network_routed_v2":                "cloudavenue_network_routed",
		"vcd_nsxt_alb_pool":                    "cloudavenue_alb_pool",
		"vcd_nsxt_app_port_profile":            "cloudavenue_edgegateway_app_port_profile",
		"vcd_nsxt_edgegateway":                 "cloudavenue_edgegateway",
		"vcd_nsxt_firewall":                    "cloudavenue_edgegateway_firewall",
		"vcd_nsxt_ip_set":                      "cloudavenue_edgegateway_ip_set",
		"vcd_nsxt_ipsec_vpn_tunnel":            "cloudavenue_edgegateway_ipsec_vpn_tunnel",
		"vcd_nsxt_nat_rule":                    "cloudavenue_edgegateway_nat_rule",
		"vcd_nsxt_network_dhcp_binding ":       "cloudavenue_network_dhcp_binding",
		"vcd_nsxt_security_group":              "cloudavenue_edgegateway_security_group",
		"vcd_vapp_network":                     "cloudavenue_vapp_isolated_network",
		"vcd_vapp_vm":                          "cloudavenue_vm",
		"vcd_vapp_access_control":              "cloudavenue_vapp_acl",
		"vcd_vm_internal_disk":                 "cloudavenue_vm_disk",
		"vcd_org_group ":                       "cloudavenue_iam_group",
		"vcd_org_user":                         "cloudavenue_iam_user",
		"vcd_org_vdc":                          "cloudavenue_vdc",
		"vcd_org_vdc_access_control":           "cloudavenue_vdc_acl",
		"vcd_role":                             "cloudavenue_iam_role",
		"vcd_security_tag":                     "cloudavenue_vm_security_tag",
		"vcd_right":                            "cloudavenue_iam_right",
		"vcd_nsxt_network_dhcp":                "cloudavenue_network_dhcp",
		"vcd_nsxt_network_dhcp_binding":        "cloudavenue_network_dhcp_binding",
		"vcd_vm":                               "cloudavenue_vm",
		"vcd_nsxt_edgegateway_dhcp_forwarding": "cloudavenue_edgegateway_dhcp_forwarding",
		"vcd_nsxt_edgegateway_static_route":    "cloudavenue_edgegateway_static_route",
	}

	vcdNotApplicableCA = []string{
		"vcd_catalog_item",
		"vcd_certificate_library",
		"vcd_edgegateway",          // NSX-V not supported on cloudavenue
		"vcd_edgegateway_settings", // NSX-V not supported on cloudavenue
		"vcd_edgegateway_vpn",      // NSX-V not supported on cloudavenue
		"vcd_external_network",     // NSX-V not supported on cloudavenue
		"vcd_external_network_v2",
		"vcd_global_role",
		"vcd_lb_app_profile",     // NSX-V not supported on cloudavenue
		"vcd_lb_app_rule",        // NSX-V not supported on cloudavenue
		"vcd_lb_server_pool",     // NSX-V not supported on cloudavenue
		"vcd_lb_service_monitor", // NSX-V not supported on cloudavenue
		"vcd_lb_virtual_server",  // NSX-V not supported on cloudavenue
		"vcd_network_direct",     // NSX-V not supported on cloudavenue
		"vcd_network_isolated",   // NSX-V not supported on cloudavenue
		"vcd_network_routed",     // NSX-V not supported on cloudavenue
		"vcd_nsxt_alb_cloud",
		"vcd_nsxt_alb_controller",
		"vcd_nsxt_alb_edgegateway_service_engine_group",
		"vcd_nsxt_alb_service_engine_group",
		"vcd_nsxt_alb_settings",
		"vcd_nsxt_alb_virtual_service",
		"vcd_nsxt_edgegateway_bgp_configuration",
		"vcd_nsxt_edgegateway_bgp_ip_prefix_list",
		"vcd_nsxt_edgegateway_bgp_neighbor",
		"vcd_nsxt_network_imported",
		"vcd_nsxt_route_advertisement",
		"vcd_nsxv_dhcp_relay",           // NSX-V not supported on cloudavenue
		"vcd_nsxv_dnat",                 // NSX-V not supported on cloudavenue
		"vcd_nsxv_firewall_rule",        // NSX-V not supported on cloudavenue
		"vcd_nsxv_ip_set",               // NSX-V not supported on cloudavenue
		"vcd_nsxv_snat",                 // NSX-V not supported on cloudavenue
		"vcd_nsxv_distributed_firewall", // NSX-V not supported on cloudavenue
		"vcd_org",
		"vcd_org_ldap",
		"vcd_org_saml",
		"vcd_org_vdc",
		"vcd_provider_vdc",
		"vcd_rights_bundle",
		"vcd_subscribed_catalog",
		"vcd_vdc_group",
		"vcd_vm",
		"vcd_vm_placement_policy",
		"vcd_vm_sizing_policy",
		"vcd_org_group",       // Manage group for LDAP and SAML
		"vcd_resource_schema", // Generic data source.
		"vcd_resource_pool",   // Require Admin Org
		"vcd_resource_list",   // Generic data source.
		"vcd_nsxt_manager",    // NSX-T manager
		"vcd_cloned_vapp",
		"vcd_ip_space",              // IP Space require Admin Org
		"vcd_ip_space_custom_quota", // IP Space require Admin Org
		"vcd_ip_space_uplink",       // IP Space require Admin Org
		"vcd_nsxt_edgegateway_dhcpv6",
		"vcd_rde",
		"vcd_rde_interface",
		"vcd_rde_interface_behavior",
		"vcd_rde_type",
		"vcd_rde_type_behavior",
		"vcd_rde_type_behavior_acl",
		"vcd_ui_plugin",
		"vcd_vm_placement_policy", // Require Admin Org
		"vcd_vm_sizing_policy",    // Require Admin Org
	}
)

//nolint:all
func main() {

	var mess string
	file, err := os.Create("./resource-ca.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Init provider Cloud Avenue
	ppCA := caProvider.New(caProvider.VCDVersion)
	// Init provider VMware Cloud Director
	ppVCD := vcdProvider.Provider()
	ppVCDVersion := vcdProvider.BuildVersion

	// Print Resources List in Orange Cloud Avenue Provider
	mess = "Checking resources and datasources of Orange Cloud Avenue provider\n"
	blue.Printf(mess)
	file.WriteString("# " + mess)
	fmt.Printf("==================================================================\n\n")

	caTFSchemaR := exportCAResources(ppCA)
	mess = fmt.Sprintf("- Found %v resources in terraform\n", len(caTFSchemaR))
	blue.Printf(mess)
	file.WriteString(mess)

	// Print DataSources List in Orange Cloud Avenue Provider
	caTFSchemaD := exportCADataSources(ppCA)
	mess = fmt.Sprintf("- Found %v datasources in terraform\n\n", len(caTFSchemaD))
	blue.Printf(mess)
	file.WriteString(mess)

	// Print Resources List in VMware Cloud Director Provider
	mess = "Checking resources and datasources of VMware Cloud Director provider (version: " + ppVCDVersion + ")\n"
	blue.Printf(mess)
	file.WriteString("# " + mess)
	fmt.Printf("====================================================================\n\n")

	// Sort Resources List in VMware Cloud Director Provider
	vcdTFSchemaR := ppVCD.ResourcesMap
	mess = fmt.Sprintf("- Found %v resources in terraform\n", len(vcdTFSchemaR))
	blue.Printf(mess)
	file.WriteString(mess)

	vcdTFSchemaD := ppVCD.DataSourcesMap
	mess = fmt.Sprintf("- Found %v datasources in terraform\n\n", len(vcdTFSchemaD))
	blue.Printf(mess)
	file.WriteString(mess)

	// Print Resources List from Orange Cloud Avenue Provider
	mess = "\n# Listing cross resources and datasources from Cloud Avenue\n"
	blue.Printf(mess)
	file.WriteString(mess)
	fmt.Printf("=======================================\n")

	findResourcesFromCA(vcdTFSchemaR, caTFSchemaR, file, "Resources")
	findResourcesFromCA(vcdTFSchemaD, caTFSchemaD, file, "Datasources")

	// Print Resources List from VMware Cloud Director Provider
	mess = "\n# Listing cross resources and datasources from VCD (version: " + ppVCDVersion + ")\n"
	blue.Printf(mess)
	file.WriteString(mess)
	fmt.Printf("=======================================\n")

	findResourcesFromVCD(vcdTFSchemaR, caTFSchemaR, ppCA, file, "Resources")
	findResourcesFromVCD(vcdTFSchemaD, caTFSchemaD, ppCA, file, "Datasources")

}

// Find and print Resources from Orange Cloud Avenue Provider.
func findResourcesFromCA(vcdTFSchemaR map[string]*schema.Resource, caTFSchemaR []string, file *os.File, typeR string) {
	numberCAResources := 1
	// Print if the Resource Name in Orange Cloud Avenue Provider is applicable for VMWARE Cloud Provider
	sortCAResources(caTFSchemaR)
	mess := "\n| Number | " + typeR + " Orange Cloud Avenue | " + typeR + " VMware VCD |\n|:--:|:--:|:--:|\n"
	blue.Printf(mess)
	wf(mess, file)

begin:
	for _, v := range caTFSchemaR {
		mess := fmt.Sprintf("| (%v) | cloudavenue%v ", numberCAResources, v)
		blue.Printf(mess)
		wf(mess, file)
		numberCAResources++

		// Print if the Resource is implemented in Orange Cloud Avenue Provider
		for i := range vcdTFSchemaR {
			if i == "vcd"+v {
				mess = fmt.Sprintf("| %v |\n", "vcd"+v)
				green.Printf(mess)
				wf(mess, file)
				continue begin
			}
		}

		// Print if the Resource have a different name in VMware VCD Provider
		for i, j := range vcdEquivalentCA {
			if "cloudavenue"+v == j {
				mess = fmt.Sprintf("| %v |\n", i)
				green.Printf(mess)
				wf(mess, file)
				continue begin
			}
		}

		mess = "|\n"
		blue.Printf(mess)
		wf(mess, file)
	}
}

// Find and print Resources from VMware Cloud Provider.
func findResourcesFromVCD(vcdTFSchemaRUnsort map[string]*schema.Resource, caTFSchemaR []string, ppCA func() provider.Provider, file *os.File, typeR string) {
	numberVCDResources := 0

	// Sort slice of keys
	sortedKeys := make([]string, 0, len(vcdTFSchemaRUnsort))
	for k := range vcdTFSchemaRUnsort {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	// Print if the Resource Name in VMWARE Cloud Provider is applicable for Orange Cloud Avenue Provider
	mess := "\n| Number | " + typeR + " VMware VCD | " + typeR + " Orange Cloud Avenue | status |\n|:--:|:--:|:--:|:--:|\n"
	blue.Printf(mess)
	wf(mess, file)

beginVCD:
	for _, k := range sortedKeys {
		v := vcdTFSchemaRUnsort[k]
		numberVCDResources++
		mess = fmt.Sprintf("| (%v) | %v | ", numberVCDResources, k)
		blue.Printf(mess)
		wf(mess, file)
		for _, j := range vcdNotApplicableCA {
			if k == j {
				mess = "Not Applicable | :heavy_multiplication_x: |\n"
				red.Printf(mess)
				wf(mess, file)
				continue beginVCD
			}
		}
		// Print if Resource is deprecated in VMware Cloud Provider
		if v.DeprecationMessage != "" {
			mess = "Deprecated | :warning: |\n"
			red.Printf("mess")
			wf(mess, file)
			continue beginVCD
		}

		// Print if the Resource is implemented in Orange Cloud Avenue Provider
		for _, v := range caTFSchemaR {
			if k == "vcd"+v {
				mess = fmt.Sprintf("%v |:white_check_mark: |\n", "cloudavenue"+v)
				green.Printf(mess)
				wf(mess, file)
				continue beginVCD
			}
		}
		// Print if the Resource is renamed in Orange Cloud Avenue Provider
		for i, j := range vcdEquivalentCA {
			if i == k {
				x := j
				// if renamed, find the name in Orange Cloud Avenue Provider
				if findCAResourceName(ppCA, x) {
					mess = fmt.Sprintf("%v | :white_check_mark: |\n", j)
					green.Printf(mess)
					wf(mess, file)
					continue beginVCD
				}
			}
		}
		mess = "Not yet implemented | :x: |\n"
		yellow.Printf(mess)
		wf(mess, file)
	}
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

// Sort Resources List Cloudavenue.
func sortCAResources(m []string) []string {
	sort.Strings(m)
	return m
}

func wf(mess string, file *os.File) {
	var err error
	_, err = file.WriteString(mess)
	if err != nil {
		panic(err)
	}
}
