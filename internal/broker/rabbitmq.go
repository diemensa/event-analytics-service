package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/diemensa/event-analytics-service/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitPublisher(uri, queueName string) (*RabbitPublisher, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err = conn.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close rabbitmq connection: %v", err)
		}
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {

		err = ch.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close publisher channel: %v", err)
		}

		err = conn.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close publisher connection: %v", err)
		}

		return nil, err
	}

	return &RabbitPublisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (rp *RabbitPublisher) Publish(ctx context.Context, e *model.Event) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return rp.channel.PublishWithContext(
		ctx,
		"",
		rp.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent},
	)

}

func (rp *RabbitPublisher) Consume() (<-chan amqp.Delivery, error) {
	messageChan, err := rp.channel.Consume(
		rp.queue.Name,
		"events-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return messageChan, nil
}

func (rp *RabbitPublisher) Close() error {

	if err := rp.channel.Close(); err != nil {
		return err

	}
	return rp.conn.Close()
}
