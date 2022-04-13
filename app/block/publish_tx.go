package block

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"

	d "github.com/itzmeanjan/ette/app/data"
	"github.com/itzmeanjan/ette/app/db"
)

// PublishTxs - Publishes all transactions in a block to kafka
func PublishTxs(blockNumber uint64, txs []*db.PackedTransaction, kafkaInfo *d.KafkaInfo) bool {

	if txs == nil {
		return false
	}

	var eventCount uint64
	var status bool

	for _, t := range txs {

		status = PublishTx(t, kafkaInfo)
		if !status {
			break
		}

		// how many events are present in this block, in total
		eventCount += uint64(len(t.Events))

	}

	if !status {
		return status
	}

	log.Printf("📎 Published %d transactions of block %d\n", len(txs), blockNumber)
	log.Printf("📎 Published %d events of block %d\n", eventCount, blockNumber)

	return status

}

// PublishTx - Publishes tx & events in tx, related data to Kafka
func PublishTx(tx *db.PackedTransaction, kafkaInfo *d.KafkaInfo) bool {

	if tx == nil {
		return false
	}

	var pTx *d.Transaction

	if tx.Tx.To == "" {
		// This is a contract creation tx
		pTx = &d.Transaction{
			Hash:      tx.Tx.Hash,
			From:      tx.Tx.From,
			Contract:  tx.Tx.Contract,
			Value:     tx.Tx.Value,
			Data:      tx.Tx.Data,
			Gas:       tx.Tx.Gas,
			GasPrice:  tx.Tx.GasPrice,
			Cost:      tx.Tx.Cost,
			Nonce:     tx.Tx.Nonce,
			State:     tx.Tx.State,
			BlockHash: tx.Tx.BlockHash,
		}
	} else {
		// This is a normal tx, so we keep contract field empty
		pTx = &d.Transaction{
			Hash:      tx.Tx.Hash,
			From:      tx.Tx.From,
			To:        tx.Tx.To,
			Value:     tx.Tx.Value,
			Data:      tx.Tx.Data,
			Gas:       tx.Tx.Gas,
			GasPrice:  tx.Tx.GasPrice,
			Cost:      tx.Tx.Cost,
			Nonce:     tx.Tx.Nonce,
			State:     tx.Tx.State,
			BlockHash: tx.Tx.BlockHash,
		}
	}

	// Marshall tx data
	ptxBinary, txMarshallErr := pTx.MarshalBinary()
	if txMarshallErr != nil {
		log.Println("Error marshalling TX Data: ", txMarshallErr.Error())
		return false
	}

	// Publish tx to Kafka
	kafkaWriteErr := kafkaInfo.KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Topic: "new-tx",
			Key:   []byte(pTx.Hash),
			Value: ptxBinary,
		},
	)
	if kafkaWriteErr != nil {

		log.Printf("❗️ Failed to publish TX data from Hash %d : %s\n", pTx.Hash, kafkaWriteErr.Error())
		return false

	}

	return PublishEvents(tx.Events, kafkaInfo)

}
