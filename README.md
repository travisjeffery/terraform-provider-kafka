# Terraform Kafka Topic Provider

This is a [Terraform][terraform] provider for managing the Kafka topics with
Terraform.


## Installation

1. Download the latest compiled binary from [GitHub releases][releases].

1. Unzip/untar the archive.

1. Move it into `$HOME/.terraform.d/plugins`:

    ```sh
    $ mkdir -p $HOME/.terraform.d/plugins
    $ mv terraform-provider-kafka-topic $HOME/.terraform.d/plugins/terraform-provider-kafka-topic
    ```

1. Create your Terraform configurations as normal, and run `terraform init`:

    ```sh
    $ terraform init
    ```

    This will find the plugin locally.


## Usage

1. Create a Terraform configuration file:

    ```hcl
    resource "kafka_topic" "example" {
      name     = "example"
      num_partitions = "8"
      replication_factor = "3"
    }
    ```

1. Run `terraform init` to pull in the provider:

    ```sh
    $ terraform init
    ```

1. Run `terraform plan` and `terraform apply` to interact with the filesystem:

    ```sh
    $ terraform plan

    $ terraform apply
    ```

## Examples

For more examples, please see the [examples][examples] folder in this
repository.

## Reference

### Filesystem Reader

#### Usage

```hcl
resource "filesystem_file_reader" "read" {
  path = "my-file.txt"
}
```

## License

MIT

---

- [travisjeffery.com](http://travisjeffery.com)
- GitHub [@travisjeffery](https://github.com/travisjeffery)
- Twitter [@travisjeffery](https://twitter.com/travisjeffery)
- Medium [@travisjeffery](https://medium.com/@travisjeffery)
