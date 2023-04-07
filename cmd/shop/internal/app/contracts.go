package app

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"xsolla/internal/dom"
)

type (
	// Repo interface for repo data repository.
	Repo interface {
		TaskRepo
		// Tx transaction for repo methods.
		// Errors: unknown.
		Tx(ctx context.Context, f func(Repo) error) error
		// SaveOrder adds new order to the repository.
		// Errors: unknown.
		SaveOrder(context.Context, Order) (*Order, error)
		// UpdateOrder order in repo.
		// Errors: unknown.
		UpdateOrder(context.Context, Order) (*Order, error)
		// GetOrder order by id.
		// Errors: ErrNotFound, unknown.
		GetOrder(context.Context, uuid.UUID) (*Order, error)
		// ListOrders returns orders by params.
		// Errors: unknown.
		ListOrders(context.Context, OrderParams) ([]Order, int, error)
		// ListProducts returns products by kind and name.
		// Errors: unknown.
		ListProducts(ctx context.Context, types, names []string) ([]Product, error)
		// SaveItem adds item in repository.
		// Errors. unknown.
		SaveItem(context.Context, Item) (*Item, error)
	}

	// TaskRepo interface for saving orders.
	TaskRepo interface {
		// SaveTask adds new task to repository.
		// Errors: unknown.
		SaveTask(context.Context, Task) (uuid.UUID, error)
		// FinishTask set column Task.FinishedAt task.
		// Errors: ErrNotFound, unknown.
		FinishTask(context.Context, uuid.UUID) error
		// ListActualTask returns list task by limit and ordered by created_at (ask).
		// Return tasks without Task.FinishedAt.
		// Errors: unknown.
		ListActualTask(context.Context, int) ([]Task, error)
	}

	// Queue sends events to queue.
	Queue interface {
		// AddOrder sends event 'EventAddOrder' to queue.
		// Errors: unknown.
		AddOrder(context.Context, EventAddOrder) error
		// UpdateOrder sends event 'EventUpdateOrder' to queue.
		// Errors: unknown.
		UpdateOrder(ctx context.Context, eventUpdate EventUpdateOrder) error
		// UpdateOrderStatus gets EventUpdateOrderStatusFromQueue from queue.
		// Errors: unknown.
		UpdateOrderStatus() <-chan dom.Event[EventUpdateOrderStatusFromQueue]
	}

	Cron interface {
		// Fetch get jobs with limit.
		// Errors: unknown.
		Fetch() <-chan uint16
	}

	AddressValidator interface {
		// CheckAddress checks that address exist.
		// Errors: ErrNotFound, unknown.
		CheckAddress(ctx context.Context, a Address) error
	}
)
