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
		OrderID uuid.UUID
		Product Product
		Count   uint16
		Comment string
	}

	// todo think to replace in another service
	Product struct {
		ID uuid.UUID
		// must be enum (examples in package dom)
		Type string
		// enum
		Name string
		// maybe desc, receipt, etc...
	}

	Address struct {
		// enum
		City     string
		Street   string
		House    string
		Entrance string
		Flat     string
		// maybe postcode, desc etc...
	}

	OrderParams struct {
		Limit  uint16
		Offset uint16
	}

	EventAddOrder struct {
		// todo add order fields
		TaskID    uuid.UUID
		ID        uuid.UUID
		Status    dom.OrderStatus
		CreatedAt time.Time
	}

	EventUpdateOrder struct {
		// todo add order fields
		TaskID    uuid.UUID
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


