package topic

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/terraform/helper/schema"
)

func resource() *schema.Resource {
	return &schema.Resource{
		Create:   create,
		Update:   update,
		Read:     read,
		Delete:   delete,
		Importer: &schema.ResourceImporter{State: importTopic},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the topic",
				ForceNew:    true,
				Required:    true,
			},

			"num_partitions": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Number of partitions.",
				Required:    true,
			},

			"replication_factor": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Replication factor.",
				Required:    true,
			},

			"config_entries": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "Config entries.",
				Optional:    true,
			},
		},
	}
}

func create(d *schema.ResourceData, meta interface{}) error {
	c, err := client(meta)
	if err != nil {
		return err
	}

	topic := d.Get("name").(string)

	d.SetId(topic)

	topicDetail := &sarama.TopicDetail{}
	topicDetail.NumPartitions = int32(d.Get("num_partitions").(int))
	topicDetail.ReplicationFactor = int16(d.Get("replication_factor").(int))
	topicDetail.ConfigEntries = make(map[string]*string)

	for name, value := range d.Get("config_entries").(map[string]interface{}) {
		strval := value.(string)
		topicDetail.ConfigEntries[name] = &strval
	}

	topicDetails := make(map[string]*sarama.TopicDetail)
	topicDetails[topic] = topicDetail

	response, err := c.CreateTopics(&sarama.CreateTopicsRequest{
		TopicDetails: topicDetails,
		Timeout:      time.Second * 15,
	})
	if err != nil || response.TopicErrors == nil {
		return err
	}
	if err := response.TopicErrors[topic]; err.Err != sarama.ErrNoError {
		return fmt.Errorf("topic error: %v", err)
	}

	return read(d, meta)
}

func update(d *schema.ResourceData, meta interface{}) error {
	c, err := client(meta)
	if err != nil {
		return err
	}

	topic := d.Get("name").(string)

	if d.HasChange("replication_factor") {
		return fmt.Errorf("can't update the replication factor currently")
	}

	if d.HasChange("num_partitions") {
		old, new := d.GetChange("num_partitions")
		if new.(int) < old.(int) {
			return fmt.Errorf("new num_partitions must be >= old num_partitions")
		}
		response, err := c.CreatePartitions(&sarama.CreatePartitionsRequest{
			Timeout: time.Second * 15,
			TopicPartitions: map[string]*sarama.TopicPartition{
				topic: {
					Count: int32(new.(int)),
				},
			},
		})
		if err != nil || response.TopicPartitionErrors == nil {
			return err
		}
		if err := response.TopicPartitionErrors[topic]; err.Err != sarama.ErrNoError {
			return fmt.Errorf("topic partition error: %v", err)
		}
	}

	if d.HasChange("config_entries") {
		_, new := d.GetChange("config_entries")

		configs := make(map[string]*string)

		for name, value := range new.(map[string]interface{}) {
			strval := value.(string)
			configs[name] = &strval
		}

		response, err := c.AlterConfigs(&sarama.AlterConfigsRequest{
			Resources: []*sarama.AlterConfigsResource{{
				Type:          sarama.TopicResource,
				Name:          topic,
				ConfigEntries: configs,
			}},
		})
		if err != nil {
			return err
		}
		for _, resource := range response.Resources {
			if resource.ErrorCode != int16(sarama.ErrNoError) {
				return fmt.Errorf(
					"resource error: code: %d, message: %s",
					resource.ErrorCode,
					resource.ErrorMsg,
				)
			}
		}
	}

	return read(d, meta)
}

func read(d *schema.ResourceData, meta interface{}) error {
	c, err := client(meta)
	if err != nil {
		return err
	}

	metadata, err := c.GetMetadata(&sarama.MetadataRequest{Topics: []string{d.Get("name").(string)}})
	if err != nil {
		return err
	}
	if len(metadata.Topics) != 1 {
		return fmt.Errorf("expected 1 topic in metadata")
	}

	topic := metadata.Topics[0]

	d.Set("name", topic.Name)
	d.Set("num_partitions", len(topic.Partitions))
	d.Set("replication_factor", len(topic.Partitions[0].Replicas)) // this work?

	if old, ok := d.GetOk("config_entries"); ok {
		read, err := configs(c, topic.Name)
		if err != nil {
			return err
		}
		new := make(map[string]interface{})
		for name, value := range read {
			if _, ok := old.(map[string]interface{})[name]; ok {
				new[name] = value
			}
		}
		d.Set("config_entries", new)
	}

	return nil
}

func delete(d *schema.ResourceData, meta interface{}) error {
	c, err := client(meta)
	if err != nil {
		return err
	}

	topic := d.Get("name").(string)

	response, err := c.DeleteTopics(&sarama.DeleteTopicsRequest{
		Topics:  []string{topic},
		Timeout: time.Second * 15,
	})
	if err != nil || response.TopicErrorCodes == nil {
		return err
	}
	if errCode := response.TopicErrorCodes[topic]; errCode != sarama.ErrNoError {
		return fmt.Errorf("topic error code: %s", errCode)
	}
	return nil
}

func importTopic(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := read(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func client(meta interface{}) (*sarama.Broker, error) {
	client := meta.(sarama.Client)
	controller, err := client.Controller()
	if err != nil {
		return nil, err
	}
	if ok, err := controller.Connected(); err != nil {
		return nil, err
	} else if ok {
		return controller, nil
	}
	if err = controller.Open(client.Config()); err != nil {
		return nil, err
	}
	return controller, nil
}

func configs(c *sarama.Broker, topic string) (map[string]string, error) {
	response, err := c.DescribeConfigs(&sarama.DescribeConfigsRequest{
		Resources: []*sarama.ConfigResource{{
			Type: sarama.TopicResource,
			Name: topic,
		}}},
	)
	if err != nil {
		return nil, err
	}
	if len(response.Resources) != 1 {
		return nil, fmt.Errorf("expected 1 resource in response")
	}
	resource := response.Resources[0]
	if resource.ErrorCode != int16(sarama.ErrNoError) {
		return nil, fmt.Errorf(
			"resource error: code: %d, message: %s",
			resource.ErrorCode,
			resource.ErrorMsg,
		)
	}

	configs := make(map[string]string)
	for _, config := range resource.Configs {
		configs[config.Name] = config.Value
	}

	return configs, nil
}
