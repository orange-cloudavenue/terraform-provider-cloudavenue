resource "cloudavenue_vapp_isolated_network" "example" {
  name                  = "MyVappNet"
  vapp_name             = cloudavenue_vapp.example.name
  gateway               = "192.168.10.1"
  netmask               = "255.255.255.0"
  dns1                  = "192.168.10.1"
  dns2                  = "192.168.10.3"
  dns_suffix            = "myvapp.biz"
  guest_vlan_allowed    = true
  retain_ip_mac_enabled = true

  static_ip_pool = [
    {
      start_address = "192.168.10.51"
      end_address   = "192.168.10.101"
    },
    {
      start_address = "192.168.10.10"
      end_address   = "192.168.10.30"
    }
  ]
}
