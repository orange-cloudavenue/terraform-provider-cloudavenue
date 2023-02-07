#  Cloud Avenue Provider for Terraform

This is the Cloud Avenue provider for Terraform. It allows you to manage Cloud Avenue resources.

Usefull links:

* [Cloud Avenue Provider documentation](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs)
* [Terraform Documentation](https://www.terraform.io/docs/language/index.html)

##  Requirements

* [Terraform](https://www.terraform.io/downloads.html) 1.x.x
* [Go](https://golang.org/doc/install) 1.19.x (to build the provider plugin)
* [Cloud Avenue Platform](https://cloud.orange-business.com/offres/infrastructure-iaas/cloud-avenue/)

##  Using the Cloud Avenue Provider

To quickly get started with the Cloud Avenue Provider, you can use the following example:

```hcl
terraform {
  required_providers {
    cloudavenue = {
      source = "orange-cloudavenue/cloudavenue"
      version = "0.1.0"
    }
  }
}

provider "cloudavenue" {
  org = "my-org"
  user = "my-user"
  password = "my-password"
}
```

##  Contributing

This provider is open source and contributions are welcome.

If you want to contribute to this provider, please read the [contributing guidelines](CONTRIBUTING.md).

You may also report issues or feature requests on the [GitHub issue tracker](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/new/choose).
