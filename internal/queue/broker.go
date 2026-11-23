// Package queue is the RabbitMQ side of the telemetry queue (ADR 0024): connect, declare the
// exchange and queues, publish with confirms, consume, and reconnect with backoff. No SQL and no
// business rules (like internal/opcua): internal/monitor decides what is published and writes
// what is consumed. Every log line goes through the Logger the caller passes in.
package queue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"7th-sense-server/internal/config"
)

const (
	// heartbeat: a broker that stops answering for about twice this long counts as lost.
	heartbeat = 10 * time.Second
	// firstRetryWait is the first wait before connecting again; it doubles up to ReconnectMax.
	firstRetryWait = time.Second
	// deadLetterExchangeSuffix / deadLetterQueueSuffix name the dead-letter exchange and queue
	// after the exchange and the queue.
	deadLetterExchangeSuffix = ".dlx"
	deadLetterQueueSuffix    = ".dead"
)

// Logger writes one line of the queue (internal/monitor passes its writeLog, logType queue).
type Logger func(level slog.Level, message string, attributes ...any)

// DeadLetterQueue is the queue the messages refused by the writer go to.
func DeadLetterQueue(settings config.RabbitMQConfig) string {
	return settings.Queue + deadLetterQueueSuffix
}

// SafeURL is the broker address without its password, for logs.
func SafeURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "(invalid RABBITMQ_URL)"
	}
	return parsed.Redacted()
}

// dial opens a connection to the broker named name in the management UI.
func dial(settings config.RabbitMQConfig, name string) (*amqp.Connection, error) {
	properties := amqp.NewConnectionProperties()
	properties.SetClientConnectionName("7th-sense " + name)
	connection, err := amqp.DialConfig(settings.URL, amqp.Config{Heartbeat: heartbeat, Properties: properties})
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", SafeURL(settings.URL), err)
	}
	return connection, nil
}

// declareTopology declares the durable direct exchange, the queue bound to it with the routing
// key, and the dead-letter exchange and queue. Declaring is idempotent, so it runs after every
// connect; a queue that already exists with other arguments is refused by the broker and the
// error says how to fix it.
func declareTopology(channel *amqp.Channel, settings config.RabbitMQConfig) error {
	deadLetterExchange := settings.Exchange + deadLetterExchangeSuffix
	if err := channel.ExchangeDeclare(settings.Exchange, amqp.ExchangeDirect, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange %s: %w", settings.Exchange, err)
	}
	// fanout: a dead letter keeps its routing key, and every one must reach the dead queue.
	if err := channel.ExchangeDeclare(deadLetterExchange, amqp.ExchangeFanout, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange %s: %w", deadLetterExchange, err)
	}
	deadQueue := DeadLetterQueue(settings)
	if _, err := channel.QueueDeclare(deadQueue, true, false, false, false, amqp.Table{amqp.QueueTypeArg: settings.QueueType}); err != nil {
		return queueDeclareError(deadQueue, err)
	}
	if err := channel.QueueBind(deadQueue, "", deadLetterExchange, false, nil); err != nil {
		return fmt.Errorf("bind queue %s: %w", deadQueue, err)
	}

	arguments := amqp.Table{amqp.QueueTypeArg: settings.QueueType, "x-dead-letter-exchange": deadLetterExchange}
	if _, err := channel.QueueDeclare(settings.Queue, true, false, false, false, arguments); err != nil {
		return queueDeclareError(settings.Queue, err)
	}
	if err := channel.QueueBind(settings.Queue, settings.RoutingKey, settings.Exchange, false, nil); err != nil {
		return fmt.Errorf("bind queue %s: %w", settings.Queue, err)
	}
	return nil
}

// queueDeclareError explains PRECONDITION_FAILED: the queue exists with other settings.
func queueDeclareError(queueName string, err error) error {
	var amqpError *amqp.Error
	if errors.As(err, &amqpError) && amqpError.Code == amqp.PreconditionFailed {
		return fmt.Errorf("declare queue %s: it already exists with other settings (type or arguments): "+
			"delete it in the management UI once it is empty: %w", queueName, err)
	}
	return fmt.Errorf("declare queue %s: %w", queueName, err)
}

// openChannel connects, opens a channel and declares the topology on it.
func openChannel(settings config.RabbitMQConfig, name string) (*amqp.Connection, *amqp.Channel, error) {
	connection, err := dial(settings, name)
	if err != nil {
		return nil, nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, nil, fmt.Errorf("open channel: %w", err)
	}
	if err := declareTopology(channel, settings); err != nil {
		connection.Close()
		return nil, nil, err
	}
	return connection, channel, nil
}

// QueueDepth returns how many messages wait in the queue and in its dead-letter queue,
// over a short connection of its own (for the writer's statistics log).
func QueueDepth(settings config.RabbitMQConfig) (waiting, dead int, err error) {
	connection, err := dial(settings, "queue depth")
	if err != nil {
		return 0, 0, err
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		return 0, 0, fmt.Errorf("open channel: %w", err)
	}
	queue, err := channel.QueueDeclarePassive(settings.Queue, true, false, false, false, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("inspect queue %s: %w", settings.Queue, err)
	}
	deadQueue, err := channel.QueueDeclarePassive(DeadLetterQueue(settings), true, false, false, false, nil)
	if err != nil {
		return queue.Messages, 0, fmt.Errorf("inspect queue %s: %w", DeadLetterQueue(settings), err)
	}
	return queue.Messages, deadQueue.Messages, nil
}

// retryWaits gives the waits between connects: 1 s, doubling up to max.
type retryWaits struct {
	current time.Duration
	max     time.Duration
}

// next returns the next wait.
func (waits *retryWaits) next() time.Duration {
	if waits.current == 0 {
		waits.current = firstRetryWait
	} else {
		waits.current = min(waits.current*2, waits.max)
	}
	return waits.current
}

// reset makes the next wait 1 s again.
func (waits *retryWaits) reset() {
	waits.current = 0
}

// sleep waits for duration or until ctx ends; it reports whether ctx is still running.
func sleep(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

// connectionState is whether a side is connected and the last error it met, for the logs.
type connectionState struct {
	mu        sync.Mutex
	connected bool
	lastError string
}

func (state *connectionState) set(connected bool, err error) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.connected = connected
	if err != nil {
		state.lastError = err.Error()
	}
}

func (state *connectionState) get() (bool, string) {
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.connected, state.lastError
}
