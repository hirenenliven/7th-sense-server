package queue

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"7th-sense-server/internal/config"
)

// Delivery is one message taken off the queue. Tag identifies it on its Session for Ack,
// Reject and Requeue.
type Delivery struct {
	Body []byte
	Tag  uint64
}

// Session is one connection's consumption of the queue. Its tags are valid only on it:
// when the connection is lost, Deliveries closes and every message not acknowledged is delivered
// again on the next session.
type Session struct {
	channel    *amqp.Channel
	deliveries chan Delivery
}

// Deliveries returns the session's messages, in the order the broker sent them; it is closed
// when the connection is lost.
func (session *Session) Deliveries() <-chan Delivery {
	return session.deliveries
}

// Ack acknowledges the message tag, or with multiple every message up to tag: they leave the queue.
func (session *Session) Ack(tag uint64, multiple bool) error {
	return session.channel.Ack(tag, multiple)
}

// Reject refuses the message tag for good: it goes to the dead-letter queue.
func (session *Session) Reject(tag uint64) error {
	return session.channel.Reject(tag, false)
}

// Requeue hands the message tag, or with multiple every message up to tag, back to the queue
// for another delivery.
func (session *Session) Requeue(tag uint64, multiple bool) error {
	return session.channel.Nack(tag, multiple, true)
}

// Consumer takes messages off the queue, at most Prefetch unacknowledged at a time.
type Consumer struct {
	settings config.RabbitMQConfig
	name     string
	log      Logger
	state    connectionState
}

// NewConsumer creates a Consumer; name tells its connection apart in the logs and the management UI.
func NewConsumer(settings config.RabbitMQConfig, name string, log Logger) *Consumer {
	return &Consumer{settings: settings, name: name, log: log}
}

// Run connects, consumes and calls handle with each session until ctx ends, reconnecting with
// backoff (1 s doubling to ReconnectMax) when the broker cannot be reached or the
// connection is lost. handle returns when the session's Deliveries closes or ctx ends; the
// connection is closed after it returns, so handle must settle its messages before.
func (consumer *Consumer) Run(ctx context.Context, handle func(ctx context.Context, session *Session)) {
	waits := retryWaits{max: consumer.settings.ReconnectMax}
	for ctx.Err() == nil {
		connection, session, err := consumer.open(ctx)
		if err != nil {
			wait := waits.next()
			consumer.state.set(false, err)
			consumer.log(slog.LevelWarn, "queue: writer cannot reach the broker, retrying", "consumer", consumer.name,
				"retryIn", wait.String(), "error", err)
			sleep(ctx, wait)
			continue
		}
		waits.reset()
		consumer.state.set(true, nil)
		consumer.log(slog.LevelInfo, "queue: writer consuming", "consumer", consumer.name, "queue", consumer.settings.Queue,
			"prefetch", consumer.settings.Prefetch)
		handle(ctx, session)
		connection.Close()
		if ctx.Err() == nil {
			consumer.state.set(false, nil)
			consumer.log(slog.LevelWarn, "queue: writer lost the broker, reconnecting", "consumer", consumer.name)
		}
	}
	consumer.state.set(false, nil)
}

// open connects, sets the prefetch and starts consuming. The broker's deliveries are copied into
// the session's channel so callers never see the AMQP types.
func (consumer *Consumer) open(ctx context.Context) (*amqp.Connection, *Session, error) {
	connection, channel, err := openChannel(consumer.settings, "writer "+consumer.name)
	if err != nil {
		return nil, nil, err
	}
	if err := channel.Qos(consumer.settings.Prefetch, 0, false); err != nil {
		connection.Close()
		return nil, nil, fmt.Errorf("set prefetch: %w", err)
	}
	brokerDeliveries, err := channel.ConsumeWithContext(ctx, consumer.settings.Queue, "7th-sense-writer-"+consumer.name,
		false, false, false, false, nil)
	if err != nil {
		connection.Close()
		return nil, nil, fmt.Errorf("consume %s: %w", consumer.settings.Queue, err)
	}
	session := &Session{channel: channel, deliveries: make(chan Delivery, consumer.settings.Prefetch)}
	go func() {
		defer close(session.deliveries)
		for delivery := range brokerDeliveries {
			session.deliveries <- Delivery{Body: delivery.Body, Tag: delivery.DeliveryTag}
		}
	}()
	return connection, session, nil
}
