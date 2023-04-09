package pb

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

// Topics.
const (
	description            = "Events from shop service for notifying about events."
	Stream                 = "kitchen"
	events                 = Stream + ".events"
	version                = events + ".v1."
	TopicUpdateOrderStatus = version + "update"
	SubscribeToAllEvents   = version + "*"
)

const (
	maxMsgReplicas  = 1
	duplicateWindow = time.Second * 30
)

// Migrate for init streams.
func Migrate(js nats.JetStreamManager) error {
	replicas := maxMsgReplicas
	eventStream := &nats.StreamConfig{
		Name:        Stream,
		Description: description,
		Subjects:    []string{TopicUpdateOrderStatus},
		Retention:   nats.LimitsPolicy,
		Storage:     nats.FileStorage,
		Replicas:    replicas,
		NoAck:       false,
		Duplicates:  duplicateWindow,
	}

	_, err := js.AddStream(eventStream)
	switch {
	case errors.Is(err, nats.ErrStreamNameAlreadyInUse):
		_, err = js.UpdateStream(eventStream)
		if err != nil {
			return fmt.Errorf("js.UpdateStream: %w", err)
		}

		return nil
	case err != nil:
		return fmt.Errorf("js.AddStream: %w", err)
	}

	return nil
}
