provider "kafka" {
  "hosts": ["localhost:9092"]
}

resource "kafka_topic" "example" {
  name: "example"
  num_partitions: "8"
  replication_factor: "1"
  config_entries: {
    "created_by": "kafka topic terraform provider"
    "retention.bytes": = 102400
    "cleanup.policy": "compact"
  }
}
