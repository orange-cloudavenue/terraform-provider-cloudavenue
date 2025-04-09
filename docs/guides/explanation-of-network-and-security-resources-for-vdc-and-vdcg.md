# **Terraform Guide: Explanation of Network and Security Resources for VDC and VDCG**

## **1. Creating Routed Networks**

### **1.1. Resource `cloudavenue_edgegateway_network_routed`**

This resource allows you to create a routed network connected to an Edge Gateway on a VDC scope.
Now you have the ability to create routed networks on Edge Gateways using the `cloudavenue_edgegateway_network_routed` resource.
This resource is specifically designed for Edge Gateways scoped only in a VDC and is not compatible with VDC Groups.

#### Example

```hcl
resource "cloudavenue_edgegateway_network_routed" "example" {
  name        = "edgegateway-network"
  edgegateway = cloudavenue_edgegateway.example.id
  gateway     = "192.168.1.1"
  prefix_length = 24
  dns1        = "8.8.8.8"
  dns2        = "8.8.4.4"

  description = "Routed network for Edge Gateway"
}
```

#### Diagram

![Edge Gateway Routed Network](https://github.com/user-attachments/assets/83b38b0d-431f-476d-9dd9-c540e95b4c6b)

---

### **1.2. Resource `cloudavenue_vdcg_network_routed`**

This resource allows you to create a routed network connected to an Edge Gateway on a VDC Group scope.
Now you have the ability to create routed networks on VDC Groups using the `cloudavenue_vdcg_network_routed` resource.
This resource is specifically designed for VDC Groups and is not compatible with an Edge Gateways on a VDC scope.
This resource is scoped to a VDC Group and allows you to create routed networks that can be used across multiple VDCs within the group.

#### Example

```hcl
resource "cloudavenue_vdcg_network_routed" "example" {
  name            = "vdcg-network"
  vdc_group_id    = cloudavenue_vdcg.example.id

  edge_gateway_id = "10.0.0.1"
  prefix_length   = 24
  
  dns1       = "1.1.1.1"
  dns2       = "1.0.0.1"
}
```

#### Diagram

![VDC Group Routed Network](https://github.com/user-attachments/assets/e105ae2e-95d6-4984-bc52-8a50c918e5b8)

---

## **2. Security Groups**

### **2.1. Resource `cloudavenue_edgegateway_security_group`**

This resource allows you to create a security group for an Edge Gateway on a VDC scope. It only accepts routed networks created with `cloudavenue_edgegateway_network_routed`.
Like the routed network resource, this security group is scoped to an Edge Gateway on a VDC and is not compatible with VDC Groups.
This resource is specifically designed for Edge Gateways in a single VDC and allows you to manage security configurations for networks connected to the Edge Gateway.

#### Example

```hcl
resource "cloudavenue_edgegateway_security_group" "example" {
  name        = "edgegateway-security-group"
  description = "Security group for Edge Gateway"
  edge_gateway_id = cloudavenue_edgegateway.example.id

  member_org_network_ids = [
    cloudavenue_edgegateway_network_routed.example.id
  ]
}
```

#### Diagram

![Edge Gateway Security Group](https://github.com/user-attachments/assets/9044b244-6545-4907-8599-a75ef98a8af1)

---

### **2.2. Resource `cloudavenue_vdcg_security_group`**

This resource allows you to create a security group for a VDC Group. It only accepts routed or isolated networks created with `cloudavenue_vdcg_network_routed` or `cloudavenue_vdcg_network_isolated`.
This resource is specifically designed for VDC Groups and allows you to manage security configurations for networks connected to the VDC Group.
This resource is scoped to a VDC Group and allows you to manage security configurations for networks connected to the VDC Group.

#### Example

```hcl
resource "cloudavenue_vdcg_security_group" "example" {
  name      = "vdcg-security-group"
  description = "Security group for VDC Group"
  
  vdc_group_id = cloudavenue_vdcg.example.id
  
  member_org_network_ids = [
    cloudavenue_vdcg_network_routed.example.id,
    cloudavenue_vdcg_network_routed.example2.id
  ]

}
```

#### Diagram

![VDC Group Security Group](https://github.com/user-attachments/assets/bda221ad-75d6-4452-ba89-9ff19052e74e)

---

## **3. Network Validation in Security Groups**

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

## **4. Complete Example**

Here is a complete example combining network and security group resources:

```hcl
resource "cloudavenue_edgegateway_network_routed" "edge_network" {
  name        = "edge-network"
  edge_gateway_id = cloudavenue_edgegateway.example.id
  gateway         = "192.168.1.1"
  prefix_length = 24
  dns1        = "8.8.8.8"
  dns2        = "8.8.4.4"
}

resource "cloudavenue_vdcg_network_routed" "vdcg_network" {
  name       = "vdcg-network"
  vdc_group_id  = cloudavenue_vdcg.example.id
  edge_gateway_id = cloudavenue_edgegateway.example.id
  gateway    = "10.0.0.1"
  prefix_length = 24
  dns1       = "1.1.1.1"
  dns2       = "1.0.0.1"
}

resource "cloudavenue_edgegateway_security_group" "edge_sg" {
  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
  name            = "example"
  description     = "This is an example security group"
  member_org_network_ids = [
    cloudavenue_edgegateway_network_routed.example.id
  ]
}

resource "cloudavenue_vdcg_security_group" "vdcg_sg" {
  name        = "example"
  description = "Example security group"

  vdc_group_id = cloudavenue_vdcg.example.id

  member_org_network_ids = [
    cloudavenue_vdcg_network_routed.example.id
  ]
}
```

---

## **5. Migration from `cloudavenue_network_routed`**

If you are currently using `cloudavenue_network_routed`, migrate to the new resources:

- For Edge Gateways: Use `cloudavenue_edgegateway_network_routed`.  
- For VDC Groups: Use `cloudavenue_vdcg_network_routed`.

---

## **Conclusion**

This guide helps you use the new network and security group resources introduced in the CloudAvenue Terraform provider. Ensure you validate the compatible networks for each resource and follow best practices for managing network and security configurations.
