resource "cloudavenue_vm_affinity_rule" "example" {
  name     = "test"
  polarity = "Affinity"

  vm_ids = [
    "urn:vcloud:vm:70b78935-cb64-4418-9607-4e3aeabbd168",
    "urn:vcloud:vm:c3912ae5-bbd1-45ae-8b1e-694d0a405a95"
  ]
}