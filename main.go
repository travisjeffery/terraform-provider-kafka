package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/deliveroo/terraform-provider-kafka/topic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: topic.Provider,
	})
}
