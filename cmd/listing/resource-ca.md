# Checking resources and datasources of Orange Cloud Avenue provider
- Found 26 resources in terraform
- Found 31 datasources in terraform

# Checking resources and datasources of VMware Cloud Director provider
- Found 98 resources in terraform
- Found 105 datasources in terraform


# Listing cross resources and datasources from Cloud Avenue

| Number | Resources Orange Cloud Avenue | Resources VMware VCD |
|:--:|:--:|:--:|
| (1) | cloudavenue_alb_pool | vcd_nsxt_alb_pool |
| (2) | cloudavenue_catalog | vcd_catalog |
| (3) | cloudavenue_edgegateway | vcd_edgegateway |
| (4) | cloudavenue_edgegateway_app_port_profile | vcd_nsxt_app_port_profile |
| (5) | cloudavenue_edgegateway_firewall | vcd_nsxt_firewall |
| (6) | cloudavenue_edgegateway_ip_set | vcd_nsxt_ip_set |
| (7) | cloudavenue_edgegateway_security_group | vcd_nsxt_security_group |
| (8) | cloudavenue_iam_role | vcd_role |
| (9) | cloudavenue_iam_user | vcd_org_user |
| (10) | cloudavenue_network_dhcp |
| (11) | cloudavenue_network_dhcp_binding | vcd_nsxt_network_dhcp_binding  |
| (12) | cloudavenue_network_isolated | vcd_network_isolated |
| (13) | cloudavenue_network_routed | vcd_network_routed |
| (14) | cloudavenue_publicip |
| (15) | cloudavenue_vapp | vcd_vapp |
| (16) | cloudavenue_vapp_acl | vcd_vapp_access_control |
| (17) | cloudavenue_vapp_isolated_network | vcd_vapp_network |
| (18) | cloudavenue_vapp_org_network | vcd_vapp_org_network |
| (19) | cloudavenue_vcda_ip |
| (20) | cloudavenue_vdc | vcd_org_vdc |
| (21) | cloudavenue_vdc_acl | vcd_org_vdc_access_control |
| (22) | cloudavenue_vm | vcd_vm |
| (23) | cloudavenue_vm_affinity_rule | vcd_vm_affinity_rule |
| (24) | cloudavenue_vm_disk | vcd_vm_internal_disk |
| (25) | cloudavenue_vm_inserted_media | vcd_inserted_media |
| (26) | cloudavenue_vm_security_tag | vcd_security_tag |

| Number | Datasources Orange Cloud Avenue | Datasources VMware VCD |
|:--:|:--:|:--:|
| (1) | cloudavenue_alb_pool | vcd_nsxt_alb_pool |
| (2) | cloudavenue_catalog | vcd_catalog |
| (3) | cloudavenue_catalog_media | vcd_catalog_media |
| (4) | cloudavenue_catalog_medias |
| (5) | cloudavenue_catalog_vapp_template | vcd_catalog_vapp_template |
| (6) | cloudavenue_catalogs |
| (7) | cloudavenue_edgegateway | vcd_edgegateway |
| (8) | cloudavenue_edgegateway_firewall | vcd_nsxt_firewall |
| (9) | cloudavenue_edgegateway_ip_set | vcd_nsxt_ip_set |
| (10) | cloudavenue_edgegateway_security_group | vcd_nsxt_security_group |
| (11) | cloudavenue_edgegateways |
| (12) | cloudavenue_iam_right |
| (13) | cloudavenue_iam_role | vcd_role |
| (14) | cloudavenue_iam_user | vcd_org_user |
| (15) | cloudavenue_network_dhcp |
| (16) | cloudavenue_network_dhcp_binding | vcd_nsxt_network_dhcp_binding  |
| (17) | cloudavenue_network_isolated | vcd_network_isolated |
| (18) | cloudavenue_network_routed | vcd_network_routed |
| (19) | cloudavenue_publicips |
| (20) | cloudavenue_storage_profile | vcd_storage_profile |
| (21) | cloudavenue_storage_profiles |
| (22) | cloudavenue_tier0_vrf |
| (23) | cloudavenue_tier0_vrfs |
| (24) | cloudavenue_vapp | vcd_vapp |
| (25) | cloudavenue_vapp_isolated_network | vcd_vapp_network |
| (26) | cloudavenue_vapp_org_network | vcd_vapp_org_network |
| (27) | cloudavenue_vdc | vcd_org_vdc |
| (28) | cloudavenue_vdc_group | vcd_vdc_group |
| (29) | cloudavenue_vdcs |
| (30) | cloudavenue_vm | vcd_vm |
| (31) | cloudavenue_vm_affinity_rule | vcd_vm_affinity_rule |

# Listing cross resources and datasources from VCD

| Number | Resources VMware VCD | Resources Orange Cloud Avenue | status |
|:--:|:--:|:--:|:--:|
| (1) | vcd_api_token | Not yet implemented | :x: |
| (2) | vcd_catalog | cloudavenue_catalog |:white_check_mark: |
| (3) | vcd_catalog_access_control | Not yet implemented | :x: |
| (4) | vcd_catalog_item | Not Applicable | :no_entry: |
| (5) | vcd_catalog_media | Not yet implemented | :x: |
| (6) | vcd_catalog_vapp_template | Not yet implemented | :x: |
| (7) | vcd_cloned_vapp | Not yet implemented | :x: |
| (8) | vcd_edgegateway | Not Applicable | :no_entry: |
| (9) | vcd_edgegateway_settings | Not Applicable | :no_entry: |
| (10) | vcd_edgegateway_vpn | Not Applicable | :no_entry: |
| (11) | vcd_external_network | Not Applicable | :no_entry: |
| (12) | vcd_external_network_v2 | Not Applicable | :no_entry: |
| (13) | vcd_global_role | Not Applicable | :no_entry: |
| (14) | vcd_independent_disk | cloudavenue_vm_disk | :white_check_mark: |
| (15) | vcd_inserted_media | cloudavenue_vm_inserted_media | :white_check_mark: |
| (16) | vcd_ip_space | Not yet implemented | :x: |
| (17) | vcd_ip_space_custom_quota | Not yet implemented | :x: |
| (18) | vcd_ip_space_ip_allocation | Not yet implemented | :x: |
| (19) | vcd_ip_space_uplink | Not yet implemented | :x: |
| (20) | vcd_lb_app_profile | Not Applicable | :no_entry: |
| (21) | vcd_lb_app_rule | Not Applicable | :no_entry: |
| (22) | vcd_lb_server_pool | Not Applicable | :no_entry: |
| (23) | vcd_lb_service_monitor | Not Applicable | :no_entry: |
| (24) | vcd_lb_virtual_server | Not Applicable | :no_entry: |
| (25) | vcd_library_certificate | Not yet implemented | :x: |
| (26) | vcd_network_direct | Not Applicable | :no_entry: |
| (27) | vcd_network_isolated | Not Applicable | :no_entry: |
| (28) | vcd_network_isolated_v2 | cloudavenue_network_isolated | :white_check_mark: |
| (29) | vcd_network_routed | Not Applicable | :no_entry: |
| (30) | vcd_network_routed_v2 | cloudavenue_network_routed | :white_check_mark: |
| (31) | vcd_nsxt_alb_cloud | Not Applicable | :no_entry: |
| (32) | vcd_nsxt_alb_controller | Not Applicable | :no_entry: |
| (33) | vcd_nsxt_alb_edgegateway_service_engine_group | Not Applicable | :no_entry: |
| (34) | vcd_nsxt_alb_pool | cloudavenue_alb_pool | :white_check_mark: |
| (35) | vcd_nsxt_alb_service_engine_group | Not Applicable | :no_entry: |
| (36) | vcd_nsxt_alb_settings | Not Applicable | :no_entry: |
| (37) | vcd_nsxt_alb_virtual_service | Not Applicable | :no_entry: |
| (38) | vcd_nsxt_app_port_profile | cloudavenue_edgegateway_app_port_profile | :white_check_mark: |
| (39) | vcd_nsxt_distributed_firewall | Not yet implemented | :x: |
| (40) | vcd_nsxt_distributed_firewall_rule | Not yet implemented | :x: |
| (41) | vcd_nsxt_dynamic_security_group | Not yet implemented | :x: |
| (42) | vcd_nsxt_edgegateway | cloudavenue_edgegateway | :white_check_mark: |
| (43) | vcd_nsxt_edgegateway_bgp_configuration | Not Applicable | :no_entry: |
| (44) | vcd_nsxt_edgegateway_bgp_ip_prefix_list | Not Applicable | :no_entry: |
| (45) | vcd_nsxt_edgegateway_bgp_neighbor | Not Applicable | :no_entry: |
| (46) | vcd_nsxt_edgegateway_dhcp_forwarding | Not yet implemented | :x: |
| (47) | vcd_nsxt_edgegateway_dhcpv6 | Not yet implemented | :x: |
| (48) | vcd_nsxt_edgegateway_rate_limiting | Not yet implemented | :x: |
| (49) | vcd_nsxt_edgegateway_static_route | Not yet implemented | :x: |
| (50) | vcd_nsxt_firewall | cloudavenue_edgegateway_firewall | :white_check_mark: |
| (51) | vcd_nsxt_ip_set | cloudavenue_edgegateway_ip_set | :white_check_mark: |
| (52) | vcd_nsxt_ipsec_vpn_tunnel | Not yet implemented | :x: |
| (53) | vcd_nsxt_nat_rule | Not yet implemented | :x: |
| (54) | vcd_nsxt_network_dhcp | Not yet implemented | :x: |
| (55) | vcd_nsxt_network_dhcp_binding | Not yet implemented | :x: |
| (56) | vcd_nsxt_network_imported | Not Applicable | :no_entry: |
| (57) | vcd_nsxt_route_advertisement | Not Applicable | :no_entry: |
| (58) | vcd_nsxt_security_group | cloudavenue_edgegateway_security_group | :white_check_mark: |
| (59) | vcd_nsxv_dhcp_relay | Not Applicable | :no_entry: |
| (60) | vcd_nsxv_distributed_firewall | Not Applicable | :no_entry: |
| (61) | vcd_nsxv_dnat | Not Applicable | :no_entry: |
| (62) | vcd_nsxv_firewall_rule | Not Applicable | :no_entry: |
| (63) | vcd_nsxv_ip_set | Not Applicable | :no_entry: |
| (64) | vcd_nsxv_snat | Not Applicable | :no_entry: |
| (65) | vcd_org | Not Applicable | :no_entry: |
| (66) | vcd_org_group | Not yet implemented | :x: |
| (67) | vcd_org_ldap | Not Applicable | :no_entry: |
| (68) | vcd_org_saml | Not yet implemented | :x: |
| (69) | vcd_org_user | cloudavenue_iam_user | :white_check_mark: |
| (70) | vcd_org_vdc | Not Applicable | :no_entry: |
| (71) | vcd_org_vdc_access_control | cloudavenue_vdc_acl | :white_check_mark: |
| (72) | vcd_provider_vdc | Not yet implemented | :x: |
| (73) | vcd_rde | Not yet implemented | :x: |
| (74) | vcd_rde_interface | Not yet implemented | :x: |
| (75) | vcd_rde_interface_behavior | Not yet implemented | :x: |
| (76) | vcd_rde_type | Not yet implemented | :x: |
| (77) | vcd_rde_type_behavior | Not yet implemented | :x: |
| (78) | vcd_rde_type_behavior_acl | Not yet implemented | :x: |
| (79) | vcd_rights_bundle | Not Applicable | :no_entry: |
| (80) | vcd_role | cloudavenue_iam_role | :white_check_mark: |
| (81) | vcd_security_tag | cloudavenue_vm_security_tag | :white_check_mark: |
| (82) | vcd_service_account | Not yet implemented | :x: |
| (83) | vcd_subscribed_catalog | Not Applicable | :no_entry: |
| (84) | vcd_ui_plugin | Not yet implemented | :x: |
| (85) | vcd_vapp | cloudavenue_vapp |:white_check_mark: |
| (86) | vcd_vapp_access_control | cloudavenue_vapp_acl | :white_check_mark: |
| (87) | vcd_vapp_firewall_rules | Not yet implemented | :x: |
| (88) | vcd_vapp_nat_rules | Not yet implemented | :x: |
| (89) | vcd_vapp_network | cloudavenue_vapp_isolated_network | :white_check_mark: |
| (90) | vcd_vapp_org_network | cloudavenue_vapp_org_network |:white_check_mark: |
| (91) | vcd_vapp_static_routing | Not yet implemented | :x: |
| (92) | vcd_vapp_vm | cloudavenue_vm | :white_check_mark: |
| (93) | vcd_vdc_group | Not Applicable | :no_entry: |
| (94) | vcd_vm | Not Applicable | :no_entry: |
| (95) | vcd_vm_affinity_rule | cloudavenue_vm_affinity_rule |:white_check_mark: |
| (96) | vcd_vm_internal_disk | cloudavenue_vm_disk | :white_check_mark: |
| (97) | vcd_vm_placement_policy | Not Applicable | :no_entry: |
| (98) | vcd_vm_sizing_policy | Not Applicable | :no_entry: |

| Number | Datasources VMware VCD | Datasources Orange Cloud Avenue | status |
|:--:|:--:|:--:|:--:|
| (1) | vcd_catalog | cloudavenue_catalog |:white_check_mark: |
| (2) | vcd_catalog_item | Not Applicable | :no_entry: |
| (3) | vcd_catalog_media | cloudavenue_catalog_media |:white_check_mark: |
| (4) | vcd_catalog_vapp_template | cloudavenue_catalog_vapp_template |:white_check_mark: |
| (5) | vcd_edgegateway | Not Applicable | :no_entry: |
| (6) | vcd_external_network | Not Applicable | :no_entry: |
| (7) | vcd_external_network_v2 | Not Applicable | :no_entry: |
| (8) | vcd_global_role | Not Applicable | :no_entry: |
| (9) | vcd_independent_disk | cloudavenue_vm_disk | :white_check_mark: |
| (10) | vcd_ip_space | Not yet implemented | :x: |
| (11) | vcd_ip_space_custom_quota | Not yet implemented | :x: |
| (12) | vcd_ip_space_ip_allocation | Not yet implemented | :x: |
| (13) | vcd_ip_space_uplink | Not yet implemented | :x: |
| (14) | vcd_lb_app_profile | Not Applicable | :no_entry: |
| (15) | vcd_lb_app_rule | Not Applicable | :no_entry: |
| (16) | vcd_lb_server_pool | Not Applicable | :no_entry: |
| (17) | vcd_lb_service_monitor | Not Applicable | :no_entry: |
| (18) | vcd_lb_virtual_server | Not Applicable | :no_entry: |
| (19) | vcd_library_certificate | Not yet implemented | :x: |
| (20) | vcd_network_direct | Not Applicable | :no_entry: |
| (21) | vcd_network_isolated | Not Applicable | :no_entry: |
| (22) | vcd_network_isolated_v2 | cloudavenue_network_isolated | :white_check_mark: |
| (23) | vcd_network_pool | Not yet implemented | :x: |
| (24) | vcd_network_routed | Not Applicable | :no_entry: |
| (25) | vcd_network_routed_v2 | cloudavenue_network_routed | :white_check_mark: |
| (26) | vcd_nsxt_alb_cloud | Not Applicable | :no_entry: |
| (27) | vcd_nsxt_alb_controller | Not Applicable | :no_entry: |
| (28) | vcd_nsxt_alb_edgegateway_service_engine_group | Not Applicable | :no_entry: |
| (29) | vcd_nsxt_alb_importable_cloud | Not yet implemented | :x: |
| (30) | vcd_nsxt_alb_pool | cloudavenue_alb_pool | :white_check_mark: |
| (31) | vcd_nsxt_alb_service_engine_group | Not Applicable | :no_entry: |
| (32) | vcd_nsxt_alb_settings | Not Applicable | :no_entry: |
| (33) | vcd_nsxt_alb_virtual_service | Not Applicable | :no_entry: |
| (34) | vcd_nsxt_app_port_profile | cloudavenue_edgegateway_app_port_profile | :white_check_mark: |
| (35) | vcd_nsxt_distributed_firewall | Not yet implemented | :x: |
| (36) | vcd_nsxt_distributed_firewall_rule | Not yet implemented | :x: |
| (37) | vcd_nsxt_dynamic_security_group | Not yet implemented | :x: |
| (38) | vcd_nsxt_edge_cluster | Not yet implemented | :x: |
| (39) | vcd_nsxt_edgegateway | cloudavenue_edgegateway | :white_check_mark: |
| (40) | vcd_nsxt_edgegateway_bgp_configuration | Not Applicable | :no_entry: |
| (41) | vcd_nsxt_edgegateway_bgp_ip_prefix_list | Not Applicable | :no_entry: |
| (42) | vcd_nsxt_edgegateway_bgp_neighbor | Not Applicable | :no_entry: |
| (43) | vcd_nsxt_edgegateway_dhcp_forwarding | Not yet implemented | :x: |
| (44) | vcd_nsxt_edgegateway_dhcpv6 | Not yet implemented | :x: |
| (45) | vcd_nsxt_edgegateway_qos_profile | Not yet implemented | :x: |
| (46) | vcd_nsxt_edgegateway_rate_limiting | Not yet implemented | :x: |
| (47) | vcd_nsxt_edgegateway_static_route | Not yet implemented | :x: |
| (48) | vcd_nsxt_firewall | cloudavenue_edgegateway_firewall | :white_check_mark: |
| (49) | vcd_nsxt_ip_set | cloudavenue_edgegateway_ip_set | :white_check_mark: |
| (50) | vcd_nsxt_ipsec_vpn_tunnel | Not yet implemented | :x: |
| (51) | vcd_nsxt_manager | Not yet implemented | :x: |
| (52) | vcd_nsxt_nat_rule | Not yet implemented | :x: |
| (53) | vcd_nsxt_network_context_profile | Not yet implemented | :x: |
| (54) | vcd_nsxt_network_dhcp | Not yet implemented | :x: |
| (55) | vcd_nsxt_network_dhcp_binding | Not yet implemented | :x: |
| (56) | vcd_nsxt_network_imported | Not Applicable | :no_entry: |
| (57) | vcd_nsxt_route_advertisement | Not Applicable | :no_entry: |
| (58) | vcd_nsxt_security_group | cloudavenue_edgegateway_security_group | :white_check_mark: |
| (59) | vcd_nsxt_tier0_router | Not yet implemented | :x: |
| (60) | vcd_nsxv_application | Not yet implemented | :x: |
| (61) | vcd_nsxv_application_finder | Not yet implemented | :x: |
| (62) | vcd_nsxv_application_group | Not yet implemented | :x: |
| (63) | vcd_nsxv_dhcp_relay | Not Applicable | :no_entry: |
| (64) | vcd_nsxv_distributed_firewall | Not Applicable | :no_entry: |
| (65) | vcd_nsxv_dnat | Not Applicable | :no_entry: |
| (66) | vcd_nsxv_firewall_rule | Not Applicable | :no_entry: |
| (67) | vcd_nsxv_ip_set | Not Applicable | :no_entry: |
| (68) | vcd_nsxv_snat | Not Applicable | :no_entry: |
| (69) | vcd_org | Not Applicable | :no_entry: |
| (70) | vcd_org_group | Not yet implemented | :x: |
| (71) | vcd_org_ldap | Not Applicable | :no_entry: |
| (72) | vcd_org_saml | Not yet implemented | :x: |
| (73) | vcd_org_saml_metadata | Not yet implemented | :x: |
| (74) | vcd_org_user | cloudavenue_iam_user | :white_check_mark: |
| (75) | vcd_org_vdc | Not Applicable | :no_entry: |
| (76) | vcd_portgroup | Not yet implemented | :x: |
| (77) | vcd_provider_vdc | Not yet implemented | :x: |
| (78) | vcd_rde | Not yet implemented | :x: |
| (79) | vcd_rde_interface | Not yet implemented | :x: |
| (80) | vcd_rde_interface_behavior | Not yet implemented | :x: |
| (81) | vcd_rde_type | Not yet implemented | :x: |
| (82) | vcd_rde_type_behavior | Not yet implemented | :x: |
| (83) | vcd_rde_type_behavior_acl | Not yet implemented | :x: |
| (84) | vcd_resource_list | Not yet implemented | :x: |
| (85) | vcd_resource_pool | Not yet implemented | :x: |
| (86) | vcd_resource_schema | Not yet implemented | :x: |
| (87) | vcd_right | Not yet implemented | :x: |
| (88) | vcd_rights_bundle | Not Applicable | :no_entry: |
| (89) | vcd_role | cloudavenue_iam_role | :white_check_mark: |
| (90) | vcd_service_account | Not yet implemented | :x: |
| (91) | vcd_storage_profile | cloudavenue_storage_profile |:white_check_mark: |
| (92) | vcd_subscribed_catalog | Not Applicable | :no_entry: |
| (93) | vcd_task | Not yet implemented | :x: |
| (94) | vcd_ui_plugin | Not yet implemented | :x: |
| (95) | vcd_vapp | cloudavenue_vapp |:white_check_mark: |
| (96) | vcd_vapp_network | cloudavenue_vapp_isolated_network | :white_check_mark: |
| (97) | vcd_vapp_org_network | cloudavenue_vapp_org_network |:white_check_mark: |
| (98) | vcd_vapp_vm | cloudavenue_vm | :white_check_mark: |
| (99) | vcd_vcenter | Not yet implemented | :x: |
| (100) | vcd_vdc_group | Not Applicable | :no_entry: |
| (101) | vcd_vm | Not Applicable | :no_entry: |
| (102) | vcd_vm_affinity_rule | cloudavenue_vm_affinity_rule |:white_check_mark: |
| (103) | vcd_vm_group | Not yet implemented | :x: |
| (104) | vcd_vm_placement_policy | Not Applicable | :no_entry: |
| (105) | vcd_vm_sizing_policy | Not Applicable | :no_entry: |
