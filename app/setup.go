package app

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
	cfg "github.com/itzmeanjan/ette/app/config"
	d "github.com/itzmeanjan/ette/app/data"
	"github.com/itzmeanjan/ette/app/db"
	k "github.com/itzmeanjan/ette/app/kafka"
	q "github.com/itzmeanjan/ette/app/queue"
	"github.com/itzmeanjan/ette/app/rest/graph"
	kafka "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

// Enum implementation of RPC and WS
// type Connector int64

// const (
// 	RPC Connector = iota
// 	RPCUrl_Backup
// 	Websocket
// 	WebsocketUrl_Backup
// )

// func (s Connector) String() int64 {
// 	switch s {
// 	case RPC:
// 		return 0
// 	case RPCUrl_Backup:
// 		return 1
// 	case Websocket:
// 		return 2
// 	case WebsocketUrl_Backup:
// 		return 3
// 	}
// 	return 0
// }

// Setting ground up i.e. acquiring resources required & determining with
// some basic checks whether we can proceed to next step or not
func bootstrap(p Params) (*d.BlockChainNodeConnection, *redis.Client, *d.RedisInfo, *gorm.DB, *d.StatusHolder, *q.BlockProcessorQueue, *kafka.Writer) {

	err := cfg.Read(p.configFile)
	if err != nil {
		log.Fatalf("[!] Failed to read `.env` : %s\n", err.Error())
	}

	if !(cfg.Get("EtteMode") == "1" || cfg.Get("EtteMode") == "2" || cfg.Get("EtteMode") == "3" || cfg.Get("EtteMode") == "4" || cfg.Get("EtteMode") == "5") {
		log.Fatalf("[!] Failed to find `EtteMode` in configuration file\n")
	}

	var RPC_connector int = 0
	var WS_connector int = 2
	if p.down {
		RPC_connector = 1
		WS_connector = 3

	}

	// Maintaining both HTTP & Websocket based connection to blockchain
	_connection := &d.BlockChainNodeConnection{
		RPC:       getClient(RPC_connector),
		Websocket: getClient(WS_connector),
	}
	_redisClient := getRedisClient()

	if _redisClient == nil {
		log.Fatalf("[!] Failed to connect to Redis Server\n")
	}

	if err := _redisClient.FlushAll(context.Background()).Err(); err != nil {
		log.Printf("[!] Failed to flush all keys from redis : %s\n", err.Error())
	}

	_db := db.Connect()

	// Populating subscription plans from `.plans.json` into
	// database table, at application start up
	db.PersistAllSubscriptionPlans(_db, p.subscriptionPlansFile)

	// Passing db handle, to graph package, so that it can be used
	// for resolving graphQL queries
	graph.GetDatabaseConnection(_db)

	_status := &d.StatusHolder{
		State: &d.SyncState{
			BlockCountAtStartUp:     db.GetBlockCount(_db),
			MaxBlockNumberAtStartUp: db.GetCurrentBlockNumber(_db),
		},
		Mutex: &sync.RWMutex{},
	}

	_redisInfo := &d.RedisInfo{
		Client:            _redisClient,
		BlockPublishTopic: "block",
		TxPublishTopic:    "transaction",
		EventPublishTopic: "event",
	}

	// This is block processor queue
	_queue := q.New(db.GetCurrentBlockNumber(_db))

	_kafkaWriter := k.Connect()

	return _connection, _redisClient, _redisInfo, _db, _status, _queue, _kafkaWriter

}
