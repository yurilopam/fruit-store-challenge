package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserEvent struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type StoredUser struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   string `gorm:"uniqueIndex"`
	Username string
	Role     string
}

func main() {
	_ = godotenv.Load()
	brokers := getenv("KAFKA_BROKERS", "kafka:9092")
	topic := getenv("KAFKA_TOPIC_USERS", "users.created")

	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	_ = db.AutoMigrate(&StoredUser{})

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(brokers, ","),
		Topic:   topic,
		GroupID: "user-consumer-group",
	})
	defer r.Close()

	log.Println("[user-consumer] Listening topic:", topic)
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("read error:", err)
			continue
		}
		var ev UserEvent
		if err := json.Unmarshal(m.Value, &ev); err != nil {
			log.Println("json error:", err)
			continue
		}
		if err := db.Create(&StoredUser{UserID: ev.ID, Username: ev.Username, Role: ev.Role}).Error; err != nil {
			log.Println("db error:", err)
			continue
		}
		log.Printf("stored user: %s (%s)\n", ev.Username, ev.ID)
	}
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
