package queue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"7th-sense-server/internal/config"
)

const (
	// publishGroupSize is how many messages are published before their confirms are awaited
	// together: one round trip per group instead of one per message. Sized so a slow confirm
	// (the broker's fsync on a busy disk) does not cap the rate: 500 per round trip filled the
	// publish buffer at 5,000 readings/s (ADR 0027). A group only takes what is buffered.
	publishGroupSize = 5000
	// confirmTimeout: a group the broker has not confirmed by then is published again.
	confirmTimeout = 30 * time.Second
)

// errPublishNacked: the broker refused a message (a broker fault, such as a full disk); it is
// published again.
var errPublishNacked = errors.New("broker refused the message")

// PublisherStats are the publisher's counters since it was created.
type PublisherStats struct {
	Connected  bool
	LastError  string
	Buffered   int   // messages waiting in memory now
	Published  int64 // messages sent to the broker (a message sent again counts again)
	Confirmed  int64 // messages the broker confirmed: safe in the queue
	Failed     int64 // sends refused or not confirmed (each is retried)
	Fallback   int64 // messages Publish refused (buffer full or publisher stopped): written directly
	Reconnects int64 // connects after the first one
}

// Publisher sends messages to the exchange with publisher confirms. Publish only puts a
// message in a bounded buffer, so a caller never waits for the broker; one goroutine (Run) owns
// the connection and sends the buffer, retrying what the broker did not confirm.
type Publisher struct {
	settings config.RabbitMQConfig
	log      Logger
	buffer   chan []byte

	acceptMu  sync.RWMutex // Publish holds it shared; stopping holds it to close the door
	accepting bool

	retry [][]byte // the group being sent, then what the broker did not confirm; only Run touches it

	state      connectionState
	published  atomic.Int64
	confirmed  atomic.Int64
	failed     atomic.Int64
	fallback   atomic.Int64
	reconnects atomic.Int64
}

// NewPublisher creates a Publisher with a buffer of PublishBuffer messages; it accepts
// messages at once, and Run sends them.
func NewPublisher(settings config.RabbitMQConfig, log Logger) *Publisher {
	return &Publisher{settings: settings, log: log, buffer: make(chan []byte, settings.PublishBuffer), accepting: true}
}

// Publish hands body to the publisher without waiting. It returns false when the buffer is full
// or the publisher has stopped: the caller then stores the data itself.
func (publisher *Publisher) Publish(body []byte) bool {
	publisher.acceptMu.RLock()
	defer publisher.acceptMu.RUnlock()
	if publisher.accepting {
		select {
		case publisher.buffer <- body:
			return true
		default:
		}
	}
	publisher.fallback.Add(1)
	return false
}

// Stats returns the publisher's counters.
func (publisher *Publisher) Stats() PublisherStats {
	connected, lastError := publisher.state.get()
	return PublisherStats{
		Connected: connected, LastError: lastError, Buffered: len(publisher.buffer),
		Published: publisher.published.Load(), Confirmed: publisher.confirmed.Load(), Failed: publisher.failed.Load(),
		Fallback: publisher.fallback.Load(), Reconnects: publisher.reconnects.Load(),
	}
}

// Run connects and sends the buffer until ctx ends, reconnecting with backoff (1 s doubling to
// ReconnectMax) whenever the broker cannot be reached. A broker that is down never
// stops it: the buffer fills meanwhile and Publish starts refusing. When ctx ends it stops
// accepting, sends what is left for at most flushTimeout and returns the messages it could not
// send, for the caller to store itself.
func (publisher *Publisher) Run(ctx context.Context, flushTimeout time.Duration) [][]byte {
	waits := retryWaits{max: publisher.settings.ReconnectMax}
	var connection *amqp.Connection
	var channel *amqp.Channel
	for attempt := 0; ctx.Err() == nil; attempt++ {
		if attempt > 0 {
			publisher.reconnects.Add(1)
		}
		var err error
		connection, channel, err = publisher.connect()
		if err != nil {
			wait := waits.next()
			publisher.state.set(false, err)
			publisher.log(slog.LevelWarn, "queue: publisher cannot reach the broker, retrying", "retryIn", wait.String(),
				"buffered", len(publisher.buffer), "error", err)
			sleep(ctx, wait)
			continue
		}
		waits.reset()
		publisher.state.set(true, nil)
		publisher.log(slog.LevelInfo, "queue: publisher connected", "url", SafeURL(publisher.settings.URL),
			"exchange", publisher.settings.Exchange, "queue", publisher.settings.Queue)
		err = publisher.sendUntilLost(ctx, connection, channel)
		if ctx.Err() != nil {
			break
		}
		publisher.state.set(false, err)
		publisher.log(slog.LevelWarn, "queue: publisher lost the broker, reconnecting", "error", err)
		connection.Close()
		connection, channel = nil, nil
	}
	return publisher.stop(connection, channel, flushTimeout)
}

// connect opens a connection and a channel in confirm mode, with the topology declared.
func (publisher *Publisher) connect() (*amqp.Connection, *amqp.Channel, error) {
	connection, channel, err := openChannel(publisher.settings, "publisher")
	if err != nil {
		return nil, nil, err
	}
	if err := channel.Confirm(false); err != nil {
		connection.Close()
		return nil, nil, fmt.Errorf("confirm mode: %w", err)
	}
	return connection, channel, nil
}

// sendUntilLost sends the buffer group by group until ctx ends or the connection fails.
func (publisher *Publisher) sendUntilLost(ctx context.Context, connection *amqp.Connection, channel *amqp.Channel) error {
	closed := connection.NotifyClose(make(chan *amqp.Error, 1))
	for {
		if len(publisher.retry) == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case closeErr := <-closed:
				return fmt.Errorf("connection closed: %v", closeErr)
			case body := <-publisher.buffer:
				publisher.retry = append(publisher.retry, body)
			}
		}
		publisher.takeBuffered()
		err := publisher.sendGroup(ctx, channel)
		if errors.Is(err, errPublishNacked) {
			// A full queue refuses everything; wait instead of sending it again at once.
			sleep(ctx, firstRetryWait)
			continue
		}
		if err != nil {
			return err
		}
	}
}

// takeBuffered moves buffered messages into the group being sent, up to publishGroupSize.
func (publisher *Publisher) takeBuffered() {
	for len(publisher.retry) < publishGroupSize {
		select {
		case body := <-publisher.buffer:
			publisher.retry = append(publisher.retry, body)
		default:
			return
		}
	}
}

// sendGroup publishes the group (publisher.retry) and waits for its confirms; the messages the
// broker did not confirm stay in publisher.retry for the next try.
func (publisher *Publisher) sendGroup(ctx context.Context, channel *amqp.Channel) error {
	confirmCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), confirmTimeout)
	defer cancel()
	confirmations := make([]*amqp.DeferredConfirmation, 0, len(publisher.retry))
	var sendErr error
	for _, body := range publisher.retry {
		confirmation, err := channel.PublishWithDeferredConfirmWithContext(confirmCtx, publisher.settings.Exchange,
			publisher.settings.RoutingKey, false, false, amqp.Publishing{
				ContentType: "application/json", DeliveryMode: amqp.Persistent, Timestamp: time.Now(), Body: body,
			})
		if err != nil {
			sendErr = fmt.Errorf("publish: %w", err)
			break
		}
		publisher.published.Add(1)
		confirmations = append(confirmations, confirmation)
	}

	var unconfirmed [][]byte
	for index, confirmation := range confirmations {
		acked, err := confirmation.WaitContext(confirmCtx)
		if err != nil || !acked {
			publisher.failed.Add(1)
			unconfirmed = append(unconfirmed, publisher.retry[index])
			if err != nil && sendErr == nil {
				sendErr = fmt.Errorf("wait for confirm: %w", err)
			}
			continue
		}
		publisher.confirmed.Add(1)
	}
	// Messages never sent (a send failed part way) are kept too, after the unconfirmed ones.
	unconfirmed = append(unconfirmed, publisher.retry[len(confirmations):]...)
	publisher.retry = unconfirmed
	if sendErr != nil {
		return sendErr
	}
	if len(unconfirmed) > 0 {
		publisher.log(slog.LevelWarn, "queue: broker refused messages, publishing them again", "count", len(unconfirmed),
			"hint", "check the broker's disk and memory alarms in the management UI")
		return errPublishNacked
	}
	return nil
}

// stop closes the door to Publish, sends what is left within flushTimeout when still connected,
// closes the connection and returns the messages it could not send.
func (publisher *Publisher) stop(connection *amqp.Connection, channel *amqp.Channel, flushTimeout time.Duration) [][]byte {
	publisher.acceptMu.Lock()
	publisher.accepting = false
	publisher.acceptMu.Unlock()

	if channel != nil {
		flushCtx, cancel := context.WithTimeout(context.Background(), flushTimeout)
		defer cancel()
		for flushCtx.Err() == nil && (len(publisher.retry) > 0 || len(publisher.buffer) > 0) {
			publisher.takeBuffered()
			if err := publisher.sendGroup(flushCtx, channel); err != nil && !errors.Is(err, errPublishNacked) {
				break
			}
		}
	}
	if connection != nil {
		connection.Close()
	}
	publisher.state.set(false, nil)

	unsent := publisher.retry
	publisher.retry = nil
	for len(publisher.buffer) > 0 {
		unsent = append(unsent, <-publisher.buffer)
	}
	publisher.log(slog.LevelInfo, "queue: publisher stopped", "unsent", len(unsent))
	return unsent
}
