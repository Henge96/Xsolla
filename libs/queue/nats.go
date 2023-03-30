package queue

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"xsolla/libs/logkey"
)

const (
	msgIDHeader   = `Nats-Msg-Id`
	drainTimeout  = 3 * time.Second // Should be less than main.shutdownDelay.
	maxReconnects = 5
	pingInterval  = time.Second // Default 2 min isn't useful because TCP keepalive is faster.
)

var _ Message = &natsMessage{}

type natsMessage struct {
	decoder Decoder
	msg     *nats.Msg
}

// ID implements Message.
func (m *natsMessage) ID() uuid.UUID {
	return uuid.Must(uuid.FromString(m.msg.Header.Get(msgIDHeader)))
}

// Unmarshal implements Message.
func (m *natsMessage) Unmarshal(a any) error {
	err := m.decoder.Unmarshal(m.msg.Data, a)
	if err != nil {
		return fmt.Errorf("m.decoder.Unmarshal: %w", err)
	}

	event, ok := a.(Validator)
	if !ok {
		return nil
	}

	err = event.ValidateAll()
	if err != nil {
		return fmt.Errorf("event.ValidateAll: %w", err)
	}

	return nil
}

// Subject implements Message.
func (m *natsMessage) Subject() string {
	return m.msg.Subject
}

// Ack implements Message.
func (m *natsMessage) Ack(ctx context.Context) error {
	return m.msg.Ack(nats.Context(ctx))
}

// Nack implements Message.
func (m *natsMessage) Nack(ctx context.Context) error {
	return m.msg.Nak(nats.Context(ctx))
}

// AsyncErrMsg is error from async publish handler.
type AsyncErrMsg struct {
	Msg *nats.Msg
	Err error
}

// NATS adds connection monitoring to nats.Conn.
type NATS struct {
	Conn            *nats.Conn
	JetStream       nats.JetStreamContext
	closed          chan struct{}
	asyncErrHandler chan AsyncErrMsg // Non-blocking on send, closes by NATS.Close.
	encoder         Encoder
	decoder         Decoder
}

// ConnectNATS adds ctx support and reasonable defaults to nats.Connect.
func ConnectNATS(ctx context.Context, urls, namespace, username, password string) (*NATS, error) {
	c := &NATS{
		closed: make(chan struct{}),
		encoder: &encoderProto{
			MarshalOptions: &proto.MarshalOptions{},
		},
		decoder: &decoderProto{
			UnmarshalOptions: &proto.UnmarshalOptions{},
		},
	}

	var err error
	for !(c.Conn != nil && err == nil) {
		errc := make(chan error)
		go func() {
			err := c.connect(ctx, urls, namespace, username, password)
			select {
			case errc <- err:
			case <-ctx.Done():
				if c.Conn != nil {
					c.Conn.Close()
				}
			}
		}()
		select {
		case err = <-errc:
			fmt.Println("couldn't connect to NATS", err)
		case <-ctx.Done():
			if err == nil {
				err = ctx.Err()
			}

			return nil, err
		}
	}

	fmt.Println("NATS connected", logkey.URL, c.Conn.ConnectedUrl())

	return c, nil
}

func (c *NATS) connect(ctx context.Context, urls, namespace, username, password string) (err error) {
	c.Conn, err = nats.Connect(urls,
		nats.Name(namespace),
		nats.UserInfo(username, password),
		nats.MaxReconnects(maxReconnects),
		nats.DrainTimeout(drainTimeout),
		nats.PingInterval(pingInterval),

		nats.NoCallbacksAfterClientClose(),
		nats.ClosedHandler(func(_ *nats.Conn) {
			close(c.closed)
		}),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err == nil {
				log.Println("NATS disconnected")
			} else {
				log.Println("NATS disconnected", zap.Error(err))
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Println("NATS reconnected", logkey.URL, nc.ConnectedUrl())
		}),
		nats.ErrorHandler(func(_ *nats.Conn, sub *nats.Subscription, err error) {
			if sub == nil {
				log.Println("NATS connection failed", err)
			} else {
				log.Println("NATS connection failed", logkey.Reason, sub.Subject, err)
			}
		}),
	)
	if err != nil {
		return fmt.Errorf("nats.Connect: %w", err)
	}

	c.JetStream, err = c.Conn.JetStream(
		nats.Context(ctx),
		nats.PublishAsyncErrHandler(func(_ nats.JetStream, msg *nats.Msg, err error) {
			c.asyncErrHandler <- AsyncErrMsg{
				Msg: msg,
				Err: err,
			}
		}),
	)
	if err != nil {
		return fmt.Errorf("c.conn.JetStream: %w", err)
	}

	return nil
}

// Monitor waits until ctx.Done or failure reconnecting NATS.
func (c *NATS) Monitor(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case <-c.closed:
		return nats.ErrConnectionClosed
	}
}

// Err returns channel for handling async error.
func (c *NATS) Err() <-chan AsyncErrMsg {
	return c.asyncErrHandler
}

// Drain starts to close process.
func (c *NATS) Drain() error {
	return c.Conn.Drain()
}

// Migrate starts callback with JetStream connection for making streams/consumers.
func (c *NATS) Migrate(f func(manager nats.JetStreamManager) error) error {
	return f(c.JetStream)
}

// Subscribe starts subscription by args.
func (c *NATS) Subscribe(
	ctx context.Context,
	subj, consumerName string,
	handler func(context.Context, Message) error,
) error {
	log := ctxzap.Extract(ctx)

	sub, err := c.JetStream.PullSubscribe(subj, consumerName, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("c.nats.JetStream.QueueSubscribe: %w", err)
	}
	defer func() {
		err := sub.Drain()
		if err != nil {
			log.Error("couldn't drain sub", zap.Error(err))
		}
	}()

	for {
		msgs, err := sub.Fetch(1, nats.Context(ctx))
		switch {
		case ctx.Err() != nil:
			return nil
		case errors.Is(err, context.DeadlineExceeded):
			continue // Because fetcher can return context error by default timeout.
		case err != nil:
			return fmt.Errorf("sub.Fetch: %w", err)
		}

		for i := range msgs {
			if ctx.Err() != nil {
				return nil
			}

			err = handler(ctx, &natsMessage{decoder: c.decoder, msg: msgs[i]})
			if err != nil {
				log.Error("couldn't handle message", zap.Error(err))
			}
		}
	}
}

// Publish send message to queue.
func (c *NATS) Publish(ctx context.Context, topic string, msgID uuid.UUID, event any) error {
	if event, ok := event.(Validator); ok {
		err := event.ValidateAll()
		if err != nil {
			return fmt.Errorf("event.ValidateAll: %w", err)
		}
	}

	buf, err := c.encoder.Marshal(event)
	if err != nil {
		return fmt.Errorf("c.encoder.Marshal: %w", err)
	}

	_, err = c.JetStream.Publish(
		topic,
		buf,
		nats.MsgId(msgID.String()),
		nats.Context(ctx),
	)
	if err != nil {
		return fmt.Errorf("c.JetStream.Publish: %w", err)
	}

	return nil
}
