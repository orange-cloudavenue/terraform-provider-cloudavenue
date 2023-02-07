# Contributing

We welcome contributions of all kinds including code, issues or documentation.

##  Contributing Code

To contribute code, please follow this steps:

1. Communicate with us on the issue you want to work on
2. Make your changes
3. Test your changes
4. Update the documentation if needed and examples too
5. Ensure to run `make generate` without issues
6. Open a pull request

Ensure to use a good commit hygiene and follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.

##  Contributing documentation

Documentation is generated from the code using [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs), you don't need to update the documentation manually in the `docs` folder.

If you must make manual changes to the documentation, please ensure to edit template in `templates` folder and run `make generate` to update the documentation.

To add examples or import command in documentation, please add them in the `examples` folder.

##  Development environment

Get git submodules:

```console
make submodules
```

Run doc generation:

```console
make generate
```

Install provider locally:

```console
make install
```

Run tests:

```console
make test
```

Run acceptance tests:

```console
export CLOUDAVENUE_ORG=your-org
export CLOUDAVENUE_USER=your-user
export CLOUDAVENUE_PASSWORD=your-password
make testacc
```
