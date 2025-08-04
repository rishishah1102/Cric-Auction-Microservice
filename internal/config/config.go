package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// LoadConfig loads the environment variables into struct
func LoadConfig() (cfg *Config) {
	godotenv.Load()

	cfg = &Config{
		MongoDB: MongoDB{
			MongoURI: os.Getenv("MONGO_URI"),
			Timeout:  10 * time.Second,
		},
		Postgres: Postgres{
			PostgresURI: os.Getenv("POSTGRES_URI"),
			Timeout:     10 * time.Second,
		},
		Redis: Redis{
			RedisURI:      os.Getenv("REDIS_URI"),
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
		},
		Kafka: Kafka{
			KafkaBroker: os.Getenv("KAFKA_BROKER"),
		},
		DbName: os.Getenv("DB_NAME"),
	}

	return cfg
}
