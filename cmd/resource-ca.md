# Checking resources and datasources of Orange Cloud Avenue provider
- Found 22 resources in terraform
- Found 27 datasources in terraform

# Checking resources and datasources of VMware Cloud Director provider
- Found 81 resources in terraform
- Found 88 datasources in terraform


# Listing cross resources and datasources from Cloud Avenue
| Number |Resources Orange Cloud Avenue | Resources VMware VCD |
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
| (20) | cloudavenue_vm_disk | vcd_independent_disk |
| (21) | cloudavenue_vm_inserted_media | vcd_inserted_media |
| (22) | cloudavenue_vm_security_tag | vcd_security_tag |

# Listing cross resources and datasources from VCD
| Number | Resources VMware VCD | Resources Orange Cloud Avenue |
|:--:|:--:|:--:|
| (1) | vcd_edgegateway_vpn | Not Applicable |
| (2) | vcd_network_isolated | Not Applicable |
| (3) | vcd_network_routed | Not Applicable |
| (4) | vcd_org_group | Not yet implemented |
| (5) | vcd_rde | Not yet implemented |
| (6) | vcd_vm_sizing_policy | Not Applicable |
| (7) | vcd_nsxt_dynamic_security_group | Not yet implemented |
| (8) | vcd_external_network_v2 | Not Applicable |
| (9) | vcd_network_isolated_v2 | cloudavenue_network_isolated |
| (10) | vcd_org_vdc | Not Applicable |
| (11) | vcd_vapp_nat_rules | Not yet implemented |
| (12) | vcd_vapp_static_routing | Not yet implemented |
| (13) | vcd_catalog_vapp_template | Not yet implemented |
| (14) | vcd_nsxt_alb_controller | Not Applicable |
| (15) | vcd_lb_server_pool | Not Applicable |
| (16) | vcd_nsxt_alb_settings | Not Applicable |
| (17) | vcd_nsxt_firewall | cloudavenue_edgegateway_firewall |
| (18) | vcd_rights_bundle | Not Applicable |
| (19) | vcd_vdc_group | Not Applicable |
| (20) | vcd_vm_internal_disk | cloudavenue_vm_disk |
| (21) | vcd_catalog_media | Not yet implemented |
| (22) | vcd_nsxt_edgegateway | cloudavenue_edgegateway |
| (23) | vcd_nsxt_edgegateway_bgp_configuration | Not Applicable |
| (24) | vcd_vapp_access_control | cloudavenue_vapp_acl |
| (25) | vcd_global_role | Not Applicable |
| (26) | vcd_nsxt_alb_cloud | Not Applicable |
| (27) | vcd_nsxt_ip_set | Not yet implemented |
| (28) | vcd_nsxv_ip_set | Not Applicable |
| (29) | vcd_inserted_media | cloudavenue_vm_inserted_media |
| (30) | vcd_network_routed_v2 | cloudavenue_network_routed |
| (31) | vcd_nsxt_alb_edgegateway_service_engine_group | Not Applicable |
| (32) | vcd_nsxt_ipsec_vpn_tunnel | Not yet implemented |
| (33) | vcd_nsxt_network_dhcp | Not yet implemented |
| (34) | vcd_catalog_access_control | Not yet implemented |
| (35) | vcd_nsxt_network_imported | Not Applicable |
| (36) | vcd_nsxv_firewall_rule | Not Applicable |
| (37) | vcd_rde_interface | Not yet implemented |
| (38) | vcd_nsxv_dnat | Not Applicable |
| (39) | vcd_rde_type | Not yet implemented |
| (40) | vcd_lb_app_rule | Not Applicable |
| (41) | vcd_lb_virtual_server | Not Applicable |
| (42) | vcd_nsxt_edgegateway_bgp_ip_prefix_list | Not Applicable |
| (43) | vcd_nsxt_edgegateway_rate_limiting | Not yet implemented |
| (44) | vcd_nsxt_security_group | Not yet implemented |
| (45) | vcd_nsxv_distributed_firewall | Not Applicable |
| (46) | vcd_vm_placement_policy | Not Applicable |
| (47) | vcd_org_user | cloudavenue_iam_user |
| (48) | vcd_org_vdc_access_control | cloudavenue_vdc_acl |
| (49) | vcd_lb_app_profile | Not Applicable |
| (50) | vcd_network_direct | Not Applicable |
| (51) | vcd_nsxt_alb_pool | cloudavenue_alb_pool |
| (52) | vcd_nsxt_distributed_firewall | Not yet implemented |
| (53) | vcd_nsxt_nat_rule | Not yet implemented |
| (54) | vcd_nsxt_route_advertisement | Not Applicable |
| (55) | vcd_vm_affinity_rule | cloudavenue_vm_affinity_rule |
| (56) | vcd_catalog | cloudavenue_catalog |
| (57) | vcd_library_certificate | Not yet implemented |
| (58) | vcd_nsxt_alb_virtual_service | Not Applicable |
| (59) | vcd_vapp_firewall_rules | Not yet implemented |
| (60) | vcd_catalog_item | Not Applicable |
| (61) | vcd_edgegateway_settings | Not Applicable |
| (62) | vcd_nsxv_snat | Not Applicable |
| (63) | vcd_org_ldap | Not Applicable |
| (64) | vcd_subscribed_catalog | Not Applicable |
| (65) | vcd_vapp | cloudavenue_vapp |
| (66) | vcd_nsxt_alb_service_engine_group | Not Applicable |
| (67) | vcd_nsxt_app_port_profile | cloudavenue_edgegateway_app_port_profile |
| (68) | vcd_role | cloudavenue_iam_role |
| (69) | vcd_vapp_org_network | cloudavenue_vapp_org_network |
| (70) | vcd_vapp_vm | cloudavenue_vm |
| (71) | vcd_vm | Not Applicable |
| (72) | vcd_independent_disk | cloudavenue_vm_disk |
| (73) | vcd_nsxt_edgegateway_bgp_neighbor | Not Applicable |
| (74) | vcd_nsxt_network_dhcp_binding | Not yet implemented |
| (75) | vcd_nsxv_dhcp_relay | Not Applicable |
| (76) | vcd_vapp_network | cloudavenue_vapp_isolated_network |
| (77) | vcd_edgegateway | Not Applicable |
| (78) | vcd_external_network | Not Applicable |
| (79) | vcd_lb_service_monitor | Not Applicable |
| (80) | vcd_org | Not Applicable |
| (81) | vcd_security_tag | cloudavenue_vm_security_tag |
