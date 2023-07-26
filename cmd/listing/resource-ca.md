# Checking resources and datasources of Orange Cloud Avenue provider
- Found 28 resources in terraform
- Found 33 datasources in terraform

# Checking resources and datasources of VMware Cloud Director provider (version: unset)
- Found 98 resources in terraform
- Found 105 datasources in terraform


# Listing cross resources and datasources from Cloud Avenue

| Number | Resources Orange Cloud Avenue | Resources VMware VCD |
|:--:|:--:|:--:|
| (1) | cloudavenue_alb_pool | vcd_nsxt_alb_pool |
| (2) | cloudavenue_catalog | vcd_catalog |
| (3) | cloudavenue_edgegateway | vcd_edgegateway |
| (4) | cloudavenue_edgegateway_app_port_profile | vcd_nsxt_app_port_profile |
| (5) | cloudavenue_edgegateway_dhcp_forwarding | vcd_nsxt_edgegateway_dhcp_forwarding |
| (6) | cloudavenue_edgegateway_firewall | vcd_nsxt_firewall |
| (7) | cloudavenue_edgegateway_ip_set | vcd_nsxt_ip_set |
| (8) | cloudavenue_edgegateway_security_group | vcd_nsxt_security_group |
| (9) | cloudavenue_edgegateway_static_route | vcd_nsxt_edgegateway_static_route |
| (10) | cloudavenue_iam_role | vcd_role |
| (11) | cloudavenue_iam_user | vcd_org_user |
| (12) | cloudavenue_network_dhcp | vcd_nsxt_network_dhcp |
| (13) | cloudavenue_network_dhcp_binding | vcd_nsxt_network_dhcp_binding |
| (14) | cloudavenue_network_isolated | vcd_network_isolated |
| (15) | cloudavenue_network_routed | vcd_network_routed |
| (16) | cloudavenue_publicip |
| (17) | cloudavenue_vapp | vcd_vapp |
| (18) | cloudavenue_vapp_acl | vcd_vapp_access_control |
| (19) | cloudavenue_vapp_isolated_network | vcd_vapp_network |
| (20) | cloudavenue_vapp_org_network | vcd_vapp_org_network |
| (21) | cloudavenue_vcda_ip |
| (22) | cloudavenue_vdc | vcd_org_vdc |
| (23) | cloudavenue_vdc_acl | vcd_org_vdc_access_control |
| (24) | cloudavenue_vm | vcd_vm |
| (25) | cloudavenue_vm_affinity_rule | vcd_vm_affinity_rule |
| (26) | cloudavenue_vm_disk | vcd_independent_disk |
| (27) | cloudavenue_vm_inserted_media | vcd_inserted_media |
| (28) | cloudavenue_vm_security_tag | vcd_security_tag |

| Number | Datasources Orange Cloud Avenue | Datasources VMware VCD |
|:--:|:--:|:--:|
| (1) | cloudavenue_alb_pool | vcd_nsxt_alb_pool |
| (2) | cloudavenue_catalog | vcd_catalog |
| (3) | cloudavenue_catalog_media | vcd_catalog_media |
| (4) | cloudavenue_catalog_medias |
| (5) | cloudavenue_catalog_vapp_template | vcd_catalog_vapp_template |
| (6) | cloudavenue_catalogs |
| (7) | cloudavenue_edgegateway | vcd_edgegateway |
| (8) | cloudavenue_edgegateway_dhcp_forwarding | vcd_nsxt_edgegateway_dhcp_forwarding |
| (9) | cloudavenue_edgegateway_firewall | vcd_nsxt_firewall |
| (10) | cloudavenue_edgegateway_ip_set | vcd_nsxt_ip_set |
| (11) | cloudavenue_edgegateway_security_group | vcd_nsxt_security_group |
| (12) | cloudavenue_edgegateway_static_route | vcd_nsxt_edgegateway_static_route |
| (13) | cloudavenue_edgegateways |
| (14) | cloudavenue_iam_right | vcd_right |
| (15) | cloudavenue_iam_role | vcd_role |
| (16) | cloudavenue_iam_user | vcd_org_user |
| (17) | cloudavenue_network_dhcp | vcd_nsxt_network_dhcp |
| (18) | cloudavenue_network_dhcp_binding | vcd_nsxt_network_dhcp_binding |
| (19) | cloudavenue_network_isolated | vcd_network_isolated |
| (20) | cloudavenue_network_routed | vcd_network_routed |
| (21) | cloudavenue_publicips |
| (22) | cloudavenue_storage_profile | vcd_storage_profile |
| (23) | cloudavenue_storage_profiles |
| (24) | cloudavenue_tier0_vrf |
| (25) | cloudavenue_tier0_vrfs |
| (26) | cloudavenue_vapp | vcd_vapp |
| (27) | cloudavenue_vapp_isolated_network | vcd_vapp_network |
| (28) | cloudavenue_vapp_org_network | vcd_vapp_org_network |
| (29) | cloudavenue_vdc | vcd_org_vdc |
| (30) | cloudavenue_vdc_group | vcd_vdc_group |
| (31) | cloudavenue_vdcs |
| (32) | cloudavenue_vm | vcd_vm |
| (33) | cloudavenue_vm_affinity_rule | vcd_vm_affinity_rule |

# Listing cross resources and datasources from VCD (version: unset)

| Number | Resources VMware VCD | Resources Orange Cloud Avenue | status |
|:--:|:--:|:--:|:--:|
| (1) | vcd_api_token | Not yet implemented | :x: |
| (2) | vcd_catalog | cloudavenue_catalog |:white_check_mark: |
| (3) | vcd_catalog_access_control | Not yet implemented | :x: |
| (4) | vcd_catalog_item | Not Applicable | :heavy_multiplication_x: |
| (5) | vcd_catalog_media | Not yet implemented | :x: |
| (6) | vcd_catalog_vapp_template | Not yet implemented | :x: |
| (7) | vcd_cloned_vapp | Not Applicable | :heavy_multiplication_x: |
| (8) | vcd_edgegateway | Not Applicable | :heavy_multiplication_x: |
| (9) | vcd_edgegateway_settings | Not Applicable | :heavy_multiplication_x: |
| (10) | vcd_edgegateway_vpn | Not Applicable | :heavy_multiplication_x: |
| (11) | vcd_external_network | Not Applicable | :heavy_multiplication_x: |
| (12) | vcd_external_network_v2 | Not Applicable | :heavy_multiplication_x: |
| (13) | vcd_global_role | Not Applicable | :heavy_multiplication_x: |
| (14) | vcd_independent_disk | cloudavenue_vm_disk | :white_check_mark: |
| (15) | vcd_inserted_media | cloudavenue_vm_inserted_media | :white_check_mark: |
| (16) | vcd_ip_space | Not Applicable | :heavy_multiplication_x: |
| (17) | vcd_ip_space_custom_quota | Not Applicable | :heavy_multiplication_x: |
| (18) | vcd_ip_space_ip_allocation | Not yet implemented | :x: |
| (19) | vcd_ip_space_uplink | Not Applicable | :heavy_multiplication_x: |
| (20) | vcd_lb_app_profile | Not Applicable | :heavy_multiplication_x: |
| (21) | vcd_lb_app_rule | Not Applicable | :heavy_multiplication_x: |
| (22) | vcd_lb_server_pool | Not Applicable | :heavy_multiplication_x: |
| (23) | vcd_lb_service_monitor | Not Applicable | :heavy_multiplication_x: |
| (24) | vcd_lb_virtual_server | Not Applicable | :heavy_multiplication_x: |
| (25) | vcd_library_certificate | Not yet implemented | :x: |
| (26) | vcd_network_direct | Not Applicable | :heavy_multiplication_x: |
| (27) | vcd_network_isolated | Not Applicable | :heavy_multiplication_x: |
| (28) | vcd_network_isolated_v2 | cloudavenue_network_isolated | :white_check_mark: |
| (29) | vcd_network_routed | Not Applicable | :heavy_multiplication_x: |
| (30) | vcd_network_routed_v2 | cloudavenue_network_routed | :white_check_mark: |
| (31) | vcd_nsxt_alb_cloud | Not Applicable | :heavy_multiplication_x: |
| (32) | vcd_nsxt_alb_controller | Not Applicable | :heavy_multiplication_x: |
| (33) | vcd_nsxt_alb_edgegateway_service_engine_group | Not Applicable | :heavy_multiplication_x: |
| (34) | vcd_nsxt_alb_pool | cloudavenue_alb_pool | :white_check_mark: |
| (35) | vcd_nsxt_alb_service_engine_group | Not Applicable | :heavy_multiplication_x: |
| (36) | vcd_nsxt_alb_settings | Not Applicable | :heavy_multiplication_x: |
| (37) | vcd_nsxt_alb_virtual_service | Not Applicable | :heavy_multiplication_x: |
| (38) | vcd_nsxt_app_port_profile | cloudavenue_edgegateway_app_port_profile | :white_check_mark: |
| (39) | vcd_nsxt_distributed_firewall | Not yet implemented | :x: |
| (40) | vcd_nsxt_distributed_firewall_rule | Not yet implemented | :x: |
| (41) | vcd_nsxt_dynamic_security_group | Not yet implemented | :x: |
| (42) | vcd_nsxt_edgegateway | cloudavenue_edgegateway | :white_check_mark: |
| (43) | vcd_nsxt_edgegateway_bgp_configuration | Not Applicable | :heavy_multiplication_x: |
| (44) | vcd_nsxt_edgegateway_bgp_ip_prefix_list | Not Applicable | :heavy_multiplication_x: |
| (45) | vcd_nsxt_edgegateway_bgp_neighbor | Not Applicable | :heavy_multiplication_x: |
| (46) | vcd_nsxt_edgegateway_dhcp_forwarding | cloudavenue_edgegateway_dhcp_forwarding | :white_check_mark: |
| (47) | vcd_nsxt_edgegateway_dhcpv6 | Not Applicable | :heavy_multiplication_x: |
| (48) | vcd_nsxt_edgegateway_rate_limiting | Not yet implemented | :x: |
| (49) | vcd_nsxt_edgegateway_static_route | cloudavenue_edgegateway_static_route | :white_check_mark: |
| (50) | vcd_nsxt_firewall | cloudavenue_edgegateway_firewall | :white_check_mark: |
| (51) | vcd_nsxt_ip_set | cloudavenue_edgegateway_ip_set | :white_check_mark: |
| (52) | vcd_nsxt_ipsec_vpn_tunnel | Not yet implemented | :x: |
| (53) | vcd_nsxt_nat_rule | Not yet implemented | :x: |
| (54) | vcd_nsxt_network_dhcp | cloudavenue_network_dhcp | :white_check_mark: |
| (55) | vcd_nsxt_network_dhcp_binding | cloudavenue_network_dhcp_binding | :white_check_mark: |
| (56) | vcd_nsxt_network_imported | Not Applicable | :heavy_multiplication_x: |
| (57) | vcd_nsxt_route_advertisement | Not Applicable | :heavy_multiplication_x: |
| (58) | vcd_nsxt_security_group | cloudavenue_edgegateway_security_group | :white_check_mark: |
| (59) | vcd_nsxv_dhcp_relay | Not Applicable | :heavy_multiplication_x: |
| (60) | vcd_nsxv_distributed_firewall | Not Applicable | :heavy_multiplication_x: |
| (61) | vcd_nsxv_dnat | Not Applicable | :heavy_multiplication_x: |
| (62) | vcd_nsxv_firewall_rule | Not Applicable | :heavy_multiplication_x: |
| (63) | vcd_nsxv_ip_set | Not Applicable | :heavy_multiplication_x: |
| (64) | vcd_nsxv_snat | Not Applicable | :heavy_multiplication_x: |
| (65) | vcd_org | Not Applicable | :heavy_multiplication_x: |
| (66) | vcd_org_group | Not Applicable | :heavy_multiplication_x: |
| (67) | vcd_org_ldap | Not Applicable | :heavy_multiplication_x: |
| (68) | vcd_org_saml | Not Applicable | :heavy_multiplication_x: |
| (69) | vcd_org_user | cloudavenue_iam_user | :white_check_mark: |
| (70) | vcd_org_vdc | Not Applicable | :heavy_multiplication_x: |
| (71) | vcd_org_vdc_access_control | cloudavenue_vdc_acl | :white_check_mark: |
| (72) | vcd_provider_vdc | Not Applicable | :heavy_multiplication_x: |
| (73) | vcd_rde | Not Applicable | :heavy_multiplication_x: |
| (74) | vcd_rde_interface | Not Applicable | :heavy_multiplication_x: |
| (75) | vcd_rde_interface_behavior | Not Applicable | :heavy_multiplication_x: |
| (76) | vcd_rde_type | Not Applicable | :heavy_multiplication_x: |
| (77) | vcd_rde_type_behavior | Not Applicable | :heavy_multiplication_x: |
| (78) | vcd_rde_type_behavior_acl | Not Applicable | :heavy_multiplication_x: |
| (79) | vcd_rights_bundle | Not Applicable | :heavy_multiplication_x: |
| (80) | vcd_role | cloudavenue_iam_role | :white_check_mark: |
| (81) | vcd_security_tag | cloudavenue_vm_security_tag | :white_check_mark: |
| (82) | vcd_service_account | Not yet implemented | :x: |
| (83) | vcd_subscribed_catalog | Not Applicable | :heavy_multiplication_x: |
| (84) | vcd_ui_plugin | Not Applicable | :heavy_multiplication_x: |
| (85) | vcd_vapp | cloudavenue_vapp |:white_check_mark: |
| (86) | vcd_vapp_access_control | cloudavenue_vapp_acl | :white_check_mark: |
| (87) | vcd_vapp_firewall_rules | Not yet implemented | :x: |
| (88) | vcd_vapp_nat_rules | Not yet implemented | :x: |
| (89) | vcd_vapp_network | cloudavenue_vapp_isolated_network | :white_check_mark: |
| (90) | vcd_vapp_org_network | cloudavenue_vapp_org_network |:white_check_mark: |
| (91) | vcd_vapp_static_routing | Not yet implemented | :x: |
| (92) | vcd_vapp_vm | cloudavenue_vm | :white_check_mark: |
| (93) | vcd_vdc_group | Not Applicable | :heavy_multiplication_x: |
| (94) | vcd_vm | Not Applicable | :heavy_multiplication_x: |
| (95) | vcd_vm_affinity_rule | cloudavenue_vm_affinity_rule |:white_check_mark: |
| (96) | vcd_vm_internal_disk | cloudavenue_vm_disk | :white_check_mark: |
| (97) | vcd_vm_placement_policy | Not Applicable | :heavy_multiplication_x: |
| (98) | vcd_vm_sizing_policy | Not Applicable | :heavy_multiplication_x: |

| Number | Datasources VMware VCD | Datasources Orange Cloud Avenue | status |
|:--:|:--:|:--:|:--:|
| (1) | vcd_catalog | cloudavenue_catalog |:white_check_mark: |
| (2) | vcd_catalog_item | Not Applicable | :heavy_multiplication_x: |
| (3) | vcd_catalog_media | cloudavenue_catalog_media |:white_check_mark: |
| (4) | vcd_catalog_vapp_template | cloudavenue_catalog_vapp_template |:white_check_mark: |
| (5) | vcd_edgegateway | Not Applicable | :heavy_multiplication_x: |
| (6) | vcd_external_network | Not Applicable | :heavy_multiplication_x: |
| (7) | vcd_external_network_v2 | Not Applicable | :heavy_multiplication_x: |
| (8) | vcd_global_role | Not Applicable | :heavy_multiplication_x: |
| (9) | vcd_independent_disk | cloudavenue_vm_disk | :white_check_mark: |
| (10) | vcd_ip_space | Not Applicable | :heavy_multiplication_x: |
| (11) | vcd_ip_space_custom_quota | Not Applicable | :heavy_multiplication_x: |
| (12) | vcd_ip_space_ip_allocation | Not yet implemented | :x: |
| (13) | vcd_ip_space_uplink | Not Applicable | :heavy_multiplication_x: |
| (14) | vcd_lb_app_profile | Not Applicable | :heavy_multiplication_x: |
| (15) | vcd_lb_app_rule | Not Applicable | :heavy_multiplication_x: |
| (16) | vcd_lb_server_pool | Not Applicable | :heavy_multiplication_x: |
| (17) | vcd_lb_service_monitor | Not Applicable | :heavy_multiplication_x: |
| (18) | vcd_lb_virtual_server | Not Applicable | :heavy_multiplication_x: |
| (19) | vcd_library_certificate | Not yet implemented | :x: |
| (20) | vcd_network_direct | Not Applicable | :heavy_multiplication_x: |
| (21) | vcd_network_isolated | Not Applicable | :heavy_multiplication_x: |
| (22) | vcd_network_isolated_v2 | cloudavenue_network_isolated | :white_check_mark: |
| (23) | vcd_network_pool | Not yet implemented | :x: |
| (24) | vcd_network_routed | Not Applicable | :heavy_multiplication_x: |
| (25) | vcd_network_routed_v2 | cloudavenue_network_routed | :white_check_mark: |
| (26) | vcd_nsxt_alb_cloud | Not Applicable | :heavy_multiplication_x: |
| (27) | vcd_nsxt_alb_controller | Not Applicable | :heavy_multiplication_x: |
| (28) | vcd_nsxt_alb_edgegateway_service_engine_group | Not Applicable | :heavy_multiplication_x: |
| (29) | vcd_nsxt_alb_importable_cloud | Not yet implemented | :x: |
| (30) | vcd_nsxt_alb_pool | cloudavenue_alb_pool | :white_check_mark: |
| (31) | vcd_nsxt_alb_service_engine_group | Not Applicable | :heavy_multiplication_x: |
| (32) | vcd_nsxt_alb_settings | Not Applicable | :heavy_multiplication_x: |
| (33) | vcd_nsxt_alb_virtual_service | Not Applicable | :heavy_multiplication_x: |
| (34) | vcd_nsxt_app_port_profile | cloudavenue_edgegateway_app_port_profile | :white_check_mark: |
| (35) | vcd_nsxt_distributed_firewall | Not yet implemented | :x: |
| (36) | vcd_nsxt_distributed_firewall_rule | Not yet implemented | :x: |
| (37) | vcd_nsxt_dynamic_security_group | Not yet implemented | :x: |
| (38) | vcd_nsxt_edge_cluster | Not yet implemented | :x: |
| (39) | vcd_nsxt_edgegateway | cloudavenue_edgegateway | :white_check_mark: |
| (40) | vcd_nsxt_edgegateway_bgp_configuration | Not Applicable | :heavy_multiplication_x: |
| (41) | vcd_nsxt_edgegateway_bgp_ip_prefix_list | Not Applicable | :heavy_multiplication_x: |
| (42) | vcd_nsxt_edgegateway_bgp_neighbor | Not Applicable | :heavy_multiplication_x: |
| (43) | vcd_nsxt_edgegateway_dhcp_forwarding | cloudavenue_edgegateway_dhcp_forwarding | :white_check_mark: |
| (44) | vcd_nsxt_edgegateway_dhcpv6 | Not Applicable | :heavy_multiplication_x: |
| (45) | vcd_nsxt_edgegateway_qos_profile | Not yet implemented | :x: |
| (46) | vcd_nsxt_edgegateway_rate_limiting | Not yet implemented | :x: |
| (47) | vcd_nsxt_edgegateway_static_route | cloudavenue_edgegateway_static_route | :white_check_mark: |
| (48) | vcd_nsxt_firewall | cloudavenue_edgegateway_firewall | :white_check_mark: |
| (49) | vcd_nsxt_ip_set | cloudavenue_edgegateway_ip_set | :white_check_mark: |
| (50) | vcd_nsxt_ipsec_vpn_tunnel | Not yet implemented | :x: |
| (51) | vcd_nsxt_manager | Not Applicable | :heavy_multiplication_x: |
| (52) | vcd_nsxt_nat_rule | Not yet implemented | :x: |
| (53) | vcd_nsxt_network_context_profile | Not yet implemented | :x: |
| (54) | vcd_nsxt_network_dhcp | cloudavenue_network_dhcp | :white_check_mark: |
| (55) | vcd_nsxt_network_dhcp_binding | cloudavenue_network_dhcp_binding | :white_check_mark: |
| (56) | vcd_nsxt_network_imported | Not Applicable | :heavy_multiplication_x: |
| (57) | vcd_nsxt_route_advertisement | Not Applicable | :heavy_multiplication_x: |
| (58) | vcd_nsxt_security_group | cloudavenue_edgegateway_security_group | :white_check_mark: |
| (59) | vcd_nsxt_tier0_router | Not yet implemented | :x: |
| (60) | vcd_nsxv_application | Not yet implemented | :x: |
| (61) | vcd_nsxv_application_finder | Not yet implemented | :x: |
| (62) | vcd_nsxv_application_group | Not yet implemented | :x: |
| (63) | vcd_nsxv_dhcp_relay | Not Applicable | :heavy_multiplication_x: |
| (64) | vcd_nsxv_distributed_firewall | Not Applicable | :heavy_multiplication_x: |
| (65) | vcd_nsxv_dnat | Not Applicable | :heavy_multiplication_x: |
| (66) | vcd_nsxv_firewall_rule | Not Applicable | :heavy_multiplication_x: |
| (67) | vcd_nsxv_ip_set | Not Applicable | :heavy_multiplication_x: |
| (68) | vcd_nsxv_snat | Not Applicable | :heavy_multiplication_x: |
| (69) | vcd_org | Not Applicable | :heavy_multiplication_x: |
| (70) | vcd_org_group | Not Applicable | :heavy_multiplication_x: |
| (71) | vcd_org_ldap | Not Applicable | :heavy_multiplication_x: |
| (72) | vcd_org_saml | Not Applicable | :heavy_multiplication_x: |
| (73) | vcd_org_saml_metadata | Not yet implemented | :x: |
| (74) | vcd_org_user | cloudavenue_iam_user | :white_check_mark: |
| (75) | vcd_org_vdc | Not Applicable | :heavy_multiplication_x: |
| (76) | vcd_portgroup | Not yet implemented | :x: |
| (77) | vcd_provider_vdc | Not Applicable | :heavy_multiplication_x: |
| (78) | vcd_rde | Not Applicable | :heavy_multiplication_x: |
| (79) | vcd_rde_interface | Not Applicable | :heavy_multiplication_x: |
| (80) | vcd_rde_interface_behavior | Not Applicable | :heavy_multiplication_x: |
| (81) | vcd_rde_type | Not Applicable | :heavy_multiplication_x: |
| (82) | vcd_rde_type_behavior | Not Applicable | :heavy_multiplication_x: |
| (83) | vcd_rde_type_behavior_acl | Not Applicable | :heavy_multiplication_x: |
| (84) | vcd_resource_list | Not Applicable | :heavy_multiplication_x: |
| (85) | vcd_resource_pool | Not Applicable | :heavy_multiplication_x: |
| (86) | vcd_resource_schema | Not Applicable | :heavy_multiplication_x: |
| (87) | vcd_right | Not yet implemented | :x: |
| (88) | vcd_rights_bundle | Not Applicable | :heavy_multiplication_x: |
| (89) | vcd_role | cloudavenue_iam_role | :white_check_mark: |
| (90) | vcd_service_account | Not yet implemented | :x: |
| (91) | vcd_storage_profile | cloudavenue_storage_profile |:white_check_mark: |
| (92) | vcd_subscribed_catalog | Not Applicable | :heavy_multiplication_x: |
| (93) | vcd_task | Not yet implemented | :x: |
| (94) | vcd_ui_plugin | Not Applicable | :heavy_multiplication_x: |
| (95) | vcd_vapp | cloudavenue_vapp |:white_check_mark: |
| (96) | vcd_vapp_network | cloudavenue_vapp_isolated_network | :white_check_mark: |
| (97) | vcd_vapp_org_network | cloudavenue_vapp_org_network |:white_check_mark: |
| (98) | vcd_vapp_vm | cloudavenue_vm | :white_check_mark: |
| (99) | vcd_vcenter | Not yet implemented | :x: |
| (100) | vcd_vdc_group | Not Applicable | :heavy_multiplication_x: |
| (101) | vcd_vm | Not Applicable | :heavy_multiplication_x: |
| (102) | vcd_vm_affinity_rule | cloudavenue_vm_affinity_rule |:white_check_mark: |
| (103) | vcd_vm_group | Not yet implemented | :x: |
| (104) | vcd_vm_placement_policy | Not Applicable | :heavy_multiplication_x: |
| (105) | vcd_vm_sizing_policy | Not Applicable | :heavy_multiplication_x: |
