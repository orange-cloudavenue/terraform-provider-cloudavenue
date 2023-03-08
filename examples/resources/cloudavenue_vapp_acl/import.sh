# if vdc is not specified, the default vdc will be used
terraform import cloudavenue_vapp_acl.example vapp_name

# if vdc is specified, the vdc will be used
terraform import cloudavenue_vapp_acl.example vdc_name.vapp_name