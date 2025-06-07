package broker

import (
	"context"
	"encoding/json"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: w}
}

func (p *KafkaProducer) Send(ctx context.Context, event *model.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
