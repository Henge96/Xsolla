package app

import (
	"github.com/gofrs/uuid/v5"
	"time"
	"xsolla/internal/dom"
)

type (
	Order struct {
		ID        uuid.UUID
		SourceID  uuid.UUID
		Status    dom.OrderStatus
		Comment   string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	UpdateOrder struct {
		// todo add order fields
		ID        uuid.UUID
		Status    dom.OrderStatus
		CreatedAt time.Time
	}

	AddOrder struct {
		Order Order
	}

	Task struct {
		ID         uuid.UUID
		Order      Order
		TaskKind   dom.TaskKind
		CreatedAt  time.Time
		UpdatedAt  time.Time
		FinishedAt time.Time
	}

	// todo rename probably
	EventUpdateOrderStatusFromQueue struct {
		SourceID        uuid.UUID
		Status          dom.OrderStatus
		SourceCreatedAt time.Time
	}
)
