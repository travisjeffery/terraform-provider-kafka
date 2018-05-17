package topic

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	tfresource "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var (
	createConfig = `
			resource "kafka_topic" "example" {
				name = "example"
				num_partitions = "3"
				replication_factor = "1"
				config_entries = {
	  			  retention.bytes = "102400"
	  			  cleanup.policy = "compact"
	  			}
			}
		`

	updateConfig = `
			resource "kafka_topic" "example" {
				name = "example"
				num_partitions = "4"
				replication_factor = "1"
				config_entries = {
	  			  retention.bytes = "1024"
	  			  cleanup.policy = "compact"
	  			}
			}
		`
)

func TestTopic(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		t.Parallel()

		tfresource.UnitTest(t, tfresource.TestCase{
			IsUnitTest: true,
			PreCheck:   func() { testAccPreCheck(t) },
			Providers:  testAccProviders,
			Steps: []tfresource.TestStep{
				{
					Config: createConfig,
					Check: func(s *terraform.State) error {
						attrs := s.RootModule().Resources["kafka_topic.example"].Primary.Attributes

						if act, exp := attrs["name"], "example"; act != exp {
							t.Errorf("expected %q to be %q", act, exp)
						}
						if act, exp := attrs["num_partitions"], "3"; act != exp {
							t.Errorf("expected %q to be %q", act, exp)
						}
						if act, exp := attrs["replication_factor"], "1"; act != exp {
							t.Errorf("expected %q to be %q", act, exp)
						}
						if act, exp := attrs["config_entries"], map[string]interface{}{
							"retention.bytes": "102400",
							"cleanup.policy":  "compact",
						}; reflect.DeepEqual(act, exp) {
							t.Errorf("expected %v to be %v", act, exp)
						}

						return nil
					},
				},
				{
					Config: updateConfig,
					Check: func(s *terraform.State) error {
						attrs := s.RootModule().Resources["kafka_topic.example"].Primary.Attributes

						if act, exp := attrs["name"], "example"; act != exp {
							t.Errorf("expected %q to be %q", act, exp)
						}
						if act, exp := attrs["num_partitions"], "4"; act != exp {
							t.Errorf("expected %q to be %q", act, exp)
						}
						if act, exp := attrs["replication_factor"], "1"; act != exp {
							t.Errorf("expected %q to be %q", act, exp)
						}
						if act, exp := attrs["config_entries"], map[string]interface{}{
							"retention.bytes": "1024",
							"cleanup.policy":  "compact",
						}; reflect.DeepEqual(act, exp) {
							t.Errorf("expected %v to be %v", act, exp)
						}

						return nil
					},
				},
			},
			CheckDestroy: func(*terraform.State) error {
				cfg := sarama.NewConfig()
				cfg.Version = sarama.V1_0_0_0
				client, err := sarama.NewClient([]string{"localhost:9092"}, cfg)
				if err != nil {
					return fmt.Errorf("failed to create kafka client: %s", err)
				}
				for count := 0; count < 3; count++ {
					if err := client.RefreshMetadata(); err != nil {
						continue
					}
					topics, err := client.Topics()
					if err != nil {
						continue
					}
					var found bool
					for _, topic := range topics {
						if topic == "example" {
							found = true
						}
					}
					if !found {
						return nil
					}
					time.Sleep(1 * time.Second)
				}
				return fmt.Errorf("topic wasn't removed")
			},
		})
	})
}
