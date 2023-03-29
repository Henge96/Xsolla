package app

import (
	"github.com/gofrs/uuid/v5"
	"time"
	"xsolla/internal/dom"
)

type (
	Order struct {
		ID        uuid.UUID
		Address   Address
		Items     []Item
		Status    dom.OrderStatus
		Comment   string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	Item struct {
		ID uuid.UUID
		// must be enum
		Type string
		// must be unique
		Name  string
		Count uint16
		// etc
	}

	Address struct {
		City     string
		Street   string
		House    string
		Entrance string
		Flat     string
		// etc
	}

	OrderParams struct {
		Limit  uint16
		Offset uint16
	}

	EventAddOrder struct {
		ID        uuid.UUID
		Status    dom.OrderStatus
		CreatedAt time.Time
	}

	EventUpdateOrder struct {
		ID        uuid.UUID
		Status    dom.OrderStatus
		CreatedAt time.Time
	}

	Task struct {
		ID         uuid.UUID
		Order      Order
		TaskKind   dom.TaskKind
		CreatedAt  time.Time
		UpdatedAt  time.Time
		FinishedAt time.Time
	}
)
