const { Kafka } = require("kafkajs")

const kafka = new Kafka({
  clientId: "ette events consumer",
  brokers: ["localhost:29092"],
})

const main = async () => {
  const consumer = kafka.consumer({ groupId: "test-group" })

  await consumer.connect()

  // to test consuming a topic from ette the topic should conntain ${contractAddress}_{eventNameHash}

  await consumer.subscribe({
    topic:
      "0xB6E3d179E941Ed21627717834C098e3e56006C85_0xa9f73b89ea96c7b79f1e66b0f72cd59b336630ae66fc1d71c477db96779109e5",
    fromBeginning: true,
  })

  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      console.log({
        value: message.value.toString(),
      })
    },
  })
}

main()
