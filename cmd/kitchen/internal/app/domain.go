package app

import (
	"github.com/gofrs/uuid/v5"
	"time"
	"xsolla/internal/dom"
)

type (
	Cooking struct {
		ID     uuid.UUID
		Status CookingStatus
		// next two fields for example what can be in real project
		Chef      uuid.UUID
		Table     uint32
		Order     Order
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	Order struct {
		ID        uuid.UUID
		SourceID  uuid.UUID
		Items     []Item
		Status    dom.OrderStatus
		Comment   string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	Item struct {
		ID      uuid.UUID
		OrderID uuid.UUID
		Product Product
		Count   uint16
		Comment string
	}

	Product struct {
		// TODO add fields
		ID   uuid.UUID
		Type string
		Name string
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

	CookingParams struct {
		// todo
		Limit  uint16
		Offset uint16
	}

	EventAddOrderFromQueue struct {
		Order Order
	}

	EventUpdateOrderStatusFromQueue struct {
		SourceID        uuid.UUID
		Status          dom.OrderStatus
		SourceCreatedAt time.Time
	}

	EventUpdateOrder struct {
		// todo add order fields
		TaskID    uuid.UUID
		ID        uuid.UUID
		Status    dom.OrderStatus
		CreatedAt time.Time
	}
)

// CookingStatus for cooking in-app.
type CookingStatus uint8

//go:generate stringer -output=stringer.CookingStatus.go -type=CookingStatus -trimprefix=CookingStatus
const (
	_ CookingStatus = iota
	OrderStatusNew
	OrderStatusCompleted
)
