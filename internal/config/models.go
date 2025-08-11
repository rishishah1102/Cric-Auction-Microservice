package config

// MongoDB is the struct for mongo configurations
type MongoDB struct {
	MongoURI string
	DbName   string
}

// Postgres is the struct for postgres configurations
type Postgres struct {
	PostgresURI string
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
