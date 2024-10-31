# Terraform Provider Paperspace

- Tutorials for creating Terraform providers can be found on the [HashiCorp Developer](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework) platform.
- To test how documents will render in the Terraform Registry, [Terraform Registry Doc Preview Tool](https://registry.terraform.io/tools/doc-preview) can be used.
- [Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework)
- [Paperspace API Reference](https://docs.digitalocean.com/reference/paperspace/pspace/api-reference/)
- [Paperspace Documentation](https://docs.digitalocean.com/products/paperspace/)


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22


Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

First, find the `GOBIN` path where Go installs your binaries. Your path may vary depending on how your Go environment variables are configured:

```shell
$ go env GOBIN
/Users/<Username>/go/bin
```

If the `GOBIN` go environment variable is not set, use the default path, `/Users/<Username>/go/bin`.


Create a new file called `.terraformrc` in your home directory (`~`), then add the `dev_overrides` block into the file. Change the `<PATH>` to the value returned from the go env `GOBIN` command above:

```
provider_installation {
  dev_overrides {
    "registry.terraform.io/metaphysic/paperspace" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

To compile the provider, run `make install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`. 

*Note:* Acceptance tests create real resources, and often cost money to run.


## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```