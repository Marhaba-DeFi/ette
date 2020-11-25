package db

import (
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"
)

// Tabler - ...
type Tabler interface {
	TableName() string
}

// Blocks - Mined block info holder table model
type Blocks struct {
	Hash                string       `gorm:"column:hash;type:char(66);primaryKey"`
	Number              uint64       `gorm:"column:number;type:bigint;not null"`
	Time                uint64       `gorm:"column:time;type:bigint;not null"`
	ParentHash          string       `gorm:"column:parenthash;type:char(66);not null"`
	Difficulty          string       `gorm:"column:difficulty;type:varchar;not null"`
	GasUsed             uint64       `gorm:"column:gasused;type:bigint;not null"`
	GasLimit            uint64       `gorm:"column:gaslimit;type:bigint;not null"`
	Nonce               uint64       `gorm:"column:nonce;type:bigint;not null"`
	Miner               string       `gorm:"column:miner;type:char(42);not null"`
	Size                float64      `gorm:"column:size;type:float(8);not null"`
	TransactionRootHash string       `gorm:"column:txroothash;type:char(66);not null"`
	ReceiptRootHash     string       `gorm:"column:receiptroothash;type:char(66);not null"`
	Transactions        Transactions `gorm:"foreignKey:blockhash"`
	Events              Events       `gorm:"foreignKey:blockhash"`
}

// TableName - Overriding default table name
func (Blocks) TableName() string {
	return "blocks"
}

// Transactions - Blockchain transaction holder table model
type Transactions struct {
	Hash      string `gorm:"column:hash;type:char(66);primaryKey"`
	From      string `gorm:"column:from;type:char(42);not null"`
	To        string `gorm:"column:to;type:char(42)"`
	Contract  string `gorm:"column:contract;type:char(42)"`
	Gas       uint64 `gorm:"column:gas;type:bigint;not null"`
	GasPrice  string `gorm:"column:gasprice;type:varchar;not null"`
	Cost      string `gorm:"column:cost;type:varchar;not null"`
	Nonce     uint64 `gorm:"column:nonce;type:bigint;not null"`
	State     uint64 `gorm:"column:state;type:smallint;not null"`
	BlockHash string `gorm:"column:blockhash;type:char(66);not null"`
	Events    Events `gorm:"foreignKey:txhash"`
}

// TableName - Overriding default table name
func (Transactions) TableName() string {
	return "transactions"
}

// Events - Events emitted from smart contracts to be held in this table
type Events struct {
	Origin          string         `gorm:"column:origin;type:char(42);not null"`
	Index           uint           `gorm:"column:index;type:integer;not null;primaryKey"`
	Topics          pq.StringArray `gorm:"column:topics;type:text[];not null"`
	Data            []byte         `gorm:"column:data;type:bytea"`
	TransactionHash string         `gorm:"column:txhash;type:char(66);not null"`
	BlockHash       string         `gorm:"column:blockhash;type:char(66);not null;primaryKey"`
}

// TableName - Overriding default table name
func (Events) TableName() string {
	return "events"
}

// Users - User address & created api key related info, holder table
type Users struct {
	Address   string    `gorm:"column:address;type:char(42);not null" json:"address"`
	APIKey    string    `gorm:"column:apikey;type:char(66);primaryKey" json:"apiKey"`
	TimeStamp time.Time `gorm:"column:ts;type:timestamp;not null" json:"timeStamp"`
	Enabled   bool      `gorm:"column:enabled;type:boolean;default:true" json:"enabled"`
}

// TableName - Overriding default table name
func (Users) TableName() string {
	return "users"
}

// ToJSON - Encodes into JSON, to be supplied when queried for apps created by user
func (u *Users) ToJSON() []byte {
	data, err := json.Marshal(u)
	if err != nil {
		log.Printf("[!] Failed to encode `ette` apps data to JSON : %s\n", err.Error())
		return nil
	}

	return data
}

// DeliveryHistory - For each request coming from client application
// we're keeping track of how much data gets sent back in response of their query
//
// This is to be used for controlling client application's access
// to resources they're requesting
type DeliveryHistory struct {
	ID         uint64    `gorm:"column:id;type:bigserial;primaryKey"`
	Client     string    `gorm:"column:client;type:char(42);not null"`
	TimeStamp  time.Time `gorm:"column:ts;type:timestamp;not null"`
	EndPoint   string    `gorm:"column:endpoint;type:varchar(100);not null"`
	DataLength uint64    `gorm:"column:datalength;type:bigint;not null"`
}

// TableName - Overriding default table name
func (DeliveryHistory) TableName() string {
	return "delivery_history"
}

// SubscriptionPlans - Allowed subscription plans, to be auto populated from
// .plans.json, at application start up
type SubscriptionPlans struct {
	ID            uint32 `gorm:"column:id;type:serial;primaryKey" json:"id"`
	Name          string `gorm:"column:name;type:varchar(20);not null;unique" json:"name"`
	DeliveryCount uint64 `gorm:"column:deliverycount;type:bigint;not null;unique" json:"deliveryCount"`
}

// TableName - Overriding default table name
func (SubscriptionPlans) TableName() string {
	return "subscription_plans"
}

// SubscriptionDetails - Keeps track of which ethereum address has subscription to which plan
// where, `subscriptionplan` refers to `id` in subscription_plans table
// and 	`address` refers to address in users table
type SubscriptionDetails struct {
	Address          string `gorm:"column:address;type:char(42);primaryKey" json:"address"`
	SubscriptionPlan uint32 `gorm:"column:subscriptionplan;type:int;not null" json:"subscriptionPlan"`
}

// TableName - Overriding default table name
func (SubscriptionDetails) TableName() string {
	return "subscription_details"
}
