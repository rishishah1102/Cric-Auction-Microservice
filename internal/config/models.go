package config

import "time"

// Config the struct for all the configuration of auction web
type Config struct {
	MongoDB  MongoDB
	Postgres Postgres
	Redis    Redis
	Kafka    Kafka
	DbName   string
}

// MongoDB is the struct for mongo configurations
type MongoDB struct {
	MongoURI string
	Timeout  time.Duration
}

// Postgres is the struct for postgres configurations
type Postgres struct {
	PostgresURI string
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
