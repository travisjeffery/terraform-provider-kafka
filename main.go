package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sethvargo/terraform-provider-kafka-topic/topic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: topic.Provider,
	})
}
