package queue

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"7th-sense-server/internal/config"
)

// testSettings returns settings for a broker at RABBITMQ_TEST_URL with an exchange and queues of
// the test's own, or skips the test when the variable is not set.
func testSettings(t *testing.T) config.RabbitMQConfig {
	t.Helper()
	brokerURL := os.Getenv("RABBITMQ_TEST_URL")
	if brokerURL == "" {
		t.Skip("RABBITMQ_TEST_URL not set: no broker to test against")
	}
	name := fmt.Sprintf("test.%d", time.Now().UnixNano())
	settings := config.RabbitMQConfig{URL: brokerURL, Exchange: name, RoutingKey: "values", Queue: name,
		QueueType: config.QueueTypeQuorum, Consumers: 1, Prefetch: 100, PublishBuffer: 1000, ReconnectMax: time.Second}
	t.Cleanup(func() { deleteTopology(t, settings) })
	return settings
}

// deleteTopology removes what the test declared.
func deleteTopology(t *testing.T, settings config.RabbitMQConfig) {
	connection, err := dial(settings, "test cleanup")
	if err != nil {
		t.Logf("cleanup: %v", err)
		return
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		t.Logf("cleanup: %v", err)
		return
	}
	channel.QueueDelete(settings.Queue, false, false, false)
	channel.QueueDelete(DeadLetterQueue(settings), false, false, false)
	channel.ExchangeDelete(settings.Exchange, false, false)
	channel.ExchangeDelete(settings.Exchange+deadLetterExchangeSuffix, false, false)
}

func testLogger(t *testing.T) Logger {
	return func(level slog.Level, message string, attributes ...any) {
		t.Log(append([]any{level.String(), message}, attributes...)...)
	}
}

// TestPublishConsume publishes messageCount messages, consumes them all and checks every one was
// confirmed and arrived once; a rejected message must land in the dead-letter queue.
func TestPublishConsume(t *testing.T) {
	settings := testSettings(t)
	const messageCount = 300

	publisher := NewPublisher(settings, testLogger(t))
	publishCtx, stopPublisher := context.WithCancel(context.Background())
	unsentResult := make(chan [][]byte, 1)
	go func() { unsentResult <- publisher.Run(publishCtx, 5*time.Second) }()

	for index := range messageCount {
		if !publisher.Publish(fmt.Appendf(nil, "message %d", index)) {
			t.Fatalf("Publish %d refused", index)
		}
	}
	publisher.Publish([]byte("dead letter"))

	received := map[string]bool{}
	rejected := false
	done := make(chan struct{})
	consumer := NewConsumer(settings, "test", testLogger(t))
	consumeCtx, stopConsumer := context.WithCancel(context.Background())
	go func() {
		defer close(done)
		consumer.Run(consumeCtx, func(ctx context.Context, session *Session) {
			for {
				select {
				case <-ctx.Done():
					return
				case delivery, ok := <-session.Deliveries():
					if !ok {
						return
					}
					if string(delivery.Body) == "dead letter" {
						session.Reject(delivery.Tag)
						rejected = true
					} else {
						received[string(delivery.Body)] = true
						session.Ack(delivery.Tag, false)
					}
					if len(received) == messageCount && rejected {
						// A quorum queue applies acks asynchronously: closing the connection at
						// once can requeue the last ones (harmless in the server, where writes are
						// idempotent, but this test counts exactly).
						time.AfterFunc(time.Second, stopConsumer)
					}
				}
			}
		})
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		stopConsumer()
		<-done
		t.Fatalf("received %d of %d messages", len(received), messageCount)
	}
	stopPublisher()
	if unsent := <-unsentResult; len(unsent) != 0 {
		t.Errorf("unsent messages: %d", len(unsent))
	}
	stats := publisher.Stats()
	if stats.Confirmed != messageCount+1 {
		t.Errorf("confirmed %d, want %d", stats.Confirmed, messageCount+1)
	}

	// Dead-lettering is asynchronous on quorum queues.
	deadline := time.Now().Add(5 * time.Second)
	for {
		waiting, dead, err := QueueDepth(settings)
		if err != nil {
			t.Fatalf("queue depth: %v", err)
		}
		if waiting == 0 && dead == 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("queue depth: waiting %d (want 0), dead %d (want 1)", waiting, dead)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// TestPublishWithoutBroker checks that a publisher that cannot reach the broker keeps accepting
// until its buffer is full, then refuses, and hands every buffered message back when stopped.
func TestPublishWithoutBroker(t *testing.T) {
	settings := config.RabbitMQConfig{URL: "amqp://guest:guest@127.0.0.1:1/", Exchange: "x", RoutingKey: "y", Queue: "z",
		QueueType: config.QueueTypeQuorum, PublishBuffer: 3, ReconnectMax: time.Second}
	publisher := NewPublisher(settings, testLogger(t))
	ctx, stop := context.WithCancel(context.Background())
	unsentResult := make(chan [][]byte, 1)
	go func() { unsentResult <- publisher.Run(ctx, time.Second) }()

	for index := range 3 {
		if !publisher.Publish([]byte{byte(index)}) {
			t.Fatalf("Publish %d refused before the buffer was full", index)
		}
	}
	if publisher.Publish([]byte{9}) {
		t.Fatal("Publish accepted a message beyond the buffer")
	}
	stop()
	if unsent := <-unsentResult; len(unsent) != 3 {
		t.Fatalf("unsent %d, want 3", len(unsent))
	}
	if publisher.Publish([]byte{10}) {
		t.Fatal("Publish accepted a message after the publisher stopped")
	}
	if stats := publisher.Stats(); stats.Fallback != 2 {
		t.Fatalf("fallback %d, want 2", stats.Fallback)
	}
}
