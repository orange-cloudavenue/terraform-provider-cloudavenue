# Checking resources and datasources of Orange Cloud Avenue provider
- Found 43 resources in terraform
- Found 49 datasources in terraform

# Checking resources and datasources of VMware Cloud Director provider (version: unset)
- Found 123 resources in terraform
- Found 141 datasources in terraform


# Listing cross resources and datasources from Cloud Avenue

| Number | Resources Orange Cloud Avenue | Resources VMware VCD |
|:--:|:--:|:--:|
| (1) | cloudavenue_alb_pool | vcd_nsxt_alb_pool |
| (2) | cloudavenue_backup |
| (3) | cloudavenue_catalog | vcd_catalog |
| (4) | cloudavenue_catalog_acl | vcd_catalog_access_control |
| (5) | cloudavenue_edgegateway | vcd_edgegateway |
| (6) | cloudavenue_edgegateway_app_port_profile | vcd_nsxt_app_port_profile |
| (7) | cloudavenue_edgegateway_dhcp_forwarding | vcd_nsxt_edgegateway_dhcp_forwarding |
| (8) | cloudavenue_edgegateway_firewall | vcd_nsxt_firewall |
| (9) | cloudavenue_edgegateway_ip_set | vcd_nsxt_ip_set |
| (10) | cloudavenue_edgegateway_nat_rule | vcd_nsxt_nat_rule |
| (11) | cloudavenue_edgegateway_security_group | vcd_nsxt_security_group |
| (12) | cloudavenue_edgegateway_static_route | vcd_nsxt_edgegateway_static_route |
| (13) | cloudavenue_edgegateway_vpn_ipsec | vcd_nsxt_ipsec_vpn_tunnel |
| (14) | cloudavenue_iam_role | vcd_role |
| (15) | cloudavenue_iam_token |
| (16) | cloudavenue_iam_user | vcd_org_user |
| (17) | cloudavenue_iam_user_saml |
| (18) | cloudavenue_network_dhcp | vcd_nsxt_network_dhcp |
| (19) | cloudavenue_network_dhcp_binding | vcd_nsxt_network_dhcp_binding |
| (20) | cloudavenue_network_isolated | vcd_network_isolated |
| (21) | cloudavenue_network_routed | vcd_network_routed |
| (22) | cloudavenue_publicip |
| (23) | cloudavenue_s3_bucket |
| (24) | cloudavenue_s3_bucket_acl |
| (25) | cloudavenue_s3_bucket_cors_configuration |
| (26) | cloudavenue_s3_bucket_lifecycle_configuration |
| (27) | cloudavenue_s3_bucket_policy |
| (28) | cloudavenue_s3_bucket_versioning_configuration |
| (29) | cloudavenue_s3_bucket_website_configuration |
| (30) | cloudavenue_s3_credential |
| (31) | cloudavenue_vapp | vcd_vapp |
| (32) | cloudavenue_vapp_acl | vcd_vapp_access_control |
| (33) | cloudavenue_vapp_isolated_network | vcd_vapp_network |
| (34) | cloudavenue_vapp_org_network | vcd_vapp_org_network |
| (35) | cloudavenue_vcda_ip |
| (36) | cloudavenue_vdc | vcd_org_vdc |
| (37) | cloudavenue_vdc_acl | vcd_org_vdc_access_control |
| (38) | cloudavenue_vdc_group | vcd_vdc_group |
| (39) | cloudavenue_vm | vcd_vm |
| (40) | cloudavenue_vm_affinity_rule | vcd_vm_affinity_rule |
| (41) | cloudavenue_vm_disk | vcd_vm_internal_disk |
| (42) | cloudavenue_vm_inserted_media | vcd_inserted_media |
| (43) | cloudavenue_vm_security_tag | vcd_security_tag |

| Number | Datasources Orange Cloud Avenue | Datasources VMware VCD |
|:--:|:--:|:--:|
| (1) | cloudavenue_alb_pool | vcd_nsxt_alb_pool |
| (2) | cloudavenue_backup |
| (3) | cloudavenue_bms |
| (4) | cloudavenue_catalog | vcd_catalog |
| (5) | cloudavenue_catalog_acl | vcd_catalog_access_control |
| (6) | cloudavenue_catalog_media | vcd_catalog_media |
| (7) | cloudavenue_catalog_medias |
| (8) | cloudavenue_catalog_vapp_template | vcd_catalog_vapp_template |
| (9) | cloudavenue_catalogs |
| (10) | cloudavenue_edgegateway | vcd_edgegateway |
| (11) | cloudavenue_edgegateway_app_port_profile | vcd_nsxt_app_port_profile |
| (12) | cloudavenue_edgegateway_dhcp_forwarding | vcd_nsxt_edgegateway_dhcp_forwarding |
| (13) | cloudavenue_edgegateway_firewall | vcd_nsxt_firewall |
| (14) | cloudavenue_edgegateway_ip_set | vcd_nsxt_ip_set |
| (15) | cloudavenue_edgegateway_nat_rule | vcd_nsxt_nat_rule |
| (16) | cloudavenue_edgegateway_security_group | vcd_nsxt_security_group |
| (17) | cloudavenue_edgegateway_static_route | vcd_nsxt_edgegateway_static_route |
| (18) | cloudavenue_edgegateway_vpn_ipsec | vcd_nsxt_ipsec_vpn_tunnel |
| (19) | cloudavenue_edgegateways |
| (20) | cloudavenue_iam_right | vcd_right |
| (21) | cloudavenue_iam_role | vcd_role |
| (22) | cloudavenue_iam_roles |
| (23) | cloudavenue_iam_user | vcd_org_user |
| (24) | cloudavenue_network_dhcp | vcd_nsxt_network_dhcp |
| (25) | cloudavenue_network_dhcp_binding | vcd_nsxt_network_dhcp_binding |
| (26) | cloudavenue_network_isolated | vcd_network_isolated |
| (27) | cloudavenue_network_routed | vcd_network_routed |
| (28) | cloudavenue_publicips |
| (29) | cloudavenue_s3_bucket |
| (30) | cloudavenue_s3_bucket_acl |
| (31) | cloudavenue_s3_bucket_cors_configuration |
| (32) | cloudavenue_s3_bucket_lifecycle_configuration |
| (33) | cloudavenue_s3_bucket_policy |
| (34) | cloudavenue_s3_bucket_versioning_configuration |
| (35) | cloudavenue_s3_bucket_website_configuration |
| (36) | cloudavenue_s3_user |
| (37) | cloudavenue_storage_profile | vcd_storage_profile |
| (38) | cloudavenue_storage_profiles |
| (39) | cloudavenue_tier0_vrf |
| (40) | cloudavenue_tier0_vrfs |
| (41) | cloudavenue_vapp | vcd_vapp |
| (42) | cloudavenue_vapp_isolated_network | vcd_vapp_network |
| (43) | cloudavenue_vapp_org_network | vcd_vapp_org_network |
| (44) | cloudavenue_vdc | vcd_org_vdc |
| (45) | cloudavenue_vdc_group | vcd_vdc_group |
| (46) | cloudavenue_vdcs |
| (47) | cloudavenue_vm | vcd_vm |
| (48) | cloudavenue_vm_affinity_rule | vcd_vm_affinity_rule |
| (49) | cloudavenue_vm_disks |

# Listing cross resources and datasources from VCD (version: unset)

| Number | Resources VMware VCD | Resources Orange Cloud Avenue | status |
|:--:|:--:|:--:|:--:|
| (1) | vcd_api_filter | Not Applicable | :heavy_multiplication_x: |
| (2) | vcd_api_token | Not Applicable | :heavy_multiplication_x: |
| (3) | vcd_catalog | cloudavenue_catalog |:white_check_mark: |
| (4) | vcd_catalog_access_control | cloudavenue_catalog_acl | :white_check_mark: |
| (5) | vcd_catalog_item | Not Applicable | :heavy_multiplication_x: |
| (6) | vcd_catalog_media | Not yet implemented | :x: |
| (7) | vcd_catalog_vapp_template | Not yet implemented | :x: |
| (8) | vcd_cloned_vapp | Not Applicable | :heavy_multiplication_x: |
| (9) | vcd_cse_kubernetes_cluster | Not Applicable | :heavy_multiplication_x: |
| (10) | vcd_dse_registry_configuration | Not Applicable | :heavy_multiplication_x: |
| (11) | vcd_dse_solution_publish | Not Applicable | :heavy_multiplication_x: |
| (12) | vcd_edgegateway | Not Applicable | :heavy_multiplication_x: |
| (13) | vcd_edgegateway_settings | Not Applicable | :heavy_multiplication_x: |
| (14) | vcd_edgegateway_vpn | Not Applicable | :heavy_multiplication_x: |
| (15) | vcd_external_endpoint | Not Applicable | :heavy_multiplication_x: |
| (16) | vcd_external_network | Not Applicable | :heavy_multiplication_x: |
| (17) | vcd_external_network_v2 | Not Applicable | :heavy_multiplication_x: |
| (18) | vcd_global_role | Not Applicable | :heavy_multiplication_x: |
| (19) | vcd_independent_disk | cloudavenue_vm_disk | :white_check_mark: |
| (20) | vcd_inserted_media | cloudavenue_vm_inserted_media | :white_check_mark: |
| (21) | vcd_ip_space | Not Applicable | :heavy_multiplication_x: |
| (22) | vcd_ip_space_custom_quota | Not Applicable | :heavy_multiplication_x: |
| (23) | vcd_ip_space_ip_allocation | Not Applicable | :heavy_multiplication_x: |
| (24) | vcd_ip_space_uplink | Not Applicable | :heavy_multiplication_x: |
| (25) | vcd_lb_app_profile | Not Applicable | :heavy_multiplication_x: |
| (26) | vcd_lb_app_rule | Not Applicable | :heavy_multiplication_x: |
| (27) | vcd_lb_server_pool | Not Applicable | :heavy_multiplication_x: |
| (28) | vcd_lb_service_monitor | Not Applicable | :heavy_multiplication_x: |
| (29) | vcd_lb_virtual_server | Not Applicable | :heavy_multiplication_x: |
| (30) | vcd_library_certificate | Not Applicable | :heavy_multiplication_x: |
| (31) | vcd_multisite_org_association | Not yet implemented | :x: |
| (32) | vcd_multisite_site_association | Not Applicable | :heavy_multiplication_x: |
| (33) | vcd_network_direct | Not Applicable | :heavy_multiplication_x: |
| (34) | vcd_network_isolated | Not Applicable | :heavy_multiplication_x: |
| (35) | vcd_network_isolated_v2 | cloudavenue_network_isolated | :white_check_mark: |
| (36) | vcd_network_pool | Not Applicable | :heavy_multiplication_x: |
| (37) | vcd_network_routed | Not Applicable | :heavy_multiplication_x: |
| (38) | vcd_network_routed_v2 | cloudavenue_network_routed | :white_check_mark: |
| (39) | vcd_nsxt_alb_cloud | Not Applicable | :heavy_multiplication_x: |
| (40) | vcd_nsxt_alb_controller | Not Applicable | :heavy_multiplication_x: |
| (41) | vcd_nsxt_alb_edgegateway_service_engine_group | Not Applicable | :heavy_multiplication_x: |
| (42) | vcd_nsxt_alb_pool | cloudavenue_alb_pool | :white_check_mark: |
| (43) | vcd_nsxt_alb_service_engine_group | Not Applicable | :heavy_multiplication_x: |
| (44) | vcd_nsxt_alb_settings | Not Applicable | :heavy_multiplication_x: |
| (45) | vcd_nsxt_alb_virtual_service | Not Applicable | :heavy_multiplication_x: |
| (46) | vcd_nsxt_alb_virtual_service_http_req_rules | Not yet implemented | :x: |
| (47) | vcd_nsxt_alb_virtual_service_http_resp_rules | Not yet implemented | :x: |
| (48) | vcd_nsxt_alb_virtual_service_http_sec_rules | Not yet implemented | :x: |
| (49) | vcd_nsxt_app_port_profile | cloudavenue_edgegateway_app_port_profile | :white_check_mark: |
| (50) | vcd_nsxt_distributed_firewall | Not yet implemented | :x: |
| (51) | vcd_nsxt_distributed_firewall_rule | Not yet implemented | :x: |
| (52) | vcd_nsxt_dynamic_security_group | Not yet implemented | :x: |
| (53) | vcd_nsxt_edgegateway | cloudavenue_edgegateway | :white_check_mark: |
| (54) | vcd_nsxt_edgegateway_bgp_configuration | Not Applicable | :heavy_multiplication_x: |
| (55) | vcd_nsxt_edgegateway_bgp_ip_prefix_list | Not Applicable | :heavy_multiplication_x: |
| (56) | vcd_nsxt_edgegateway_bgp_neighbor | Not Applicable | :heavy_multiplication_x: |
| (57) | vcd_nsxt_edgegateway_dhcp_forwarding | cloudavenue_edgegateway_dhcp_forwarding | :white_check_mark: |
| (58) | vcd_nsxt_edgegateway_dhcpv6 | Not Applicable | :heavy_multiplication_x: |
| (59) | vcd_nsxt_edgegateway_dns | Not yet implemented | :x: |
| (60) | vcd_nsxt_edgegateway_l2_vpn_tunnel | Not yet implemented | :x: |
| (61) | vcd_nsxt_edgegateway_rate_limiting | Not Applicable | :heavy_multiplication_x: |
| (62) | vcd_nsxt_edgegateway_static_route | cloudavenue_edgegateway_static_route | :white_check_mark: |
| (63) | vcd_nsxt_firewall | cloudavenue_edgegateway_firewall | :white_check_mark: |
| (64) | vcd_nsxt_global_default_segment_profile_template | Not Applicable | :heavy_multiplication_x: |
| (65) | vcd_nsxt_ip_set | cloudavenue_edgegateway_ip_set | :white_check_mark: |
| (66) | vcd_nsxt_ipsec_vpn_tunnel | cloudavenue_edgegateway_vpn_ipsec | :white_check_mark: |
| (67) | vcd_nsxt_nat_rule | cloudavenue_edgegateway_nat_rule | :white_check_mark: |
| (68) | vcd_nsxt_network_dhcp | cloudavenue_network_dhcp | :white_check_mark: |
| (69) | vcd_nsxt_network_dhcp_binding | cloudavenue_network_dhcp_binding | :white_check_mark: |
| (70) | vcd_nsxt_network_imported | Not Applicable | :heavy_multiplication_x: |
| (71) | vcd_nsxt_network_segment_profile | Not Applicable | :heavy_multiplication_x: |
| (72) | vcd_nsxt_route_advertisement | Not Applicable | :heavy_multiplication_x: |
| (73) | vcd_nsxt_security_group | cloudavenue_edgegateway_security_group | :white_check_mark: |
| (74) | vcd_nsxt_segment_profile_template | Not Applicable | :heavy_multiplication_x: |
| (75) | vcd_nsxv_dhcp_relay | Not Applicable | :heavy_multiplication_x: |
| (76) | vcd_nsxv_distributed_firewall | Not Applicable | :heavy_multiplication_x: |
| (77) | vcd_nsxv_dnat | Not Applicable | :heavy_multiplication_x: |
| (78) | vcd_nsxv_firewall_rule | Not Applicable | :heavy_multiplication_x: |
| (79) | vcd_nsxv_ip_set | Not Applicable | :heavy_multiplication_x: |
| (80) | vcd_nsxv_snat | Not Applicable | :heavy_multiplication_x: |
| (81) | vcd_org | Not Applicable | :heavy_multiplication_x: |
| (82) | vcd_org_group | Not Applicable | :heavy_multiplication_x: |
| (83) | vcd_org_ldap | Not Applicable | :heavy_multiplication_x: |
| (84) | vcd_org_oidc | Not yet implemented | :x: |
| (85) | vcd_org_saml | Not Applicable | :heavy_multiplication_x: |
| (86) | vcd_org_user | cloudavenue_iam_user | :white_check_mark: |
| (87) | vcd_org_vdc | Not Applicable | :heavy_multiplication_x: |
| (88) | vcd_org_vdc_access_control | cloudavenue_vdc_acl | :white_check_mark: |
| (89) | vcd_org_vdc_nsxt_network_profile | Not Applicable | :heavy_multiplication_x: |
| (90) | vcd_org_vdc_template | Not Applicable | :heavy_multiplication_x: |
| (91) | vcd_org_vdc_template_instance | Not Applicable | :heavy_multiplication_x: |
| (92) | vcd_provider_vdc | Not Applicable | :heavy_multiplication_x: |
| (93) | vcd_rde | Not Applicable | :heavy_multiplication_x: |
| (94) | vcd_rde_interface | Not Applicable | :heavy_multiplication_x: |
| (95) | vcd_rde_interface_behavior | Not Applicable | :heavy_multiplication_x: |
| (96) | vcd_rde_type | Not Applicable | :heavy_multiplication_x: |
| (97) | vcd_rde_type_behavior | Not Applicable | :heavy_multiplication_x: |
| (98) | vcd_rde_type_behavior_acl | Not Applicable | :heavy_multiplication_x: |
| (99) | vcd_rights_bundle | Not Applicable | :heavy_multiplication_x: |
| (100) | vcd_role | cloudavenue_iam_role | :white_check_mark: |
| (101) | vcd_security_tag | cloudavenue_vm_security_tag | :white_check_mark: |
| (102) | vcd_service_account | Not Applicable | :heavy_multiplication_x: |
| (103) | vcd_solution_add_on | Not Applicable | :heavy_multiplication_x: |
| (104) | vcd_solution_add_on_instance | Not Applicable | :heavy_multiplication_x: |
| (105) | vcd_solution_add_on_instance_publish | Not Applicable | :heavy_multiplication_x: |
| (106) | vcd_solution_landing_zone | Not Applicable | :heavy_multiplication_x: |
| (107) | vcd_subscribed_catalog | Not Applicable | :heavy_multiplication_x: |
| (108) | vcd_ui_plugin | Not Applicable | :heavy_multiplication_x: |
| (109) | vcd_vapp | cloudavenue_vapp |:white_check_mark: |
| (110) | vcd_vapp_access_control | cloudavenue_vapp_acl | :white_check_mark: |
| (111) | vcd_vapp_firewall_rules | Not Applicable | :heavy_multiplication_x: |
| (112) | vcd_vapp_nat_rules | Not Applicable | :heavy_multiplication_x: |
| (113) | vcd_vapp_network | cloudavenue_vapp_isolated_network | :white_check_mark: |
| (114) | vcd_vapp_org_network | cloudavenue_vapp_org_network |:white_check_mark: |
| (115) | vcd_vapp_static_routing | Not Applicable | :heavy_multiplication_x: |
| (116) | vcd_vapp_vm | cloudavenue_vm | :white_check_mark: |
| (117) | vcd_vdc_group | Not Applicable | :heavy_multiplication_x: |
| (118) | vcd_vm | Not Applicable | :heavy_multiplication_x: |
| (119) | vcd_vm_affinity_rule | cloudavenue_vm_affinity_rule |:white_check_mark: |
| (120) | vcd_vm_internal_disk | cloudavenue_vm_disk | :white_check_mark: |
| (121) | vcd_vm_placement_policy | Not Applicable | :heavy_multiplication_x: |
| (122) | vcd_vm_sizing_policy | Not Applicable | :heavy_multiplication_x: |
| (123) | vcd_vm_vgpu_policy | Not Applicable | :heavy_multiplication_x: |

| Number | Datasources VMware VCD | Datasources Orange Cloud Avenue | status |
|:--:|:--:|:--:|:--:|
| (1) | vcd_api_filter | Not Applicable | :heavy_multiplication_x: |
| (2) | vcd_catalog | cloudavenue_catalog |:white_check_mark: |
| (3) | vcd_catalog_access_control | cloudavenue_catalog_acl | :white_check_mark: |
| (4) | vcd_catalog_item | Not Applicable | :heavy_multiplication_x: |
| (5) | vcd_catalog_media | cloudavenue_catalog_media |:white_check_mark: |
| (6) | vcd_catalog_vapp_template | cloudavenue_catalog_vapp_template |:white_check_mark: |
| (7) | vcd_cse_kubernetes_cluster | Not Applicable | :heavy_multiplication_x: |
| (8) | vcd_dse_registry_configuration | Not Applicable | :heavy_multiplication_x: |
| (9) | vcd_dse_solution_publish | Not Applicable | :heavy_multiplication_x: |
| (10) | vcd_edgegateway | Not Applicable | :heavy_multiplication_x: |
| (11) | vcd_external_endpoint | Not Applicable | :heavy_multiplication_x: |
| (12) | vcd_external_network | Not Applicable | :heavy_multiplication_x: |
| (13) | vcd_external_network_v2 | Not Applicable | :heavy_multiplication_x: |
| (14) | vcd_global_role | Not Applicable | :heavy_multiplication_x: |
| (15) | vcd_independent_disk | cloudavenue_vm_disk | :white_check_mark: |
| (16) | vcd_ip_space | Not Applicable | :heavy_multiplication_x: |
| (17) | vcd_ip_space_custom_quota | Not Applicable | :heavy_multiplication_x: |
| (18) | vcd_ip_space_ip_allocation | Not Applicable | :heavy_multiplication_x: |
| (19) | vcd_ip_space_uplink | Not Applicable | :heavy_multiplication_x: |
| (20) | vcd_lb_app_profile | Not Applicable | :heavy_multiplication_x: |
| (21) | vcd_lb_app_rule | Not Applicable | :heavy_multiplication_x: |
| (22) | vcd_lb_server_pool | Not Applicable | :heavy_multiplication_x: |
| (23) | vcd_lb_service_monitor | Not Applicable | :heavy_multiplication_x: |
| (24) | vcd_lb_virtual_server | Not Applicable | :heavy_multiplication_x: |
| (25) | vcd_library_certificate | Not Applicable | :heavy_multiplication_x: |
| (26) | vcd_multisite_org_association | Not yet implemented | :x: |
| (27) | vcd_multisite_org_data | Not yet implemented | :x: |
| (28) | vcd_multisite_site | Not Applicable | :heavy_multiplication_x: |
| (29) | vcd_multisite_site_association | Not Applicable | :heavy_multiplication_x: |
| (30) | vcd_multisite_site_data | Not Applicable | :heavy_multiplication_x: |
| (31) | vcd_network_direct | Not Applicable | :heavy_multiplication_x: |
| (32) | vcd_network_isolated | Not Applicable | :heavy_multiplication_x: |
| (33) | vcd_network_isolated_v2 | cloudavenue_network_isolated | :white_check_mark: |
| (34) | vcd_network_pool | Not Applicable | :heavy_multiplication_x: |
| (35) | vcd_network_routed | Not Applicable | :heavy_multiplication_x: |
| (36) | vcd_network_routed_v2 | cloudavenue_network_routed | :white_check_mark: |
| (37) | vcd_nsxt_alb_cloud | Not Applicable | :heavy_multiplication_x: |
| (38) | vcd_nsxt_alb_controller | Not Applicable | :heavy_multiplication_x: |
| (39) | vcd_nsxt_alb_edgegateway_service_engine_group | Not Applicable | :heavy_multiplication_x: |
| (40) | vcd_nsxt_alb_importable_cloud | Not yet implemented | :x: |
| (41) | vcd_nsxt_alb_pool | cloudavenue_alb_pool | :white_check_mark: |
| (42) | vcd_nsxt_alb_service_engine_group | Not Applicable | :heavy_multiplication_x: |
| (43) | vcd_nsxt_alb_settings | Not Applicable | :heavy_multiplication_x: |
| (44) | vcd_nsxt_alb_virtual_service | Not Applicable | :heavy_multiplication_x: |
| (45) | vcd_nsxt_alb_virtual_service_http_req_rules | Not yet implemented | :x: |
| (46) | vcd_nsxt_alb_virtual_service_http_resp_rules | Not yet implemented | :x: |
| (47) | vcd_nsxt_alb_virtual_service_http_sec_rules | Not yet implemented | :x: |
| (48) | vcd_nsxt_app_port_profile | cloudavenue_edgegateway_app_port_profile | :white_check_mark: |
| (49) | vcd_nsxt_distributed_firewall | Not yet implemented | :x: |
| (50) | vcd_nsxt_distributed_firewall_rule | Not yet implemented | :x: |
| (51) | vcd_nsxt_dynamic_security_group | Not yet implemented | :x: |
| (52) | vcd_nsxt_edge_cluster | Not yet implemented | :x: |
| (53) | vcd_nsxt_edgegateway | cloudavenue_edgegateway | :white_check_mark: |
| (54) | vcd_nsxt_edgegateway_bgp_configuration | Not Applicable | :heavy_multiplication_x: |
| (55) | vcd_nsxt_edgegateway_bgp_ip_prefix_list | Not Applicable | :heavy_multiplication_x: |
| (56) | vcd_nsxt_edgegateway_bgp_neighbor | Not Applicable | :heavy_multiplication_x: |
| (57) | vcd_nsxt_edgegateway_dhcp_forwarding | cloudavenue_edgegateway_dhcp_forwarding | :white_check_mark: |
| (58) | vcd_nsxt_edgegateway_dhcpv6 | Not Applicable | :heavy_multiplication_x: |
| (59) | vcd_nsxt_edgegateway_dns | Not yet implemented | :x: |
| (60) | vcd_nsxt_edgegateway_l2_vpn_tunnel | Not yet implemented | :x: |
| (61) | vcd_nsxt_edgegateway_qos_profile | Not Applicable | :heavy_multiplication_x: |
| (62) | vcd_nsxt_edgegateway_rate_limiting | Not Applicable | :heavy_multiplication_x: |
| (63) | vcd_nsxt_edgegateway_static_route | cloudavenue_edgegateway_static_route | :white_check_mark: |
| (64) | vcd_nsxt_firewall | cloudavenue_edgegateway_firewall | :white_check_mark: |
| (65) | vcd_nsxt_global_default_segment_profile_template | Not Applicable | :heavy_multiplication_x: |
| (66) | vcd_nsxt_ip_set | cloudavenue_edgegateway_ip_set | :white_check_mark: |
| (67) | vcd_nsxt_ipsec_vpn_tunnel | cloudavenue_edgegateway_vpn_ipsec | :white_check_mark: |
| (68) | vcd_nsxt_manager | Not Applicable | :heavy_multiplication_x: |
| (69) | vcd_nsxt_nat_rule | cloudavenue_edgegateway_nat_rule | :white_check_mark: |
| (70) | vcd_nsxt_network_context_profile | Not Applicable | :heavy_multiplication_x: |
| (71) | vcd_nsxt_network_dhcp | cloudavenue_network_dhcp | :white_check_mark: |
| (72) | vcd_nsxt_network_dhcp_binding | cloudavenue_network_dhcp_binding | :white_check_mark: |
| (73) | vcd_nsxt_network_imported | Not Applicable | :heavy_multiplication_x: |
| (74) | vcd_nsxt_network_segment_profile | Not Applicable | :heavy_multiplication_x: |
| (75) | vcd_nsxt_route_advertisement | Not Applicable | :heavy_multiplication_x: |
| (76) | vcd_nsxt_security_group | cloudavenue_edgegateway_security_group | :white_check_mark: |
| (77) | vcd_nsxt_segment_ip_discovery_profile | Not Applicable | :heavy_multiplication_x: |
| (78) | vcd_nsxt_segment_mac_discovery_profile | Not Applicable | :heavy_multiplication_x: |
| (79) | vcd_nsxt_segment_profile_template | Not Applicable | :heavy_multiplication_x: |
| (80) | vcd_nsxt_segment_qos_profile | Not Applicable | :heavy_multiplication_x: |
| (81) | vcd_nsxt_segment_security_profile | Not Applicable | :heavy_multiplication_x: |
| (82) | vcd_nsxt_segment_spoof_guard_profile | Not Applicable | :heavy_multiplication_x: |
| (83) | vcd_nsxt_tier0_router | Not Applicable | :heavy_multiplication_x: |
| (84) | vcd_nsxt_tier0_router_interface | Not Applicable | :heavy_multiplication_x: |
| (85) | vcd_nsxv_application | Not Applicable | :heavy_multiplication_x: |
| (86) | vcd_nsxv_application_finder | Not Applicable | :heavy_multiplication_x: |
| (87) | vcd_nsxv_application_group | Not Applicable | :heavy_multiplication_x: |
| (88) | vcd_nsxv_dhcp_relay | Not Applicable | :heavy_multiplication_x: |
| (89) | vcd_nsxv_distributed_firewall | Not Applicable | :heavy_multiplication_x: |
| (90) | vcd_nsxv_dnat | Not Applicable | :heavy_multiplication_x: |
| (91) | vcd_nsxv_firewall_rule | Not Applicable | :heavy_multiplication_x: |
| (92) | vcd_nsxv_ip_set | Not Applicable | :heavy_multiplication_x: |
| (93) | vcd_nsxv_snat | Not Applicable | :heavy_multiplication_x: |
| (94) | vcd_org | Not Applicable | :heavy_multiplication_x: |
| (95) | vcd_org_group | Not Applicable | :heavy_multiplication_x: |
| (96) | vcd_org_ldap | Not Applicable | :heavy_multiplication_x: |
| (97) | vcd_org_oidc | Not yet implemented | :x: |
| (98) | vcd_org_saml | Not Applicable | :heavy_multiplication_x: |
| (99) | vcd_org_saml_metadata | Not yet implemented | :x: |
| (100) | vcd_org_user | cloudavenue_iam_user | :white_check_mark: |
| (101) | vcd_org_vdc | Not Applicable | :heavy_multiplication_x: |
| (102) | vcd_org_vdc_nsxt_network_profile | Not Applicable | :heavy_multiplication_x: |
| (103) | vcd_org_vdc_template | Not Applicable | :heavy_multiplication_x: |
| (104) | vcd_portgroup | Not Applicable | :heavy_multiplication_x: |
| (105) | vcd_provider_vdc | Not Applicable | :heavy_multiplication_x: |
| (106) | vcd_rde | Not Applicable | :heavy_multiplication_x: |
| (107) | vcd_rde_behavior_invocation | Not Applicable | :heavy_multiplication_x: |
| (108) | vcd_rde_interface | Not Applicable | :heavy_multiplication_x: |
| (109) | vcd_rde_interface_behavior | Not Applicable | :heavy_multiplication_x: |
| (110) | vcd_rde_type | Not Applicable | :heavy_multiplication_x: |
| (111) | vcd_rde_type_behavior | Not Applicable | :heavy_multiplication_x: |
| (112) | vcd_rde_type_behavior_acl | Not Applicable | :heavy_multiplication_x: |
| (113) | vcd_resource_list | Not Applicable | :heavy_multiplication_x: |
| (114) | vcd_resource_pool | Not Applicable | :heavy_multiplication_x: |
| (115) | vcd_resource_schema | Not Applicable | :heavy_multiplication_x: |
| (116) | vcd_right | Not yet implemented | :x: |
| (117) | vcd_rights_bundle | Not Applicable | :heavy_multiplication_x: |
| (118) | vcd_role | cloudavenue_iam_role | :white_check_mark: |
| (119) | vcd_service_account | Not Applicable | :heavy_multiplication_x: |
| (120) | vcd_solution_add_on | Not Applicable | :heavy_multiplication_x: |
| (121) | vcd_solution_add_on_instance | Not Applicable | :heavy_multiplication_x: |
| (122) | vcd_solution_add_on_instance_publish | Not Applicable | :heavy_multiplication_x: |
| (123) | vcd_solution_landing_zone | Not Applicable | :heavy_multiplication_x: |
| (124) | vcd_storage_profile | cloudavenue_storage_profile |:white_check_mark: |
| (125) | vcd_subscribed_catalog | Not Applicable | :heavy_multiplication_x: |
| (126) | vcd_task | Not Applicable | :heavy_multiplication_x: |
| (127) | vcd_ui_plugin | Not Applicable | :heavy_multiplication_x: |
| (128) | vcd_vapp | cloudavenue_vapp |:white_check_mark: |
| (129) | vcd_vapp_network | cloudavenue_vapp_isolated_network | :white_check_mark: |
| (130) | vcd_vapp_org_network | cloudavenue_vapp_org_network |:white_check_mark: |
| (131) | vcd_vapp_vm | cloudavenue_vm | :white_check_mark: |
| (132) | vcd_vcenter | Not Applicable | :heavy_multiplication_x: |
| (133) | vcd_vdc_group | Not Applicable | :heavy_multiplication_x: |
| (134) | vcd_version | Not Applicable | :heavy_multiplication_x: |
| (135) | vcd_vgpu_profile | Not Applicable | :heavy_multiplication_x: |
| (136) | vcd_vm | Not Applicable | :heavy_multiplication_x: |
| (137) | vcd_vm_affinity_rule | cloudavenue_vm_affinity_rule |:white_check_mark: |
| (138) | vcd_vm_group | Not Applicable | :heavy_multiplication_x: |
| (139) | vcd_vm_placement_policy | Not Applicable | :heavy_multiplication_x: |
| (140) | vcd_vm_sizing_policy | Not Applicable | :heavy_multiplication_x: |
| (141) | vcd_vm_vgpu_policy | Not Applicable | :heavy_multiplication_x: |
