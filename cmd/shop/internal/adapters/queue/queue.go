package queue

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
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
	Client struct {
		nats *queue.NATS
		m    Metrics
	}
)

// New build and returns new queue instance.
func New(ctx context.Context, reg *prometheus.Registry, namespace string, cfg Config) (*Client, error) {
	const subsystem = "queue"
	m := NewMetrics(reg, namespace, subsystem, []string{})

	client, err := queue.ConnectNATS(ctx, strings.Join(cfg.URLs, ","), namespace, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("queue.ConnectNATS: %w", err)
	}

	err = pb.Migrate(client.JetStream)
	if err != nil {
		return nil, fmt.Errorf("post.Migrate: %w", err)
	}

	return &Client{
		nats: client,
		m:    m,
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
