package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadMongoConfig loads the mongo variables into struct
func LoadMongoConfig() (cfg *MongoDB) {
	godotenv.Load()
	cfg = &MongoDB{
		MongoURI: os.Getenv("MONGO_URI"),
		DbName:   os.Getenv("DB_NAME"),
	}
	return cfg
}

// LoadPostgresConfig loads the postgres variables into struct
func LoadPostgresConfig() (cfg *Postgres) {
	godotenv.Load()
	cfg = &Postgres{
		PostgresURI: os.Getenv("POSTGRES_URI"),
	}
	return cfg
}

// LoadRedisConfig loads the redis variables into struct
func LoadRedisConfig() (cfg *Redis) {
	godotenv.Load()
	cfg = &Redis{
		RedisURI:      os.Getenv("REDIS_URI"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
	return cfg
}

// LoadKafkaConfig loads the kafka variables into struct
func LoadKafkaConfig() (cfg *Kafka) {
	godotenv.Load()
	cfg = &Kafka{
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
	}
	return cfg
}
