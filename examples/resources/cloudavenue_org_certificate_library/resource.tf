resource "cloudavenue_org_certificate_library" "example" {
  name        = "company-certificate"
  description = "dev environment certificate"
  certificate = chomp(file("/my/certificate/path/cert.pem"))
  private_key = chomp(file("/my/certificate/path/key.pem"))
}
