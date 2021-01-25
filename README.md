# Terraform Kafka Topic Provider

## Forked
This repo has been forked from https://github.com/travisjeffery/terraform-provider-kafka

## Description

This is a Terraform provider for managing Kafka topics with
Terraform.

Why use this Kafka provider?

- Supports adding partitions and altering configs
- Supports TLS/SASL
- Uses Kafka's new admin APIs rather than shelling out to old bash scripts

## Installation

1. Download the latest compiled binary from [GitHub releases](https://github.com/travisjeffery/terraform-provider-kafka/releases).

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
provider "kafka" {
  hosts = ["localhost:9092"]
}

resource "kafka_topic" "example" {
  name: "example"
  num_partitions: "8"
  replication_factor: "1"
  config_entries: {
      retention.bytes: "102400"
      cleanup.policy: "compact
  }
}
```

[There's parameters to set if you use TLS/SASL](https://github.com/travisjeffery/terraform-provider-kafka/blob/58dfc2e47748eb6a4f817a3e93d9848c1668c164/topic/provider.go#L18-L46).

1. Run `terraform init` to pull in the provider:

```sh
$ terraform init
```

1. Run `terraform plan` and `terraform apply` to interact with the filesystem:

```sh
$ terraform plan

$ terraform apply
```

## Importing topics

This provider supports importing externally created topics by their name. Assuming you've already created a topic declaration like the one above, you can get Terraform to manage the state of the existing topic:

```sh
$ terraform import kafka_topic.example example
```

## Examples

For more examples, please see the [examples](https://github.com/travisjeffery/terraform-provider-kafka/tree/master/examples) folder in this
repository.

## License

MIT

---

- [travisjeffery.com](http://travisjeffery.com)
- GitHub [@travisjeffery](https://github.com/travisjeffery)
- Twitter [@travisjeffery](https://twitter.com/travisjeffery)
- Medium [@travisjeffery](https://medium.com/@travisjeffery)
