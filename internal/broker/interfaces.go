package broker

import (
	"context"
	"github.com/diemensa/event-analytics-service/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Producer - Kafka
type Producer interface {
	Send(ctx context.Context, e *model.Event) error
	Close() error
}

// Publisher - RabbitMQ
type Publisher interface {
	Publish(ctx context.Context, e *model.Event) error
	Consume(queueName string) (<-chan amqp.Delivery, error)
	Close() error
}

// Умер старый тимлид. Много лет писал он код на бигтех.
// Все уважали старого тимлида. Однако была у него одна странность.
// Стоя в переговорной, прежде чем отдать команду, он всегда доставал
// из кармана маленькую бумажку, смотрел на нее, и затем уже отдавал команду.
// Никто никогда не видел, что написано в этой бумажке.
// После похорон программиста мидлы, снедаемые любопытством, собрались
// в переговорной. Мидл+ торжественно подошел к рубашке
// тимлида, вынул бумажку и развернул ее.
// На бумажке было написано:
// "Producer - Kafka, Publisher - RabbitMQ"
