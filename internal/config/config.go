package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port      string
	JWTSecret string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr       string
	RedisDB         int
	RedisPassword   string
	RedisTTLSeconds int

	KafkaBrokers    string
	KafkaTopicUsers string

	SeedAdminUsername string
	SeedAdminPassword string
}

func MustLoad() Config {
	redisDB := 0
	if v := os.Getenv("REDIS_DB"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			redisDB = i
		}
	}
	redisTTL := 60
	if v := os.Getenv("REDIS_TTL_SECONDS"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			redisTTL = i
		}
	}

	cfg := Config{
		Port:              env("PORT", "8080"),
		JWTSecret:         env("JWT_SECRET", "supersecret"),
		DBHost:            env("DB_HOST", "localhost"),
		DBPort:            env("DB_PORT", "5432"),
		DBUser:            env("DB_USER", "postgres"),
		DBPassword:        env("DB_PASSWORD", "postgres"),
		DBName:            env("DB_NAME", "fruitstore"),
		RedisAddr:         env("REDIS_ADDR", "localhost:6379"),
		RedisDB:           redisDB,
		RedisPassword:     os.Getenv("REDIS_PASSWORD"),
		RedisTTLSeconds:   redisTTL,
		KafkaBrokers:      env("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopicUsers:   env("KAFKA_TOPIC_USERS", "users.created"),
		SeedAdminUsername: env("SEED_ADMIN_USERNAME", "admin"),
		SeedAdminPassword: env("SEED_ADMIN_PASSWORD", "admin"),
	}
	log.Println("Config loaded")
	return cfg
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
