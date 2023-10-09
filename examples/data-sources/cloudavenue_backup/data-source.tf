data "cloudavenue_backup" "example" {
  type        = "vdc"
  target_name = data.cloudavenue_vdc.example.name
}