# if vdc is not specified, the default vdc will be used
terraform import cloudavenue_vapp_isolated_network.example vapp_name.network_name

# if vdc is specified, the vdc will be used
terraform import cloudavenue_vapp_isolated_network.example vdc.vapp_name.network_name