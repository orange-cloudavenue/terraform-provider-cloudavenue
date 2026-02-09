<p align="center">
  <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue">
    <img src="https://avatars.githubusercontent.com/u/1506386?s=150&v=4" alt="terraform-provider-cloudavenue" width="150">
  </a>
  <h3 align="center" style="font-weight: bold">Terraform Provider for CloudAvenue Iaas public offer</h3>
  <p align="center">
      <a href="https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs">
      <img alt="Terraform" src="https://img.shields.io/badge/terraform-%235835CC.svg?style=for-the-badge&logo=terraform&logoColor=white" />
    </a>
    <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/releases/latest">
      <img alt="Latest release" src="https://img.shields.io/github/v/release/orange-cloudavenue/terraform-provider-cloudavenue?style=for-the-badge&logo=starship&color=C9CBFF&logoColor=D9E0EE&labelColor=302D41&include_prerelease&sort=semver" />
    </a>
    <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/pulse">
      <img alt="Last commit" src="https://img.shields.io/github/last-commit/orange-cloudavenue/terraform-provider-cloudavenue?style=for-the-badge&logo=starship&color=8bd5ca&logoColor=D9E0EE&labelColor=302D41"/>
    </a>
    <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/stargazers">
      <img alt="Stars" src="https://img.shields.io/github/stars/orange-cloudavenue/terraform-provider-cloudavenue?style=for-the-badge&logo=starship&color=c69ff5&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues">
      <img alt="Issues" src="https://img.shields.io/github/issues/orange-cloudavenue/terraform-provider-cloudavenue?style=for-the-badge&logo=bilibili&color=F5E0DC&logoColor=D9E0EE&labelColor=302D41" />
    </a>
  </p>
  <p align="center">
    <a href="https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs"><strong>:green_book: Explore the docs</strong></a>
  </p>
</p>

> [!IMPORTANT]  
> All releases below [**0.35.0**](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/releases/tag/v0.35.0) are considered outdated after 01/10/2026. **Please update to the latest version** to avoid any issues with the end of life of the legacy authentication method.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [About this project](#about-this-project)
- [Supported Versions](#supported-versions)
- [Using the Cloud Avenue Provider](#using-the-cloud-avenue-provider)
- [Contributing](#contributing)
  - [Top contributors](#top-contributors)

## About this project

A [Terraform](https://www.terraform.io) provider to manage [CloudAvenue Iaas offer](https://cloud.orange-business.com/offres/infrastructure-iaas/cloud-avenue/).

Made with <span style="color: #e25555;">&#9829;</span> using [Go](https://golang.org/).

## Supported Versions

- Terraform v1.5
- Go v1.25

It doesn't mean that this provider won't run on previous versions of Terraform or Go, though.
It just means that we can't guarantee backward compatibility.

## Using the Cloud Avenue Provider

To quickly get started with the Cloud Avenue Provider, you can use the following example:

```hcl
terraform {
  required_providers {
    cloudavenue = {
      source = "orange-cloudavenue/cloudavenue"
      version = ">= 0.35.0"
    }
  }
}

provider "cloudavenue" {
  org = "my-org"
  user = "my-user"
  password = "my-password"
}
```

For more information, please refer to the [Cloud Avenue Provider documentation](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs).

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please read the [CONTRIBUTING.md](CONTRIBUTING.md) for more details on our code of conduct, and the process for submitting pull requests to us.

### Top contributors

<a href="https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=orange-cloudavenue/terraform-provider-cloudavenue" alt="contrib.rocks image" />
</a>
