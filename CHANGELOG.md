## 0.6.0 (Unreleased)
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
