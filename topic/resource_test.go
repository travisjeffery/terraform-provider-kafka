package topic

import (
	"fmt"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	tfresource "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestTopic(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		t.Parallel()

		config := `
			resource "kafka_topic" "example" {
				name = "example"
				num_partitions = "3"
				replication_factor = "1"
			}
		`

		tfresource.UnitTest(t, tfresource.TestCase{
			IsUnitTest: true,
			PreCheck:   func() { testAccPreCheck(t) },
			Providers:  testAccProviders,
			Steps: []tfresource.TestStep{
				{
					Config: config,
					Check: func(s *terraform.State) error {
						attrs := s.RootModule().Resources["kafka_topic.example"].Primary.Attributes

						name := "example"
						if act, exp := attrs["name"], name; act != exp {
							t.Errorf("expected %q to be %q", act, exp)
						}
						if act, exp := attrs["num_partitions"], "3"; act != exp {
							t.Errorf("expected %q to be %q", act, exp)
						}
						if act, exp := attrs["replication_factor"], "1"; act != exp {
							t.Errorf("expected %q to be %q", act, exp)
						}
						return nil
					},
				},
			},
			CheckDestroy: func(*terraform.State) error {
				cfg := sarama.NewConfig()
				cfg.Version = sarama.V0_11_0_0
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
