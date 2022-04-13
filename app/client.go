package app

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	cfg "github.com/itzmeanjan/ette/app/config"
	d "github.com/itzmeanjan/ette/app/data"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
)

// Connect to blockchain node, either using HTTP or Websocket connection
// depending upon true/ false, passed to function, respectively
func getClient(isRPC bool) *ethclient.Client {
	var client *ethclient.Client
	var err error

	if isRPC {
		client, err = ethclient.Dial(cfg.Get("RPCUrl"))
	} else {
		client, err = ethclient.Dial(cfg.Get("WebsocketUrl"))
	}

	if err != nil {
		log.Fatalf("[!] Failed to connect to blockchain : %s\n", err.Error())
	}

	return client
}

func getKafkaClient(hostAddress string, hostPort string) *d.KafkaInfo {

	// Create a new kafka producer
	return &d.KafkaInfo{
		KafkaWriter: &kafka.Writer{
			Addr:     kafka.TCP(net.JoinHostPort(hostAddress, hostPort)),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// Creates connection to Redis server & returns that handle to be used for further communication
func getRedisClient() *redis.Client {

	var options *redis.Options

	// If password is given in config file
	if cfg.Get("RedisPassword") != "" {

		options = &redis.Options{
			Network:  cfg.Get("RedisConnection"),
			Addr:     cfg.Get("RedisAddress"),
			Password: cfg.Get("RedisPassword"),
			DB:       0,
		}

	} else {
		// If password is not given, attempting to connect with out it
		//
		// Though this is not recommended
		options = &redis.Options{
			Network: cfg.Get("RedisConnection"),
			Addr:    cfg.Get("RedisAddress"),
			DB:      0,
		}

	}

	_redis := redis.NewClient(options)
	// Checking whether connection was successful or not
	if err := _redis.Ping(context.Background()).Err(); err != nil {
		return nil
	}

	return _redis

}
