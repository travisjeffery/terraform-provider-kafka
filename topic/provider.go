package topic

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns the actual provider instance.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hosts": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Your Kafka host addresses.",
			},
			"tls_enable": {
				Type:        schema.TypeBool,
				Description: "Whether or not to use TLS when connecting to the broker.",
			},
			"sasl_enable": {
				Type:        schema.TypeBool,
				Description: "Whether or not to use SASL auth when connecting to the broker.",
			},
			"sasl_username": {
				Type:        schema.TypeString,
				Description: "Username for SASL/Plain authentication.",
			},
			"sasl_password": {
				Type:        schema.TypeString,
				Description: "Password for SASL/Plain authentication.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kafka_topic": resource(),
		},
		ConfigureFunc: configure,
	}
}

func configure(d *schema.ResourceData) (interface{}, error) {
	cfg := sarama.NewConfig()

	if v, ok := d.GetOk("tls_enable"); ok {
		cfg.Net.TLS.Enable = v.(bool)
	}

	if v, ok := d.GetOk("sasl_enable"); ok {
		cfg.Net.SASL.Enable = v.(bool)
	}

	if v, ok := d.GetOk("sasl_username"); ok {
		cfg.Net.SASL.User = v.(string)
	}

	if v, ok := d.GetOk("sasl_password"); ok {
		cfg.Net.SASL.Password = v.(string)
	}

	client, err := sarama.NewClient(d.Get("hosts").([]string), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %s", err)
	}

	return client, nil
}
