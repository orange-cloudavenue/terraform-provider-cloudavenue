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
7. If needed, make a changelog of your changes

Ensure to use a good commit hygiene and follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.

##  Contributing documentation

Documentation is generated from the code using [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs), you don't need to update the documentation manually in the `docs/` folder.

If you must make manual changes to the documentation, please ensure to edit template in `templates/` folder and run `make generate` to update the documentation.

To add examples or import command in documentation, please add them in the `examples/` folder.

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
export CLOUDAVENUE_VDC=your-vdc
TF_ACC=1 go test -v -count=1 ./internal/tests/your_test_folder
```

##  Changelog format

We use the go-changelog to generate and update the changelog from files created in the .changelog/ directory. It is important that when you raise your Pull Request, there is a changelog entry which describes the changes your contribution makes. Not all changes require an entry in the changelog, guidance follows on what changes do.

The changelog format requires an entry in the following format, where HEADER corresponds to the changelog category, and the entry is the changelog entry itself. The entry should be included in a file in the .changelog directory with the naming convention {PR-NUMBER}.txt. For example, to create a changelog entry for pull request 1234, there should be a file named .changelog/1234.txt.

``````markdown
```release-note:{HEADER}
{ENTRY}
```
``````

## Pull request types to CHANGELOG

The CHANGELOG is intended to show operator-impacting changes to the codebase for a particular version. If every change or commit to the code resulted in an entry, the CHANGELOG would become less useful for operators. The lists below are general guidelines and examples for when a decision needs to be made to decide whether a change should have an entry.

### Changes that should have a CHANGELOG entry

#### New resource

A new resource entry should only contain the name of the resource, and use the `release-note:new-resource` header.

``````markdown
```release-note:new-resource
cloudavenue_alb_pool
```
``````

#### New data source

A new data source entry should only contain the name of the data source, and use the `release-note:new-data-source` header.

``````markdown
```release-note:new-data-source
cloudavenue_alb_pool
```
``````

#### Resource and provider bug fixes

A new bug entry should use the `release-note:bug` header and have a prefix indicating the resource or data source it corresponds to, a colon, then followed by a brief summary. Use a `provider` prefix for provider level fixes.

``````markdown
```release-note:bug
resource/cloudavenue_alb_pool: Fix argument being optional
```
``````

#### Resource and provider enhancements

A new enhancement entry should use the `release-note:enhancement` header and have a prefix indicating the resource or data source it corresponds to, a colon, then followed by a brief summary. Use a `provider` prefix for provider level enhancements.

``````markdown
```release-note:enhancement
resource/cloudavenue_alb_pool: Add new argument
```
``````

#### Deprecations

A deprecation entry should use the `release-note:note` header and have a prefix indicating the resource or data source it corresponds to, a colon, then followed by a brief summary. Use a `provider` prefix for provider level changes.

``````markdown
```release-note:note
resource/cloudavenue_alb_pool: The old_attribute is being deprecated in favor of the new_attribute to support new feature
```
``````

#### Breaking changes and removals

A breaking-change entry should use the `release-note:breaking-change` header and have a prefix indicating the resource or data source it corresponds to, a colon, then followed by a brief summary. Use a `provider` prefix for provider level changes.

``````markdown
```release-note:breaking-change
resource/cloudavenue_alb_pool: This is a breaking change
```
``````

### Changes that should _not_ have a CHANGELOG entry

- Resource and provider documentation updates
- Testing updates
- Code refactoring
