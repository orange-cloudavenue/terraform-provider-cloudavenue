
# if vdc is not specified, the default vdc will be used
terraform import cloudavenue_vm_disk.example vapp_name.vm_name.id

# if vdc is specified, the vdc will be used
terraform import cloudavenue_vm_disk.example vdc.vapp_name.vm_name.id
