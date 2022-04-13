package block

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"

	d "github.com/itzmeanjan/ette/app/data"
	"github.com/itzmeanjan/ette/app/db"
)

// PublishEvents - Iterate over all events & try to publish them on kafka
func PublishEvents(events []*db.Events, kafkaInfo *d.KafkaInfo) bool {

	if events == nil {
		return false
	}

	var status bool

	for _, e := range events {

		status = PublishEvent(e, kafkaInfo)
		if !status {
			break
		}

	}

	return status

}

// PublishEvent - Publishing event/ log entry to kafka topic, to be captured by subscribers
// and sent to client application, who are interested in this piece of data
// after applying filter
func PublishEvent(event *db.Events, kafkaInfo *d.KafkaInfo) bool {

	if event == nil {
		return false
	}

	data := &d.Event{
		Origin:          event.Origin,
		Index:           event.Index,
		Topics:          event.Topics,
		Data:            event.Data,
		TransactionHash: event.TransactionHash,
		BlockHash:       event.BlockHash,
	}

	// Marshall tx data
	eventBinary, eventMarshallErr := data.MarshalBinary()
	if eventMarshallErr != nil {
		log.Println("Error marshalling Event Data: ", eventMarshallErr.Error())
		return false
	}

	// Publish tx to Kafka
	kafkaWriteErr := kafkaInfo.KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Topic: "new-event",
			Key:   []byte(data.TransactionHash),
			Value: eventBinary,
		},
	)
	if kafkaWriteErr != nil {

		log.Printf("❗️ Failed to publish event data from TX hash %d : %s\n", data.TransactionHash, kafkaWriteErr.Error())
		return false

	}

	return true

}
