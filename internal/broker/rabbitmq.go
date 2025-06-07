package broker

import (
	"context"
	"encoding/json"
	"github.com/diemensa/event-analytics-service/internal/model"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewRabbitPublisher(uri, queueName string) (*RabbitPublisher, error) {
	conn, err := amqp091.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
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
		ch.Close()
		conn.Close()
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
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)

}

func (rp *RabbitPublisher) Consume(queueName string) (<-chan amqp091.Delivery, error) {
	messages, err := rp.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (rp *RabbitPublisher) Close() error {

	if err := rp.channel.Close(); err != nil {
		return err

	}
	return rp.conn.Close()
}
