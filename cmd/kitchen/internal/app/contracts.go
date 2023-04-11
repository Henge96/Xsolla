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
		// SaveCooking in repo.
		// Errors: unknown.
		SaveCooking(context.Context, Cooking) (*Cooking, error)
		// UpdateCooking order in repo.
		// Errors: unknown.
		UpdateCooking(context.Context, Cooking) (*Cooking, error)
		// UpdateCookingStatusByOrderID in repo.
		// Errors: unknown.
		UpdateCookingStatusByOrderID(ctx context.Context, orderID uuid.UUID, status CookingStatus) (*Cooking, error)
		// GetCooking by id.
		// Errors: ErrNotFound, unknown.
		GetCooking(context.Context, uuid.UUID) (*Cooking, error)
		// ListCooking returns cooking by params.
		// Errors: unknown.
		ListCooking(context.Context, CookingParams) ([]Cooking, int, error)
		// SaveOrder adds new order to the repository.
		// Errors: unknown.
		SaveOrder(context.Context, Order) (*Order, error)
		// UpdateOrder order in repo.
		// Errors: unknown.
		UpdateOrder(context.Context, Order) (*Order, error)
		// GetOrder order by id.
		// Errors: ErrNotFound, unknown.
		GetOrder(context.Context, uuid.UUID) (*Order, error)
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
		// UpdateOrder sends event 'EventUpdateOrder' to queue.
		// Errors: unknown.
		UpdateOrder(ctx context.Context, eventUpdate EventUpdateOrder) error
		// UpdateOrderStatus gets EventUpdateOrderStatusFromQueue from queue.
		// Errors: unknown.
		UpdateOrderStatus() <-chan dom.Event[EventUpdateOrderStatusFromQueue]
		// AddOrder gets EventAddOrderFromQueue from queue.
		// Errors: unknown.
		AddOrder() <-chan dom.Event[EventAddOrderFromQueue]
	}

	// Cron for gets jobs.
	Cron interface {
		// Fetch get jobs with limit.
		// Errors: unknown.
		Fetch() <-chan uint16
	}
)
