---
** Using Networks in CloudAvenue **
---
This page is a brief overview of the different types of networks available in CloudAvenue. For more detailed information, please refer to the [CloudAvenue documentation](https://docs.cloudavenue.com).

## **1. Creating Routed Networks**

Routed networks are used to connect different networks together. They allow for communication between different subnets and can be used to connect to external networks.
Two kind of routed networks are available in CloudAvenue:
- `cloudavenue_edgegateway_network_routed` - This resource is used to create a routed network on VDC scope only.
- `cloudavenue_vdcg_network_routed` - This resource is used to create a routed network in the CloudAvenue VDC Group. The network is available to all VDCs in the group.

### **1.1. Creating Routed Networks in a single VDC**

This resource allows you to create a routed network connected to an Edge Gateway on a VDC scope.
Now you have the ability to create routed networks on Edge Gateways using the `cloudavenue_edgegateway_network_routed` resource.
This resource is specifically designed for Edge Gateways scoped only in a VDC and is not compatible with VDC Groups.

#### Example Usage

```terraform
resource "cloudavenue_edgegateway_network_routed" "example" {
  name = "example"

  edge_gateway_id = cloudavenue_edgegateway.example.id

  gateway       = "192.168.1.254"
  prefix_length = 24

  dns1 = "1.1.1.1"
  dns2 = "8.8.8.8"

  dns_suffix = "example"

  static_ip_pool = [
    {
      start_address = "192.168.1.10"
      end_address   = "192.168.1.20"
    }
  ]
}
```

### **1.2. Associated Resource**

All resources that are created on the Edge Gateway are associated with the `cloudavenue_edgegateway` resource.
This means that you can create a routed network on an Edge Gateway and used it in conjunction with a security group of resource `cloudavenue_edgegateway_security_group`. This allows you to create a routed network that is secured by a security group on VDC scope ONLY.

#### Example Usage

```terraform
resource "cloudavenue_edgegateway_security_group" "example" {
  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
  name            = "example"
  description     = "This is an example security group"
  member_org_network_ids = [
    cloudavenue_network_routed.example.id
  ]
}

data "cloudavenue_edgegateways" "example" {}

resource "cloudavenue_network_routed" "example" {
  name        = "MyOrgNet"
  description = "This is an example Net"

  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id

  gateway       = "192.168.1.254"
  prefix_length = 24

  dns1 = "1.1.1.1"
  dns2 = "8.8.8.8"

  dns_suffix = "example"

  static_ip_pool = [
    {
      start_address = "192.168.1.10"
      end_address   = "192.168.1.20"
    }
  ]
}
```

---

### **1.3. Creating Routed Networks in VDC Groups**

This resource allows you to create a routed network connected to an Edge Gateway on a VDC Group scope.
Now you have the ability to create routed networks on VDC Groups using the `cloudavenue_vdcg_network_routed` resource.
This resource is specifically designed for VDC Groups and is not compatible with an Edge Gateways on a VDC scope.
This resource is scoped to a VDC Group and allows you to create routed networks that can be used across multiple VDCs within the group.

#### Example Usage

```terraform
resource "cloudavenue_vdcg_network_routed" "example" {
  name = "example"

  vdc_group_id    = cloudavenue_vdcg.example.id
  edge_gateway_id = cloudavenue_edgegateway.example.id

  gateway       = "192.168.1.254"
  prefix_length = 24

  dns1 = "1.1.1.1"
  dns2 = "8.8.8.8"

  dns_suffix = "example"

  static_ip_pool = [
    {
      start_address = "192.168.1.10"
      end_address   = "192.168.1.20"
    }
  ]
}
```

### **1.4. Associated Resource**

All resources that are created on the VDC Group are associated with the `cloudavenue_vdcg` resource.
This means that you can create a routed network on a VDC Group and used it in conjunction with a security group of resource `cloudavenue_vdcg_security_group`. This allows you to create a routed network that is secured by a security group on VDC Group scope ONLY.
This resource is specifically designed for VDC Groups and is not compatible with Edge Gateways on a VDC scope.

#### Example Usage

```terraform
resource "cloudavenue_vdcg_security_group" "example" {
  name        = "example"
  description = "Example security group"

  vdc_group_id = cloudavenue_vdcg.example.id

  member_org_network_ids = [
    cloudavenue_vdcg_network_isolated.example.id,
    cloudavenue_vdcg_network_routed.example.id
  ]
}
```

---

## **2. Creating Isolated Networks**

Isolated networks are used to create a private network that is not connected to any external networks. They are typically used for internal communication between virtual machines and can be used to create a secure environment for sensitive data.
Two kind of isolated networks are available in CloudAvenue:
- `cloudavenue_vdc_network_isolated` - This resource is used to create an isolated network on VDC scope only.
- `cloudavenue_vdcg_network_isolated` - This resource is used to create an isolated network in the CloudAvenue VDC Group. The network is available to all VDCs in the group.

### **2.1. Creating Isolated Networks in a single VDC**

This resource allows you to create an isolated network connected in a VDC.
Now you have the ability to create isolated networks on Edge Gateways using the `cloudavenue_vdc_network_isolated` resource.
This resource is specifically designed for Edge Gateways scoped only in a VDC and is not compatible with VDC Groups.

#### Example Usage

```terraform
resource "cloudavenue_vdc_network_isolated" "example" {
  name = "my-isolated-network"
  vdc  = cloudavenue_vdc.example.name

  gateway       = "192.168.0.1"
  prefix_length = 24

  dns1       = "192.168.0.2"
  dns2       = "192.168.0.3"
  dns_suffix = "example.local"
}
```

---
### **2.3. Creating Isolated Networks in VDC Groups**

This resource allows you to create an isolated network connected to an Edge Gateway on a VDC Group scope.
Now you have the ability to create isolated networks on VDC Groups using the `cloudavenue_vdcg_network_isolated` resource.
This resource is specifically designed for VDC Groups and is not compatible with an Edge Gateways on a VDC scope.
This resource is scoped to a VDC Group and allows you to create isolated networks that can be used across multiple VDCs within the group.

#### Example Usage

```terraform
resource "cloudavenue_vdcg_network_isolated" "example" {
  name           = "my-isolated-network"
  vdc_group_name = cloudavenue_vdcg.example.name

  gateway       = "192.168.0.1"
  prefix_length = 24

  dns1       = "192.168.0.2"
  dns2       = "192.168.0.3"
  dns_suffix = "example.local"
}
```

---

## **3. A brief Network Validation in Security Groups**

### **3.1. Validation for `cloudavenue_edgegateway_security_group`**

- ✅ Accepts: `cloudavenue_edgegateway_network_routed`
- ❌ Rejects:
  - `cloudavenue_vdcg_network_routed`
  - `cloudavenue_vdcg_network_isolated`
  - `cloudavenue_vdc_network_isolated`

### **3.2. Validation for `cloudavenue_vdcg_security_group`**

- ✅ Accepts:
  - `cloudavenue_vdcg_network_routed`
  - `cloudavenue_vdcg_network_isolated`
- ❌ Rejects:
  - `cloudavenue_vdc_network_isolated`
  - `cloudavenue_edgegateway_network_routed`

---