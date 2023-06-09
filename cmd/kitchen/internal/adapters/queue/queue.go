package queue

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/sync/errgroup"
	"time"
	kitchen_pb "xsolla/api/kitchen/v1"
	shop_pb "xsolla/api/shop/v1"
	"xsolla/cmd/kitchen/internal/app"
	"xsolla/internal/dom"
	"xsolla/libs/queue"
)

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
		chAddOrder          chan dom.Event[app.EventAddOrderFromQueue]
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

// Process starts worker for collecting events from queue.
func (c *Client) Process(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	subjects := []string{
		shop_pb.SubscribeToAllEvents,
	}

	for i := range subjects {
		// todo add after implement sub logic
		_ = i
		group.Go(func() error {
			//return c.nats.Subscribe(ctx, subjects[i], c.consumerName, c.handleEvent)

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case _ = <-time.Tick(1 * time.Hour):
					continue
				}
			}
		})
	}

	return group.Wait()
}

func (c *Client) UpdateCooking(ctx context.Context, eventUpdate app.EventUpdateCooking) error {
	// todo other handler logic
	_ = c.nats.Publish(ctx, kitchen_pb.TopicUpdate, eventUpdate.TaskID, eventUpdate)
	return nil
}

func (c *Client) UpdateOrderStatus() <-chan dom.Event[app.EventUpdateOrderStatusFromQueue] {
	return c.chUpdateOrderStatus
}

func (c *Client) AddOrder() <-chan dom.Event[app.EventAddOrderFromQueue] {
	return c.chAddOrder
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
	case msg.Subject() == shop_pb.TopicUpdate:
		err = c.handleUpdateOrderStatus(ctx, ack, msg.ID(), msg)
	case msg.Subject() == shop_pb.TopicAdd:
		err = c.handleAddOrder(ctx, ack, msg.ID(), msg)
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

func (c *Client) handleAddOrder(ctx context.Context, ack chan dom.AcknowledgeKind, msgID uuid.UUID, msg queue.Message) error {
	// todo service structure + convert logic
	event := struct {
		ID        uuid.UUID       `json:"id,omitempty"`
		Items     []app.Item      `json:"items,omitempty"`
		Status    dom.OrderStatus `json:"status,omitempty"`
		Comment   string          `json:"comment,omitempty"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedAt time.Time       `json:"updated_at"`
	}{}
	err := msg.Unmarshal(&event)
	if err != nil {
		return fmt.Errorf("msg.Unmarshal: %w", err)
	}

	arg := dom.NewEvent(msgID, ack, app.EventAddOrderFromQueue{
		Order: app.Order{
			SourceID:  event.ID,
			Items:     event.Items,
			Status:    event.Status,
			Comment:   event.Comment,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		},
	})

	select {
	case <-ctx.Done():
		return nil
	case c.chAddOrder <- *arg:
	}

	return nil
}

func (c *Client) handleUpdateOrderStatus(ctx context.Context, ack chan dom.AcknowledgeKind, msgID uuid.UUID, msg queue.Message) error {
	// todo service structure + convert logic
	event := struct {
		ID        uuid.UUID       `json:"id,omitempty"`
		Status    dom.OrderStatus `json:"status,omitempty"`
		CreatedAt time.Time       `json:"created_at"`
	}{}
	err := msg.Unmarshal(&event)
	if err != nil {
		return fmt.Errorf("msg.Unmarshal: %w", err)
	}

	arg := dom.NewEvent(msgID, ack, app.EventUpdateOrderStatusFromQueue{
		SourceID:        event.ID,
		Status:          event.Status,
		SourceCreatedAt: event.CreatedAt,
	})

	select {
	case <-ctx.Done():
		return nil
	case c.chUpdateOrderStatus <- *arg:
	}

	return nil
}
