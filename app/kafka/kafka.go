package kafka

import (
	cfg "github.com/itzmeanjan/ette/app/config"
	kafka "github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func Connect() *kafka.Writer {
	_writer := newKafkaWriter(cfg.Get("KAFKA_URL"), cfg.Get("KAFKA_TOPIC"))

	return _writer
}
