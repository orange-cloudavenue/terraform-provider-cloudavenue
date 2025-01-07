# SPDX-FileCopyrightText: Copyright (c) 2025 Orange
# SPDX-License-Identifier: Mozilla Public License 2.0
#
# This software is distributed under the MPL-2.0 license.
# the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
# or see the "LICENSE" file for more details.

# if vdc is not specified, the default vdc will be used
terraform import cloudavenue_vapp_isolated_network.example vapp_name.network_name

# if vdc is specified, the vdc will be used
terraform import cloudavenue_vapp_isolated_network.example vdc.vapp_name.network_name