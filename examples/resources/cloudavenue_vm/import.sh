# SPDX-FileCopyrightText: Copyright (c) 2025 Orange
# SPDX-License-Identifier: Mozilla Public License 2.0
#
# This software is distributed under the MPL-2.0 license.
# the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
# or see the "LICENSE" file for more details.

# If `vDC` is not specified, the default `vDC` will be used
# The `myVMID` is the ID of the VM. 
terraform import cloudavenue_vm.example myVAPP.myVMID

# or you can specify the vDC
terraform import cloudavenue_vm.example myVDC.myVAPP.myVMID
