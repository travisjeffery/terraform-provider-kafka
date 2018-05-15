package topic

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/terraform/helper/schema"
)

func resource() *schema.Resource {
	return &schema.Resource{
		Create: create,
		Read:   read,
		Delete: delete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the topic",
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
				Required:    false,
			},
		},
	}
}

func create(d *schema.ResourceData, meta interface{}) error {
	c, err := client(meta)
	if err != nil {
		return err
	}

	topic := d.Get("topic").(string)

	topicDetail := &sarama.TopicDetail{}
	topicDetail.NumPartitions = d.Get("num_partitions").(int32)
	topicDetail.ReplicationFactor = d.Get("replication_factor").(int16)
	for name, value := range d.Get("config_entries").(map[string]string) {
		topicDetail.ConfigEntries[name] = &value
	}

	topicDetails := make(map[string]*sarama.TopicDetail)
	topicDetails[topic] = topicDetail

	response, err := c.CreateTopics(&sarama.CreateTopicsRequest{
		TopicDetails: topicDetails,
	})
	if err != nil || response.TopicErrors == nil {
		return err
	}
	if err := response.TopicErrors[topic]; err != nil {
		return fmt.Errorf("topic error: %s", err.ErrMsg)
	}
	return nil
}

func read(d *schema.ResourceData, meta interface{}) error {
	c, err := client(meta)
	if err != nil {
		return err
	}

	metadata, err := c.GetMetadata(&sarama.MetadataRequest{Topics: []string{d.Get("topic").(string)}})
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

	response, err := c.DescribeConfigs(&sarama.DescribeConfigsRequest{
		Resources: []*sarama.ConfigResource{{
			Type: sarama.TopicResource,
			Name: topic.Name,
		}}},
	)
	if err != nil {
		return err
	}
	if len(response.Resources) != 1 {
		return fmt.Errorf("expected 1 resource in response")
	}
	resource := response.Resources[0]
	if resource.ErrorCode != int16(sarama.ErrNoError) {
		return fmt.Errorf(
			"resource error: code: %d, message: %s",
			resource.ErrorCode,
			resource.ErrorMsg,
		)
	}

	configs := make(map[string]string)
	for _, config := range resource.Configs {
		configs[config.Name] = config.Value
	}
	d.Set("config_entries", configs)

	return nil
}

func delete(d *schema.ResourceData, meta interface{}) error {
	c, err := client(meta)
	if err != nil {
		return err
	}

	topic := d.Get("topic").(string)

	response, err := c.DeleteTopics(&sarama.DeleteTopicsRequest{
		Topics: []string{topic},
	})
	if err != nil || response.TopicErrorCodes == nil {
		return err
	}
	if errCode := response.TopicErrorCodes[topic]; errCode != sarama.ErrNoError {
		return fmt.Errorf("topic error code: %s", errCode)
	}
	return nil
}

func client(meta interface{}) (*sarama.Broker, error) {
	client := meta.(sarama.Client)
	return client.Controller()
}
