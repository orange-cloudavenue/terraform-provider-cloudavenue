<div align="center">
    <a href="https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs">
      <img alt="Terraform" src="https://img.shields.io/badge/terraform-%235835CC.svg?style=for-the-badge&logo=terraform&logoColor=white" />
    </a>
    <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/releases/latest">
      <img alt="Latest release" src="https://img.shields.io/github/v/release/orange-cloudavenue/terraform-provider-cloudavenue?style=for-the-badge&logo=starship&color=C9CBFF&logoColor=D9E0EE&labelColor=302D41&include_prerelease&sort=semver" />
    </a>
    <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/pulse">
      <img alt="Last commit" src="https://img.shields.io/github/last-commit/orange-cloudavenue/terraform-provider-cloudavenue?style=for-the-badge&logo=starship&color=8bd5ca&logoColor=D9E0EE&labelColor=302D41"/>
    </a>
    <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/blob/main/LICENSE">
      <img alt="License" src="https://img.shields.io/github/license/orange-cloudavenue/terraform-provider-cloudavenue?style=for-the-badge&logo=starship&color=ee999f&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/stargazers">
      <img alt="Stars" src="https://img.shields.io/github/stars/orange-cloudavenue/terraform-provider-cloudavenue?style=for-the-badge&logo=starship&color=c69ff5&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues">
      <img alt="Issues" src="https://img.shields.io/github/issues/orange-cloudavenue/terraform-provider-cloudavenue?style=for-the-badge&logo=bilibili&color=F5E0DC&logoColor=D9E0EE&labelColor=302D41" />
    </a>
</div>

# Cloud Avenue Provider for Terraform

This is the Cloud Avenue provider for Terraform. It allows you to manage Cloud Avenue resources.

Useful links:

* [Cloud Avenue Provider documentation](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs)
* [Terraform Documentation](https://www.terraform.io/docs/language/index.html)

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 1.x.x
* [Go](https://golang.org/doc/install) 1.20.x (to build the provider plugin)
* [Cloud Avenue Platform](https://cloud.orange-business.com/offres/infrastructure-iaas/cloud-avenue/)

## Using the Cloud Avenue Provider

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

## Contributing

This provider is open source and contributions are welcome.

If you want to contribute to this provider, please read the [contributing guidelines](CONTRIBUTING.md).

You may also report issues or feature requests on the [GitHub issue tracker](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/new/choose).
