# terraform-sops

A Terraform plugin for using files encrypted with [Mozilla sops](https://github.com/mozilla/sops).

**NOTE:** To prevent plaintext secrets from being written to disk, you *must* set up a secure remote state backend. See the [official docs](https://www.terraform.io/docs/state/sensitive-data.html) on _Sensitive Data in State_ for more information.

## Example

**NOTE:** All examples assume Terraform 0.13 or newer. For information about usage on older versions, see the [legacy usage docs](docs/legacy_usage.md).

Encrypt a file using Sops: `sops demo-secret.enc.json`

```json
{
  "password": "foo",
  "db": {"password": "bar"}
}
```
### sops_file

```hcl
terraform {
  required_providers {
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
  }
}

provider "sops" {}

data "sops_file" "demo-secret" {
  source_file = "demo-secret.enc.json"
}

output "root-value-password" {
  # Access the password variable from the map
  value = data.sops_file.demo-secret.data["password"]
}

output "mapped-nested-value" {
  # Access the password variable that is under db via the terraform map of data
  value = data.sops_file.demo-secret.data["db.password"]
}

output "nested-json-value" {
  # Access the password variable that is under db via the terraform object
  value = jsondecode(data.sops_file.demo-secret.raw).db.password
}
```

Sops also supports encrypting the entire file when in other formats. Such files can also be used by specifying `input_type = "raw"`:

```hcl
data "sops_file" "some-file" {
  source_file = "secret-data.txt"
  input_type = "raw"
}

output "do-something" {
  value = data.sops_file.some-file.raw
}
```

### sops_external
For use with reading files that might not be local. 

> `input_type` is required with this data source.

```hcl
terraform {
  required_providers {
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
  }
}

provider "sops" {}

# using sops/test-fixtures/basic.yaml as an example
data "local_file" "yaml" {
  filename = "basic.yaml"
}

data "sops_external" "demo-secret" {
  source     = data.local_file.yaml.content
  input_type = "yaml"
}

output "root-value-hello" {
  value = data.sops_external.demo-secret.data.hello
}

output "nested-yaml-value" {
  # Access the password variable that is under db via the terraform object
  value = yamldecode(data.sops_file.demo-secret.raw).db.password
}
```

## Install

For Terraform 0.13 and later, specify the source and version in a `required_providers` block:

```hcl
terraform {
  required_providers {
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
  }
}
```

## Migrating state from older Terraform versions
To migrate a state from Terraform 0.12 or older, there is a need to change how the provider is referenced. After adding the `required_providers` block as shown above, Terraform provides a command to do the migration:

```shell
terraform state replace-provider registry.terraform.io/-/sops registry.terraform.io/carlpett/sops
```

## Development
Building and testing is most easily performed with `make build` and `make test` respectively.

The PGP key used for encrypting the test cases is found in `test/testing-key.pgp`. You can import it with `gpg --import test/testing-key.pgp`.

## Transitioning to Terraform 0.13 (0.12.20+) provider required blocks.

With Terraform 0.13 (0.12.20+), providers are available/downloaded via the [terraform registry](https://registry.terraform.io/providers/carlpett/sops/latest) via a required_providers block (Yay!)

```
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    sops = {
      source = "carlpett/sops"
      version = "0.5.1"
    }
  }
}
```
A prerequisite when converting is that you must from remove the data source block from the previous SOP provider in your `terraform.state` file. 

If not you will be greeted with: 

```
- Finding latest version of -/sops...

Error: Failed to query available provider packages

Could not retrieve the list of available versions for provider -/sops:
provider registry registry.terraform.io does not have a provider named
registry.terraform.io/-/sops
```