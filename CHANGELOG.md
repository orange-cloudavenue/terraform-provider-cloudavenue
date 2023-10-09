## 0.13.0 (Unreleased)

### :dependabot: **Dependencies**

* deps: bumps stefanzweifel/git-auto-commit-action from 4 to 5 ([GH-559](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/559))

## 0.12.0 (October  2, 2023)
### :rotating_light: **Breaking Changes**

* `datasource/cloudavenue_vdc_group` - Remove attributes `local_egress`, `error_message`, `dfw_enabled`, `network_pool_id`, `network_pool_universal_id`, `network_provider_type`, `universal_networking_enabled`, `vdcs`, `fault_domain_tag`, `is_remote_org`, `name`, `network_provider_scope`, `site_id`, `site_name` from the datasource. 
The attribute `vdc_ids` is added to the datasource and return the list of VDC IDs of the VDC Group. ([GH-442](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/442))
* `datasource/cloudavenue_vdc` - The `vdc_group` attribute is now **removed**. ([GH-447](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/447))
* `resource/cloudavenue_vdc` - Announced in the release [v0.0.9](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/releases/tag/v0.9.0) the attribute `vdc_group` is now **removed**. ([GH-447](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/447))

### :rocket: **New Features**

* `client/cloudavenue` - Add `NetBackup` credentials in provider configuration. ([GH-546](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/546))
### :information_source: **Notes**

* `provider` - Improve documentation for provider authentication and configuration. ([GH-548](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/548))
* `resource/cloudavenue_vdc_group` - Add import documentation. ([GH-442](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/442))
* `resource/cloudavenue_vdc` - Fix values in documentation for attributes `cpu_allocated`, `memory_allocated`, `cpu_speed_in_mhz`. ([GH-530](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/530))
* `resource/cloudavenue_vdc` - Improve documentation. ([GH-533](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/533))
* `resource/cloudavenue_vm_disk` - Improve documentation. ([GH-533](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/533))

### :dependabot: **Dependencies**

* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-superschema from 1.5.3 to 1.5.4 ([GH-533](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/533))
* deps: bumps github.com/rs/zerolog from 1.30.0 to 1.31.0 ([GH-541](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/541))

## 0.11.0 (September 22, 2023)
### :rotating_light: **Breaking Changes**

* `datasource/cloudavenue_vapp_org_network` - `is_fenced` and `retain_ip_mac_enabled` are now removed from the schema. ([GH-538](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/538))
* `resource/cloudavenue_network_routed` - Change attribute field for import resource by EdgeGatewayName or EdgeGatewayID. ([GH-526](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/526))
* `resource/cloudavenue_vapp_org_network` - `is_fenced` and `retain_ip_mac_enabled` are now removed from the schema. ([GH-538](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/538))
### :information_source: **Notes**

* `resource/cloudavenue_vdc` - Only attributs `cpu_allocated`, `memory_allocated`, `storage_profile`, `cpu_speed_in_mhz` and description can be modified. ForceNew for other. ([GH-524](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/524))

## 0.10.4 (September 15, 2023)

### :bug: **Bug Fixes**

* `resource/cloudavenue_vdc` - Fix bug in vdc about vdcgroup field or vdc_group resource. ([GH-521](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/521))

## 0.10.3 (September 14, 2023)

### :bug: **Bug Fixes**

* `resource/cloudavenue_vcd` - Fix bug to impossible to update a vcd resource without a resource vcd_group define. ([GH-518](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/518))

### :dependabot: **Dependencies**

* deps: bumps crazy-max/ghaction-import-gpg from 5.4.0 to 6.0.0 ([GH-516](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/516))
* deps: bumps github.com/hashicorp/terraform-plugin-framework from 1.3.5 to 1.4.0 ([GH-512](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/512))
* deps: bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.28.0 to 2.29.0 ([GH-513](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/513))
* deps: bumps github.com/orange-cloudavenue/cloudavenue-sdk-go from 0.1.2 to 0.1.3 ([GH-518](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/518))
* deps: bumps goreleaser/goreleaser-action from 4.6.0 to 5.0.0 ([GH-517](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/517))

## 0.10.2 (September  8, 2023)

### :tada: **Improvements**

* `resource/cloudavenue_alb_pool` - Improve documentation. ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))
* `resource/cloudavenue_catalog_acl` - Improve documentation. ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))
* `resource/cloudavenue_edgegateway_nat_rule` - Improve documentation. ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))
* `resource/cloudavenue_edgegateway_vpn_ipsec` - Improve documentation. ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))
* `resource/cloudavenue_network_dhcp` - Improve documentation. ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))
* `resource/cloudavenue_vm_disk` - Improve documentation. ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))
* `resource/cloudavenue_vm` - Improve documentation. ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))

### :bug: **Bug Fixes**

* `resource/cloudavenue_vcd` - Fix bug to impossible to update a vcd resource with a resource vcd_group define. ([GH-511](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/511))

### :dependabot: **Dependencies**

* deps: bumps actions/checkout from 3 to 4 ([GH-508](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/508))
* deps: bumps crazy-max/ghaction-import-gpg from 5.3.0 to 5.4.0 ([GH-507](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/507))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers from 1.3.2 to 1.3.3 ([GH-505](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/505))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-superchema from v1.5.2 to v1.5.3 ([GH-500](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/500))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-superschema from v1.5.1 to v1.5.2 ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-supertypes from v0.0.5 to v0.1.0 ([GH-500](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/500))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-validators from v1.7.0 to v1.8.0 ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))
* deps: bumps github.com/google/uuid from 1.3.0 to 1.3.1 ([GH-499](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/499))
* deps: bumps github.com/hashicorp/terraform-plugin-framework from 1.3.4 to 1.3.5 ([GH-497](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/497))
* deps: bumps github.com/hashicorp/terraform-plugin-framework-validators from 0.10.0 to 0.11.0 ([GH-476](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/476))
* deps: bumps github.com/hashicorp/terraform-plugin-framework-validators from 0.11.0 to 0.12.0 ([GH-506](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/506))
* deps: bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.27.0 to 2.28.0 ([GH-502](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/502))
* deps: bumps golang.org/x/net from 0.14.0 to 0.15.0 ([GH-510](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/510))
* deps: bumps golangci/golangci-lint-action from 3.6.0 to 3.7.0 ([GH-494](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/494))
* deps: bumps goreleaser/goreleaser-action from 4.4.0 to 4.6.0 ([GH-509](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/509))

## 0.10.0 (August 11, 2023)

### :rocket: **New Features**

* **New Data Source:** `datasource/cloudavenue_edgegateway_vpn_ipsec` - New data source to read Cloud Avenue IPsec VPN Tunnel. ([GH-353](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/353))
* **New Data Source:** `datasource/cloudavenue_vm_disks` - New datasource to get the list of disks available on vApp and VM level. ([GH-475](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/475))
* **New Resource:** `resource/cloudavenue_edgegateway_vpn_ipsec` - New resource to manage Cloud Avenue IPSec VPN Tunnel. ([GH-352](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/352))
* **New Resource:** `resource/cloudavenue_iam_token` - New resource to create user token. ([GH-423](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/423))

### :tada: **Improvements**

* `resource/cloudavenue_vdc` - Improve example in documentation. ([GH-481](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/481))
* `resource/cloudavenue_vm_disk` - Add import feature and improve documentation. ([GH-344](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/344))
* `resource/cloudavenue_vm_disk` - Now the attributes `vm_id` and `vm_name` causes replacement if `is_detachable` is set to `false`. ([GH-480](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/480))

### :bug: **Bug Fixes**

* `resource/cloudavenue_vm_disk` - Fix bug not possible to create detachable disk with no `vm_id` or `vm_name` specified. ([GH-478](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/478))

### :dependabot: **Dependencies**

* deps: bumps golang.org/x/net from 0.13.0 to 0.14.0 ([GH-477](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/477))
* deps: bumps goreleaser/goreleaser-action from 4.3.0 to 4.4.0 ([GH-487](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/487))

## 0.9.0 (August  7, 2023)
### :warning: **Deprecations**

* `resource/cloudavenue_vdc` - The `vdc_group` attribute has been deprecated and will be removed in a `v0.12.0` release. Please use `cloudavenue_vdc_group` resource instead. ([GH-448](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/448))

### :rocket: **New Features**

* **New Data Source:** `datasource/cloudavenue_catalog_acl` - New data source to get the ACL of a catalog ([GH-472](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/472))
* **New Resource:** `resource/cloudavenue_catalog_acl` - New resource to manage Cloud Avenue Catalog ACLs ([GH-453](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/453))
* **New Resource:** `resource/cloudavenue_vdc_group` - Add new resource to manage VDC Group in Cloud Avenue. ([GH-445](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/445))

### :tada: **Improvements**

* `datasource/cloudavenue_catalog_vapp_template` - Now the `template_name` and `template_id` attributes are always returned and improve the documentation. ([GH-443](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/443))
* `datasource/cloudavenue_catalogs` - Improve the documentation. ([GH-443](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/443))
* `resource/catalog` - Improve example in documentation
`datasource/edgegateway_nat_rule` - Improve example in documentation
`datasource/vapp_isolated_network` - Improve example in documentation
`datasource/vm` - Improve example in documentation ([GH-451](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/451))
### :information_source: **Notes**

* `datasource/cloudavenue_catalog_vapp_template` - Now the `id` attribute return URN instead of UUID. ([GH-443](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/443))

### :dependabot: **Dependencies**

* deps: bumps github.com/hashicorp/terraform-plugin-framework from 1.3.3 to 1.3.4 ([GH-461](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/461))
* deps: bumps github.com/hashicorp/terraform-plugin-superschema from 1.4.1 to 1.5.1 ([GH-448](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/448))
* deps: bumps golang.org/x/net from 0.12.0 to 0.13.0 ([GH-446](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/446))

## 0.8.0 (August  1, 2023)

### :rocket: **New Features**

* **New Data Source:** `datasource/cloudavenue_edgegateway_dhcp_forwarding` - New data source to get DHCP forwarding configuration from Edge Gateway. ([GH-422](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/422))
* **New Data Source:** `datasource/cloudavenue_edgegateway_nat_rule` - New datasource to get a NAT Rule in edge gateway. ([GH-356](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/356))
* **New Data Source:** `datasource/cloudavenue_edgegateway_static_route` - New data source to fetch static route details from Cloud Avenue Edge Gateway. ([GH-428](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/428))
* **New Resource:** `resource/cloudavenue_edgegateway_dhcp_forwarding` - New resource to manage DHCP forwarding on Edge Gateway. ([GH-421](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/421))
* **New Resource:** `resource/cloudavenue_edgegateway_nat_rule` - New resource to manage a NAT Rule in Edge Gateway. ([GH-355](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/355))
* **New Resource:** `resource/cloudavenue_edgegateway_static_route` - New resource to manage static routes on edge gateway. ([GH-427](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/427))

### :tada: **Improvements**

* `resource/cloudavenue_edgegateway_ip_set` - Improve example in documentation.
`datasource/cloudavenue_edgegateway_ip_set` - Improve example in documentation. ([GH-431](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/431))
* `resource/cloudavenue_vapp_acl` - Now the attribute `access_level` support `ReadOnly`, `Change` and `FullControl` options. The documentation has been updated accordingly. ([GH-407](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/407))

### :dependabot: **Dependencies**

* deps: bumps github.com/rs/zerolog from 1.29.1 to 1.30.0 ([GH-438](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/438))
* deps: bumps github.com/vmware/go-vcloud-director/v2 from 2.20.0 to 2.21.0 ([GH-403](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/403))

## 0.7.0 (July 25, 2023)

### :rocket: **New Features**

* **New Data Source:** `datasource/cloudavenue_edgegateway_ip_set` - New datasource to get the IP set of an edge gateway. ([GH-354](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/354))
* **New Data Source:** `datasource/cloudavenue_network_dhcp_binding` - New data source to get DHCP binding information from Org Network. ([GH-358](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/358))
* **New Resource:** `resource/cloudavenue_edgegateway_ip_set` - New resource to manage Edge Gateway IP Sets. ([GH-350](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/350))
* **New Resource:** `resource/cloudavenue_network_dhcp_binding` - New resource to manage DHCP bindings. ([GH-357](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/357))

### :tada: **Improvements**

* `datasource/cloudavenue_catalog_medias` - Improve documentation. ([GH-384](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/384))

### :dependabot: **Dependencies**

* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-superschema from 1.3.3 to 1.4.1 ([GH-397](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/397))
* deps: bumps github.com/hashicorp/terraform-plugin-framework from 1.3.2 to 1.3.3 ([GH-402](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/402))
* deps: bumps golang.org/x/net from 0.11.0 to 0.12.0 ([GH-409](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/409))

## 0.6.1 (July 18, 2023)

### :bug: **Bug Fixes**

* `resource/cloudavenue_publicip` - Fix bug in `public_ip` attribute. Now it is possible to set multiple publicip with the good attribute ip. ([GH-389](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/389))

## 0.6.0 (July 13, 2023)

### :rocket: **New Features**

* **New Data Source:** `datasource/cloudavenue_edgegateway_security_group` - New data source to fetch security group details from Cloud Avenue Edge Gateway. ([GH-351](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/351))
* **New Data Source:** `datasource/cloudavenue_network_dhcp` - New data source to get DHCP information for an organization network. ([GH-349](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/349))
* **New Resource:** `resource/cloudavenue_edgegateway_security_group` - New resource to manage Edge Gateway Security Group in Cloud Avenue. ([GH-342](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/342))
* **New Resource:** `resource/cloudavenue_network_dhcp` - New resource to manage DHCP for a organization network. ([GH-348](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/348))

### :tada: **Improvements**

* `resource/cloudavenue_vm_affinity_rule` - Add notice in documentation about polarity attribute. ([GH-380](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/380))

### :bug: **Bug Fixes**

* `resource/cloudavenue_vm_affinity_rule` - Fix bug in `vm_ids` attribute. Now it is possible to set more than 2 VMs IDs. ([GH-380](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/380))

### :dependabot: **Dependencies**

* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers from 1.3.1 to 1.3.2 ([GH-369](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/369))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-superschema from 1.3.2 to 1.3.3 ([GH-377](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/377))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-validators from 1.6.4 to 1.7.0 ([GH-379](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/379))
* deps: bumps github.com/hashicorp/terraform-plugin-docs from 0.15.0 to 0.16.0 ([GH-370](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/370))
* deps: bumps github.com/hashicorp/terraform-plugin-framework-timeouts from 0.4.0 to 0.4.1 ([GH-376](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/376))

## 0.5.2 (July  4, 2023)

### :bug: **Bug Fixes**

* `resource/cloudavenue_vapp_org_network` - Fixed a bug where failed to delete resource if the vapp status is RESOLVED. ([GH-365](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/365))

### :dependabot: **Dependencies**

* deps: bumps dependabot/fetch-metadata from 1.5.1 to 1.6.0 ([GH-359](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/359))
* deps: bumps github.com/hashicorp/terraform-plugin-framework from 1.3.1 to 1.3.2 ([GH-363](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/363))
* deps: bumps github.com/hashicorp/terraform-plugin-go from 0.16.0 to 0.17.0 ([GH-362](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/362))
* deps: bumps github.com/hashicorp/terraform-plugin-go from 0.17.0 to 0.18.0 ([GH-366](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/366))
* deps: bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.26.1 to 2.27.0 ([GH-361](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/361))

## 0.5.1 (June 27, 2023)

### :information_source: **Notes**

* `resource/cloudavenue_edgegateway_app_port_profile` is a resource moved from `cloudavenue_network_app_port_profile`. ([GH-347](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/347))

## 0.5.0 (June 27, 2023)

### :rocket: **New Features**

* **New Resource:** `resource/cloudavenue_edgegateway_firewall` - New resource to create a Edge Gateway firewall. ([GH-340](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/340))
* **New Resource:** `resource/cloudavenue_network_app_port_profile` - Is a new resource type that allows you to create a port profile for a network application. ([GH-319](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/319))
* **New Resource:** `resource/cloudavenue_publicip` is a new resource that can be used to manage public IP addresses in Cloud Avenue. ([GH-336](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/336))
* **New Resource:** `resource/cloudavenue_vcda_ip` - Is a new resource allows you to declare or remove your on-premises IP address for the DRaaS service ([GH-335](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/335))
* **New Datasource:** `datasource/cloudavenue_vm` - Data source allows you to read information about a virtual machine. ([GH-322](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/322))

### :tada: **Improvements**

* `resource/cloudavenue_vm_disk` - Add update support for disk not detachable. ([GH-332](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/332))
* `resource/cloudavenue_vm` - Add in the documentation the fields that require a VM restart if they are modified. ([GH-308](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/308))

### :dependabot: **Dependencies**

* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers from 1.3.0 to 1.3.1 ([GH-343](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/343))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-superschema from 1.3.1 to 1.3.2 ([GH-338](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/338))
* deps: bumps github.com/hashicorp/terraform-plugin-framework-timeouts from 0.3.1 to 0.4.0 ([GH-324](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/324))

## 0.4.0 (June 21, 2023)

NOTES:

* `cloudavenue_vm` - The attributes `settings.customization.force`, `settings.customization.change_sid`, `settings.customization.allow_local_admin_password`, `settings.customization.must_change_password_on_first_login`, `settings.customization.join_domain` and `settings.customization.join_org_domain` have now a default value of `false`. ([GH-320](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/320))

BREAKING CHANGES:

* `cloudavenue_vm` - Now the attributes `settings.customization.auto_generate_password` and `settings.customization.admin_password` are mutually exclusive and are no longer exactly one of. ([GH-320](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/320))
* `cloudavenue_vm` - The default value for attribute `deploy_os.accept_all_eulas` has been removed. ([GH-320](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/320))

FEATURES:

* `cloudavenue_vm` - Add import of VM. ([GH-320](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/320))
* **New Resource:** `cloudavenue_vm_security_tag` resource is added to manage security tags on VMs. ([GH-294](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/294))

BUG FIXES:

* `cloudavenue_vm` - Fix bugs in `settings.customization` and fix the ability to perform actions on multiple VMs simultaneously. ([GH-320](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/320))

DEPENDENCIES:

* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers from 1.2.2 to 1.3.0 ([GH-317](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/317))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-validators from 1.6.3 to 1.6.4 ([GH-315](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/315))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-superschema from 1.3.0 to 1.3.1 ([GH-316](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/316))

## 0.3.0 (June 20, 2023)

BREAKING CHANGES:

* Deletion of `power_on` attribute in schema for `cloudavenue_vapp` resource and datasource. ([GH-286](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/286))
* New major change on schema `cloudavenue_vm_disk` resource. ([GH-286](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/286))
* New major change on schema for `cloudavenue_vm` resource. ([GH-286](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/286))

FEATURES:

* **New Data Source:** cloudavenue_network_routed ([GH-249](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/249))
* **New Data Source:** cloudavenue_vapp_isolated_network ([GH-291](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/291))

BUG FIXES:

* Force to `power_off` a `cloudavenue_vapp` when you delete a `cloudavenue_vapp_org_network` resource. ([GH-286](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/286))

DEPENDENCIES:

* Update VMware Cloud Director API from v37.0 to v37.1 ([GH-286](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/286))
* deps: bumps actions/setup-go from 4.0.0 to 4.0.1 ([GH-314](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/314))
* deps: bumps crazy-max/ghaction-import-gpg from 5.2.0 to 5.3.0 ([GH-293](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/293))
* deps: bumps dependabot/fetch-metadata from 1.3.6 to 1.4.0 ([GH-289](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/289))
* deps: bumps dependabot/fetch-metadata from 1.4.0 to 1.5.0 ([GH-296](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/296))
* deps: bumps dependabot/fetch-metadata from 1.5.0 to 1.5.1 ([GH-297](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/297))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-validators from 1.4.0 to 1.5.0 ([GH-287](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/287))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-validators from 1.5.0 to 1.5.1 ([GH-290](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/290))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-validators from 1.5.1 to 1.6.3 ([GH-286](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/286))
* deps: bumps github.com/cloudflare/circl from 1.3.2 to 1.3.3 ([GH-295](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/295))
* deps: bumps github.com/hashicorp/terraform-plugin-docs from 0.14.1 to 0.15.0 ([GH-300](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/300))
* deps: bumps github.com/hashicorp/terraform-plugin-framework from 1.2.0 to 1.3.0 ([GH-301](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/301))
* deps: bumps github.com/hashicorp/terraform-plugin-framework from 1.3.0 to 1.3.1 ([GH-306](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/306))
* deps: bumps github.com/hashicorp/terraform-plugin-go from 0.15.0 to 0.16.0 ([GH-313](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/313))
* deps: bumps github.com/hashicorp/terraform-plugin-log from 0.8.0 to 0.9.0 ([GH-298](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/298))
* deps: bumps github.com/vmware/go-vcloud-director/v2 from 2.19.0 to 2.20.0 ([GH-292](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/292))
* deps: bumps golangci/golangci-lint-action from 3.4.0 to 3.5.0 ([GH-299](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/299))
* deps: bumps golangci/golangci-lint-action from 3.5.0 to 3.6.0 ([GH-307](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/307))
* deps: bumps goreleaser/goreleaser-action from 4.2.0 to 4.3.0 ([GH-304](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/304))

## 0.2.0 (April 7, 2023)

FEATURES:

* **New Data Source:** cloudavenue_alb_pool ([GH-246](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/246))
* **New Data Source:** cloudavenue_network_isolated ([GH-248](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/248))
* **New Resource:** cloudavenue_alb_pool ([GH-246](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/246))

DEPENDENCIES:

* deps: bumps actions/checkout from 2 to 3 ([GH-263](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/263))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-superschema from 1.2.0 to 1.3.0 ([GH-265](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/265))
* deps: bumps github.com/FrangipaneTeam/terraform-plugin-framework-validators from 1.3.1 to 1.4.0 ([GH-271](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/271))

## 0.1.1 (April 3, 2023)

IMPROVEMENTS:

* Docs: Fix categories

## 0.1.0 (April 3, 2023)

Initial release
