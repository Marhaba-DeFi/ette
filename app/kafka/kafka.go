package kafka

import (
	cfg "github.com/itzmeanjan/ette/app/config"
	kafka "github.com/segmentio/kafka-go"
)

func Connect() *kafka.Writer {
	_writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Get("KAFKA_URL")),
		Balancer: &kafka.LeastBytes{},
	}

	return _writer
}
