# If `vDC` is not specified, the default `vDC` will be used
# The `myVMID` is the ID of the VM. 
terraform import cloudavenue_vm.example myVAPP.myVMID

# or you can specify the vDC
terraform import cloudavenue_vm.example myVDC.myVAPP.myVMID
