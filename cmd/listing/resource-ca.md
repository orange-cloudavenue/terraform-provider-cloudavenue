# Checking resources and datasources of Orange Cloud Avenue provider
- Found 22 resources in terraform
- Found 27 datasources in terraform

# Checking resources and datasources of VMware Cloud Director provider
- Found 81 resources in terraform
- Found 88 datasources in terraform


# Listing cross resources and datasources from Cloud Avenue

| Number | Resources Orange Cloud Avenue | Resources VMware VCD |
|:--:|:--:|:--:|
| (1) | cloudavenue_alb_pool | vcd_nsxt_alb_pool |
| (2) | cloudavenue_catalog | vcd_catalog |
| (3) | cloudavenue_edgegateway | vcd_edgegateway |
| (4) | cloudavenue_edgegateway_app_port_profile | vcd_nsxt_app_port_profile |
| (5) | cloudavenue_edgegateway_firewall | vcd_nsxt_firewall |
| (6) | cloudavenue_iam_role | vcd_role |
| (7) | cloudavenue_iam_user | vcd_org_user |
| (8) | cloudavenue_network_isolated | vcd_network_isolated |
| (9) | cloudavenue_network_routed | vcd_network_routed |
| (10) | cloudavenue_publicip |
| (11) | cloudavenue_vapp | vcd_vapp |
| (12) | cloudavenue_vapp_acl | vcd_vapp_access_control |
| (13) | cloudavenue_vapp_isolated_network | vcd_vapp_network |
| (14) | cloudavenue_vapp_org_network | vcd_vapp_org_network |
| (15) | cloudavenue_vcda_ip |
| (16) | cloudavenue_vdc | vcd_org_vdc |
| (17) | cloudavenue_vdc_acl | vcd_org_vdc_access_control |
| (18) | cloudavenue_vm | vcd_vm |
| (19) | cloudavenue_vm_affinity_rule | vcd_vm_affinity_rule |
| (20) | cloudavenue_vm_disk | vcd_vm_internal_disk |
| (21) | cloudavenue_vm_inserted_media | vcd_inserted_media |
| (22) | cloudavenue_vm_security_tag | vcd_security_tag |

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
| (9) | cloudavenue_edgegateways |
| (10) | cloudavenue_iam_right |
| (11) | cloudavenue_iam_role | vcd_role |
| (12) | cloudavenue_iam_user | vcd_org_user |
| (13) | cloudavenue_network_isolated | vcd_network_isolated |
| (14) | cloudavenue_network_routed | vcd_network_routed |
| (15) | cloudavenue_publicips |
| (16) | cloudavenue_storage_profile | vcd_storage_profile |
| (17) | cloudavenue_storage_profiles |
| (18) | cloudavenue_tier0_vrf |
| (19) | cloudavenue_tier0_vrfs |
| (20) | cloudavenue_vapp | vcd_vapp |
| (21) | cloudavenue_vapp_isolated_network | vcd_vapp_network |
| (22) | cloudavenue_vapp_org_network | vcd_vapp_org_network |
| (23) | cloudavenue_vdc | vcd_org_vdc |
| (24) | cloudavenue_vdc_group | vcd_vdc_group |
| (25) | cloudavenue_vdcs |
| (26) | cloudavenue_vm | vcd_vm |
| (27) | cloudavenue_vm_affinity_rule | vcd_vm_affinity_rule |

# Listing cross resources and datasources from VCD

| Number | Resources VMware VCD | Resources Orange Cloud Avenue | status |
|:--:|:--:|:--:|:--:|
| (1) | vcd_catalog | cloudavenue_catalog |:white_check_mark: |
| (2) | vcd_catalog_access_control | Not yet implemented | :x: |
| (3) | vcd_catalog_item | Not Applicable | :no_entry: |
| (4) | vcd_catalog_media | Not yet implemented | :x: |
| (5) | vcd_catalog_vapp_template | Not yet implemented | :x: |
| (6) | vcd_edgegateway | Not Applicable | :no_entry: |
| (7) | vcd_edgegateway_settings | Not Applicable | :no_entry: |
| (8) | vcd_edgegateway_vpn | Not Applicable | :no_entry: |
| (9) | vcd_external_network | Not Applicable | :no_entry: |
| (10) | vcd_external_network_v2 | Not Applicable | :no_entry: |
| (11) | vcd_global_role | Not Applicable | :no_entry: |
| (12) | vcd_independent_disk | cloudavenue_vm_disk | :white_check_mark: |
| (13) | vcd_inserted_media | cloudavenue_vm_inserted_media | :white_check_mark: |
| (14) | vcd_lb_app_profile | Not Applicable | :no_entry: |
| (15) | vcd_lb_app_rule | Not Applicable | :no_entry: |
| (16) | vcd_lb_server_pool | Not Applicable | :no_entry: |
| (17) | vcd_lb_service_monitor | Not Applicable | :no_entry: |
| (18) | vcd_lb_virtual_server | Not Applicable | :no_entry: |
| (19) | vcd_library_certificate | Not yet implemented | :x: |
| (20) | vcd_network_direct | Not Applicable | :no_entry: |
| (21) | vcd_network_isolated | Not Applicable | :no_entry: |
| (22) | vcd_network_isolated_v2 | cloudavenue_network_isolated | :white_check_mark: |
| (23) | vcd_network_routed | Not Applicable | :no_entry: |
| (24) | vcd_network_routed_v2 | cloudavenue_network_routed | :white_check_mark: |
| (25) | vcd_nsxt_alb_cloud | Not Applicable | :no_entry: |
| (26) | vcd_nsxt_alb_controller | Not Applicable | :no_entry: |
| (27) | vcd_nsxt_alb_edgegateway_service_engine_group | Not Applicable | :no_entry: |
| (28) | vcd_nsxt_alb_pool | cloudavenue_alb_pool | :white_check_mark: |
| (29) | vcd_nsxt_alb_service_engine_group | Not Applicable | :no_entry: |
| (30) | vcd_nsxt_alb_settings | Not Applicable | :no_entry: |
| (31) | vcd_nsxt_alb_virtual_service | Not Applicable | :no_entry: |
| (32) | vcd_nsxt_app_port_profile | cloudavenue_edgegateway_app_port_profile | :white_check_mark: |
| (33) | vcd_nsxt_distributed_firewall | Not yet implemented | :x: |
| (34) | vcd_nsxt_dynamic_security_group | Not yet implemented | :x: |
| (35) | vcd_nsxt_edgegateway | cloudavenue_edgegateway | :white_check_mark: |
| (36) | vcd_nsxt_edgegateway_bgp_configuration | Not Applicable | :no_entry: |
| (37) | vcd_nsxt_edgegateway_bgp_ip_prefix_list | Not Applicable | :no_entry: |
| (38) | vcd_nsxt_edgegateway_bgp_neighbor | Not Applicable | :no_entry: |
| (39) | vcd_nsxt_edgegateway_rate_limiting | Not yet implemented | :x: |
| (40) | vcd_nsxt_firewall | cloudavenue_edgegateway_firewall | :white_check_mark: |
| (41) | vcd_nsxt_ip_set | Not yet implemented | :x: |
| (42) | vcd_nsxt_ipsec_vpn_tunnel | Not yet implemented | :x: |
| (43) | vcd_nsxt_nat_rule | Not yet implemented | :x: |
| (44) | vcd_nsxt_network_dhcp | Not yet implemented | :x: |
| (45) | vcd_nsxt_network_dhcp_binding | Not yet implemented | :x: |
| (46) | vcd_nsxt_network_imported | Not Applicable | :no_entry: |
| (47) | vcd_nsxt_route_advertisement | Not Applicable | :no_entry: |
| (48) | vcd_nsxt_security_group | Not yet implemented | :x: |
| (49) | vcd_nsxv_dhcp_relay | Not Applicable | :no_entry: |
| (50) | vcd_nsxv_distributed_firewall | Not Applicable | :no_entry: |
| (51) | vcd_nsxv_dnat | Not Applicable | :no_entry: |
| (52) | vcd_nsxv_firewall_rule | Not Applicable | :no_entry: |
| (53) | vcd_nsxv_ip_set | Not Applicable | :no_entry: |
| (54) | vcd_nsxv_snat | Not Applicable | :no_entry: |
| (55) | vcd_org | Not Applicable | :no_entry: |
| (56) | vcd_org_group | Not yet implemented | :x: |
| (57) | vcd_org_ldap | Not Applicable | :no_entry: |
| (58) | vcd_org_user | cloudavenue_iam_user | :white_check_mark: |
| (59) | vcd_org_vdc | Not Applicable | :no_entry: |
| (60) | vcd_org_vdc_access_control | cloudavenue_vdc_acl | :white_check_mark: |
| (61) | vcd_rde | Not yet implemented | :x: |
| (62) | vcd_rde_interface | Not yet implemented | :x: |
| (63) | vcd_rde_type | Not yet implemented | :x: |
| (64) | vcd_rights_bundle | Not Applicable | :no_entry: |
| (65) | vcd_role | cloudavenue_iam_role | :white_check_mark: |
| (66) | vcd_security_tag | cloudavenue_vm_security_tag | :white_check_mark: |
| (67) | vcd_subscribed_catalog | Not Applicable | :no_entry: |
| (68) | vcd_vapp | cloudavenue_vapp |:white_check_mark: |
| (69) | vcd_vapp_access_control | cloudavenue_vapp_acl | :white_check_mark: |
| (70) | vcd_vapp_firewall_rules | Not yet implemented | :x: |
| (71) | vcd_vapp_nat_rules | Not yet implemented | :x: |
| (72) | vcd_vapp_network | cloudavenue_vapp_isolated_network | :white_check_mark: |
| (73) | vcd_vapp_org_network | cloudavenue_vapp_org_network |:white_check_mark: |
| (74) | vcd_vapp_static_routing | Not yet implemented | :x: |
| (75) | vcd_vapp_vm | cloudavenue_vm | :white_check_mark: |
| (76) | vcd_vdc_group | Not Applicable | :no_entry: |
| (77) | vcd_vm | Not Applicable | :no_entry: |
| (78) | vcd_vm_affinity_rule | cloudavenue_vm_affinity_rule |:white_check_mark: |
| (79) | vcd_vm_internal_disk | cloudavenue_vm_disk | :white_check_mark: |
| (80) | vcd_vm_placement_policy | Not Applicable | :no_entry: |
| (81) | vcd_vm_sizing_policy | Not Applicable | :no_entry: |

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
| (10) | vcd_lb_app_profile | Not Applicable | :no_entry: |
| (11) | vcd_lb_app_rule | Not Applicable | :no_entry: |
| (12) | vcd_lb_server_pool | Not Applicable | :no_entry: |
| (13) | vcd_lb_service_monitor | Not Applicable | :no_entry: |
| (14) | vcd_lb_virtual_server | Not Applicable | :no_entry: |
| (15) | vcd_library_certificate | Not yet implemented | :x: |
| (16) | vcd_network_direct | Not Applicable | :no_entry: |
| (17) | vcd_network_isolated | Not Applicable | :no_entry: |
| (18) | vcd_network_isolated_v2 | cloudavenue_network_isolated | :white_check_mark: |
| (19) | vcd_network_routed | Not Applicable | :no_entry: |
| (20) | vcd_network_routed_v2 | cloudavenue_network_routed | :white_check_mark: |
| (21) | vcd_nsxt_alb_cloud | Not Applicable | :no_entry: |
| (22) | vcd_nsxt_alb_controller | Not Applicable | :no_entry: |
| (23) | vcd_nsxt_alb_edgegateway_service_engine_group | Not Applicable | :no_entry: |
| (24) | vcd_nsxt_alb_importable_cloud | Not yet implemented | :x: |
| (25) | vcd_nsxt_alb_pool | cloudavenue_alb_pool | :white_check_mark: |
| (26) | vcd_nsxt_alb_service_engine_group | Not Applicable | :no_entry: |
| (27) | vcd_nsxt_alb_settings | Not Applicable | :no_entry: |
| (28) | vcd_nsxt_alb_virtual_service | Not Applicable | :no_entry: |
| (29) | vcd_nsxt_app_port_profile | cloudavenue_edgegateway_app_port_profile | :white_check_mark: |
| (30) | vcd_nsxt_distributed_firewall | Not yet implemented | :x: |
| (31) | vcd_nsxt_dynamic_security_group | Not yet implemented | :x: |
| (32) | vcd_nsxt_edge_cluster | Not yet implemented | :x: |
| (33) | vcd_nsxt_edgegateway | cloudavenue_edgegateway | :white_check_mark: |
| (34) | vcd_nsxt_edgegateway_bgp_configuration | Not Applicable | :no_entry: |
| (35) | vcd_nsxt_edgegateway_bgp_ip_prefix_list | Not Applicable | :no_entry: |
| (36) | vcd_nsxt_edgegateway_bgp_neighbor | Not Applicable | :no_entry: |
| (37) | vcd_nsxt_edgegateway_qos_profile | Not yet implemented | :x: |
| (38) | vcd_nsxt_edgegateway_rate_limiting | Not yet implemented | :x: |
| (39) | vcd_nsxt_firewall | cloudavenue_edgegateway_firewall | :white_check_mark: |
| (40) | vcd_nsxt_ip_set | Not yet implemented | :x: |
| (41) | vcd_nsxt_ipsec_vpn_tunnel | Not yet implemented | :x: |
| (42) | vcd_nsxt_manager | Not yet implemented | :x: |
| (43) | vcd_nsxt_nat_rule | Not yet implemented | :x: |
| (44) | vcd_nsxt_network_context_profile | Not yet implemented | :x: |
| (45) | vcd_nsxt_network_dhcp | Not yet implemented | :x: |
| (46) | vcd_nsxt_network_dhcp_binding | Not yet implemented | :x: |
| (47) | vcd_nsxt_network_imported | Not Applicable | :no_entry: |
| (48) | vcd_nsxt_route_advertisement | Not Applicable | :no_entry: |
| (49) | vcd_nsxt_security_group | Not yet implemented | :x: |
| (50) | vcd_nsxt_tier0_router | Not yet implemented | :x: |
| (51) | vcd_nsxv_application | Not yet implemented | :x: |
| (52) | vcd_nsxv_application_finder | Not yet implemented | :x: |
| (53) | vcd_nsxv_application_group | Not yet implemented | :x: |
| (54) | vcd_nsxv_dhcp_relay | Not Applicable | :no_entry: |
| (55) | vcd_nsxv_distributed_firewall | Not Applicable | :no_entry: |
| (56) | vcd_nsxv_dnat | Not Applicable | :no_entry: |
| (57) | vcd_nsxv_firewall_rule | Not Applicable | :no_entry: |
| (58) | vcd_nsxv_ip_set | Not Applicable | :no_entry: |
| (59) | vcd_nsxv_snat | Not Applicable | :no_entry: |
| (60) | vcd_org | Not Applicable | :no_entry: |
| (61) | vcd_org_group | Not yet implemented | :x: |
| (62) | vcd_org_ldap | Not Applicable | :no_entry: |
| (63) | vcd_org_user | cloudavenue_iam_user | :white_check_mark: |
| (64) | vcd_org_vdc | Not Applicable | :no_entry: |
| (65) | vcd_portgroup | Not yet implemented | :x: |
| (66) | vcd_provider_vdc | Not yet implemented | :x: |
| (67) | vcd_rde | Not yet implemented | :x: |
| (68) | vcd_rde_interface | Not yet implemented | :x: |
| (69) | vcd_rde_type | Not yet implemented | :x: |
| (70) | vcd_resource_list | Not yet implemented | :x: |
| (71) | vcd_resource_schema | Not yet implemented | :x: |
| (72) | vcd_right | Not yet implemented | :x: |
| (73) | vcd_rights_bundle | Not Applicable | :no_entry: |
| (74) | vcd_role | cloudavenue_iam_role | :white_check_mark: |
| (75) | vcd_storage_profile | cloudavenue_storage_profile |:white_check_mark: |
| (76) | vcd_subscribed_catalog | Not Applicable | :no_entry: |
| (77) | vcd_task | Not yet implemented | :x: |
| (78) | vcd_vapp | cloudavenue_vapp |:white_check_mark: |
| (79) | vcd_vapp_network | cloudavenue_vapp_isolated_network | :white_check_mark: |
| (80) | vcd_vapp_org_network | cloudavenue_vapp_org_network |:white_check_mark: |
| (81) | vcd_vapp_vm | cloudavenue_vm | :white_check_mark: |
| (82) | vcd_vcenter | Not yet implemented | :x: |
| (83) | vcd_vdc_group | Not Applicable | :no_entry: |
| (84) | vcd_vm | Not Applicable | :no_entry: |
| (85) | vcd_vm_affinity_rule | cloudavenue_vm_affinity_rule |:white_check_mark: |
| (86) | vcd_vm_group | Not yet implemented | :x: |
| (87) | vcd_vm_placement_policy | Not Applicable | :no_entry: |
| (88) | vcd_vm_sizing_policy | Not Applicable | :no_entry: |
