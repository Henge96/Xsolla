package adapters

import (
	"context"
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

