package queue

import (
	"context"
	pb "xsolla/api/shop/v1"
	"xsolla/cmd/shop/internal/app"
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
		nats *queue.NATS
	}
)

// New build and returns new queue instance.
func New(ctx context.Context,  namespace string, cfg Config) (*Client, error) {
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

func (c *Client) Close() error {
	return c.nats.Drain()
}
