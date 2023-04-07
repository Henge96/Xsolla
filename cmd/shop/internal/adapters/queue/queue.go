package queue

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/sync/errgroup"
	"time"
	pb "xsolla/api/shop/v1"
	"xsolla/cmd/shop/internal/app"
	"xsolla/internal/dom"
	"xsolla/libs/queue"
)

var _ app.Queue = &Client{}

type (
	// Config provide connection info for message broker.
	Config struct {
		URLs        []string
		Username    string
		Password    string
		ClusterMode bool
	}
	// Client provided data from and to message broker.
	// todo add metrics
	Client struct {
		consumerName        string
		nats                *queue.NATS
		chUpdateOrderStatus chan dom.Event[app.EventUpdateOrderStatusFromQueue]
	}
)

// New build and returns new queue instance.
func New(ctx context.Context, namespace string, cfg Config) (*Client, error) {
	//todo add real connect
	//client, err := queue.ConnectNATS(ctx, strings.Join(cfg.URLs, ","), namespace, cfg.Username, cfg.Password)
	//if err != nil {
	//	return nil, fmt.Errorf("queue.ConnectNATS: %w", err)
	//}
	//
	//err = pb.Migrate(client.JetStream)
	//if err != nil {
	//	return nil, fmt.Errorf("post.Migrate: %w", err)
	//}

	return &Client{
		chUpdateOrderStatus: make(chan dom.Event[app.EventUpdateOrderStatusFromQueue]),
		//nats: client,
	}, nil
}

func (c *Client) AddOrder(ctx context.Context, eventAdd app.EventAddOrder) error {
	// just example about publishing
	_ = c.nats.Publish(ctx, pb.TopicAdd, eventAdd.TaskID, eventAdd)
	return nil
}

func (c *Client) UpdateOrder(ctx context.Context, eventUpdate app.EventUpdateOrder) error {
	// just example about publishing
	_ = c.nats.Publish(ctx, pb.TopicUpdate, eventUpdate.TaskID, eventUpdate)
	return nil
}

func (c *Client) UpdateOrderStatus() <-chan dom.Event[app.EventUpdateOrderStatusFromQueue] {
	return c.chUpdateOrderStatus
}

// Process starts worker for collecting events from queue.
func (c *Client) Process(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	subjects := []string{
		//TODO events from another services
		"EventUpdateOrderStatus",
	}

	for i := range subjects {
		// todo add after implement sub logic
		_ = i

		group.Go(func() error {
			//return c.nats.Subscribe(ctx, subjects[i], c.consumerName, c.handleEvent)

			for {
				select {
				case <- ctx.Done():
					return ctx.Err()
				case _ = <- time.Tick(1 * time.Hour):
					continue
				}
			}
		})
	}

	return group.Wait()
}

func (c *Client) Close() error {
	return c.nats.Drain()
}

func (c *Client) handleEvent(ctx context.Context, msg queue.Message) error {
	ack := make(chan dom.AcknowledgeKind)

	var err error
	switch {
	case ctx.Err() != nil:
		return nil
		//TODO events from another services
	case msg.Subject() == "EventUpdateOrderStatus":
		err = c.handleUpdateOrderStatus(ctx, ack, msg.ID(), msg)
	default:
		err = fmt.Errorf("%w: unknown topic %s", app.ErrInvalidArgument, msg.Subject())
	}

	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return nil
	case ackKind := <-ack:
		switch ackKind {
		case dom.AcknowledgeKindAck:
			err = msg.Ack(ctx)
		case dom.AcknowledgeKindNack:
			err = msg.Nack(ctx)
		}
		if err != nil {
			return fmt.Errorf("msg.Ack|Nack: %w", err)
		}
	}

	return nil
}

func (c *Client) handleUpdateOrderStatus(ctx context.Context, ack chan dom.AcknowledgeKind, msgID uuid.UUID, msg queue.Message) error {
	// todo service structure + convert logic
	event := struct{}{}
	err := msg.Unmarshal(&event)
	if err != nil {
		return fmt.Errorf("msg.Unmarshal: %w", err)
	}

	arg := dom.NewEvent(msgID, ack, app.EventUpdateOrderStatusFromQueue{
		SourceID:        uuid.Must(uuid.NewV4()),
		Status:          dom.OrderStatusCooking,
		SourceCreatedAt: time.Now(),
	})

	select {
	case <-ctx.Done():
		return nil
	case c.chUpdateOrderStatus <- *arg:
	}

	return nil
}
