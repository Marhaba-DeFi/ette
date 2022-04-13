package block

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"

	d "github.com/itzmeanjan/ette/app/data"
	"github.com/itzmeanjan/ette/app/db"
)

// PublishBlock - Attempts to publish block data to Kafka pubsub channel
func PublishBlock(block *db.PackedBlock, kafkaInfo *d.KafkaInfo) bool {

	if block == nil {
		return false
	}

	_block := &d.Block{
		Hash:                block.Block.Hash,
		Number:              block.Block.Number,
		Time:                block.Block.Time,
		ParentHash:          block.Block.ParentHash,
		Difficulty:          block.Block.Difficulty,
		GasUsed:             block.Block.GasUsed,
		GasLimit:            block.Block.GasLimit,
		Nonce:               block.Block.Nonce,
		Miner:               block.Block.Miner,
		Size:                block.Block.Size,
		StateRootHash:       block.Block.StateRootHash,
		UncleHash:           block.Block.UncleHash,
		TransactionRootHash: block.Block.TransactionRootHash,
		ReceiptRootHash:     block.Block.ReceiptRootHash,
		ExtraData:           block.Block.ExtraData,
	}

	// Marshall block data
	byteBlock, marshallErr := _block.MarshalBinary()
	if marshallErr != nil {
		log.Println("Error marshalling block: ", marshallErr.Error())
		return false
	}

	kafkaWriteErr := kafkaInfo.KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Topic: "new-block",
			Key:   []byte("new-block"),
			Value: byteBlock,
		},
	)

	if kafkaWriteErr != nil {

		log.Printf("❗️ Failed to publish block %d : %s\n", block.Block.Number, kafkaWriteErr.Error())
		return false

	}

	log.Printf("📎 Published block %d\n", block.Block.Number)

	return PublishTxs(block.Block.Number, block.Transactions, kafkaInfo)

}
