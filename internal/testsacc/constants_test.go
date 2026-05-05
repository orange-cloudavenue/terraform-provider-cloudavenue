/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

// Common test name constants used across multiple acceptance test files.
const (
	// testNameExample is the default test name used in most acceptance tests.
	testNameExample = "example"

	// testNameExampleWithVDCGroup is the test name for VDC group tests.
	testNameExampleWithVDCGroup = "example_with_vdc_group"

	// testNameExampleWithID is the test name for tests using resource ID.
	testNameExampleWithID = "example_with_id"

	// Common attribute keys used in ImportStateIDBuilder slices.
	testAttrEdgeGatewayID   = "edge_gateway_id"
	testAttrEdgeGatewayName = "edge_gateway_name"
	testAttrName            = "name"

	// Common check attribute keys.
	testAttrPortsHash  = "ports.#"
	testAttrValuesHash = "values.#"
	testAttrValue      = "value"

	// Protocol constants used in app port profile tests.
	testProtocolTCP    = "TCP"
	testProtocolUDP    = "UDP"
	testProtocolICMPv4 = "ICMPv4"

	// ELB HTTP policy criteria constants.
	testCriteriaContains   = "CONTAINS"
	testCriteriaBeginsWith = "BEGINS_WITH"

	// ELB HTTP policy header name constants.
	testHeaderXExample = "X-EXAMPLE"
	testHeaderXCustom  = "X-CUSTOM"
	testHeaderXSecure  = "X-SECURE"

	// ELB HTTP policy action constants.
	testActionADD    = "ADD"
	testActionREMOVE = "REMOVE"

	// Network IP constants used in network_routed tests.
	testNetworkIP10  = "192.168.40.10"
	testNetworkIP20  = "192.168.40.20"
	testNetworkIP100 = "192.168.40.100"
	testNetworkIP130 = "192.168.40.130"

	// Network IP constants used in vdcg_network_isolated tests.
	testIsolatedNetworkIP10  = "192.168.0.10"
	testIsolatedNetworkIP20  = "192.168.0.20"
	testIsolatedNetworkIP100 = "192.168.0.100"
	testIsolatedNetworkIP130 = "192.168.0.130"

	// Network IP constants used in vdcg_network_routed tests.
	testRoutedNetworkIP10  = "192.168.100.10"
	testRoutedNetworkIP20  = "192.168.100.20"
	testRoutedNetworkIP100 = "192.168.100.100"
	testRoutedNetworkIP130 = "192.168.100.130"

	// ELB pool member IP constant.
	testELBPoolMemberIP = "192.168.0.1"

	// VDC group attribute keys.
	testAttrVDCGroupID   = "vdc_group_id"
	testAttrVDCGroupName = "vdc_group_name"

	// vApp attribute key.
	testAttrVAppName = "vapp_name"

	// Network attribute keys.
	testAttrStartAddress = "start_address"
	testAttrEndAddress   = "end_address"

	// Common boolean string constants.
	testValueTrue  = "true"
	testValueFalse = "false"

	// Common attribute keys used in check functions.
	testAttrEnabled = "enabled"
	testAttrPort    = "port"
	testAttrRatio   = "ratio"
	testAttrBucket  = "bucket"

	// Firewall rule attribute constants.
	testFirewallActionALLOW      = "ALLOW"
	testFirewallActionREJECT     = "REJECT"
	testFirewallDirectionINOUT   = "IN_OUT"
	testFirewallDirectionOUT     = "OUT"
	testFirewallIPProtocolIPV4   = "IPV4"
	testFirewallRuleNameAllowAll = "allow all IPv4 traffic"
	testFirewallRuleNameAllowOut = "allow out IPv4 traffic"
	testFirewallRuleNameRejectIn = "reject in IPv4 traffic"

	// Firewall rule attribute keys.
	testAttrAppPortProfileIDsHash     = "app_port_profile_ids.#"
	testAttrSourceIDsHash             = "source_ids.#"
	testAttrDestinationIDsHash        = "destination_ids.#"
	testAttrSourceGroupsExcluded      = "source_groups_excluded"
	testAttrDestinationGroupsExcluded = "destination_groups_excluded"

	// Dynamic security group criteria constants.
	testCriteriaStartsWith = "STARTS_WITH"

	// Dynamic security group member type constants.
	testMemberTypeVMName = "VM_NAME"
	testMemberTypeVMTag  = "VM_TAG"

	// Dynamic security group test value constants.
	testDSGValueTest     = "test"
	testDSGValueFront    = "front"
	testDSGValueProd     = "prod"
	testDSGValueWebFront = "web-front"

	// S3 website configuration constants.
	testS3RedirectHostname = "redirect.hostname"
	testS3RedirectHTTPCode = "redirect.http_redirect_code"
	testS3SchemeHTTPS      = "https"

	// Additional test name constants.
	testNameExample2 = "example_2"

	// Common attribute keys used in check functions (additional).
	testAttrProtocol         = "protocol"
	testAttrCriteria         = "criteria"
	testAttrAction           = "action"
	testAttrTakeOwnership    = "take_ownership"
	testAttrIPAddress        = "ip_address"
	testAttrTargetName       = "target_name"
	testAttrTargetID         = "target_id"
	testAttrType             = "type"
	testAttrOperator         = "operator"
	testAttrDirection        = "direction"
	testAttrIPProtocol       = "ip_protocol"
	testAttrIPAllocationMode = "ip_allocation_mode"
	testAttrIsPrimary        = "is_primary"
	testAttrOrg              = "org"
	testAttrVDC              = "vdc"
	testAttrClass            = "class"
	testAttrDefault          = "default"
	testAttrLimit            = "limit"

	// Validator type constants.
	testValidatorUUID4 = "uuid4"

	// S3 website configuration additional constants.
	testS3WWWExample       = "www.example.com"
	testS3RedirectProtocol = "redirect.protocol"
)
