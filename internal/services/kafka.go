package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct{ writer *kafka.Writer }

type KafkaOption func(*kafka.Writer)

func NewKafkaProducer(brokersCSV, topic string, opts ...KafkaOption) *KafkaProducer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(strings.Split(brokersCSV, ",")...),
		Topic:        topic,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		BatchTimeout: time.Millisecond * 10,
	}
	for _, opt := range opts {
		opt(w)
	}
	return &KafkaProducer{writer: w}
}

type UserEvent struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (p *KafkaProducer) PublishUser(ctx context.Context, ev UserEvent) error {
	b, _ := json.Marshal(ev)
	return p.writer.WriteMessages(ctx, kafka.Message{Value: b})
}
