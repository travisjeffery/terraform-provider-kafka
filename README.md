# Terraform Kafka Topic Provider

This is a [Terraform][terraform] provider for managing the Kafka topics with
Terraform.

Why use this Kafka provider over others?

- Uses Kafka's new admin APIs rather than
shelling out to old bash scripts
- Supports adding partitions and altering
configs

## Installation

1. Download the latest compiled binary from [GitHub releases][releases].

1. Unzip/untar the archive.

1. Move it into `$HOME/.terraform.d/plugins`:

    ```sh
    $ mkdir -p $HOME/.terraform.d/plugins
    $ mv terraform-provider-kafka $HOME/.terraform.d/plugins/terraform-provider-kafka
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
	  name: "example"
	  num_partitions: "8"
	  replication_factor: "1"
	  config_entries: {
	    "retention.bytes": "102400"
	    "cleanup.policy": "compact"
	  }
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

## License

MIT

---

- [travisjeffery.com](http://travisjeffery.com)
- GitHub [@travisjeffery](https://github.com/travisjeffery)
- Twitter [@travisjeffery](https://twitter.com/travisjeffery)
- Medium [@travisjeffery](https://medium.com/@travisjeffery)
