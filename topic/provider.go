package topic

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

type threadsafeClient struct {
	sarama.Client
	*sync.Mutex
}

// Provider returns the actual provider instance.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hosts": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Your Kafka host addresses.",
				DefaultFunc: func() (interface{}, error) {
					return getHosts()
				},
			},
			"tls_enable": {
				Type:        schema.TypeBool,
				Description: "Whether or not to use TLS when connecting to the broker.",
				Optional:    true,
			},
			"sasl_enable": {
				Type:        schema.TypeBool,
				Description: "Whether or not to use SASL auth when connecting to the broker.",
				Optional:    true,
			},
			"sasl_username": {
				Type:        schema.TypeString,
				Description: "Username for SASL/Plain authentication.",
				Optional:    true,
			},
			"sasl_password": {
				Type:        schema.TypeString,
				Description: "Password for SASL/Plain authentication.",
				Optional:    true,
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
	cfg.Version = sarama.V1_0_0_0

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

	hostsConfig := d.Get("hosts").([]interface{})
	if hostsConfig == nil || len(hostsConfig) < 1 {
		hostsConfig, _ = getHosts()
	}
	var hosts []string
	for _, host := range hostsConfig {
		hosts = append(hosts, host.(string))
	}

	log.Printf("[INFO] Initializing Kafka client with hosts: %v\n", hosts)

	client, err := sarama.NewClient(hosts, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kafka client")
	}

	return &threadsafeClient{client, new(sync.Mutex)}, nil
}

func getHosts() ([]interface{}, error) {
	hosts := strings.Split(os.Getenv("KAFKA_HOSTS"), ",")
	log.Printf("[INFO] hosts: %v\n", hosts)
	if len(hosts) < 1 {
		return []interface{}{}, nil
	}
	res := make([]interface{}, len(hosts))
	for i, v := range hosts {
		res[i] = v
	}
	return res, nil
}
