package config

import "time"

// MongoDB is the struct for mongo configurations
type MongoDB struct {
	MongoURI string
	DbName   string
	Timeout  time.Duration
}

// Postgres is the struct for postgres configurations
type Postgres struct {
	PostgresURI string
	DbName      string
	Timeout     time.Duration
}

// Redis is the struct for redis configurations
type Redis struct {
	RedisURI      string
	RedisPassword string
}

// Kafka is the struct for kafka configurations
type Kafka struct {
	KafkaBroker string
}
