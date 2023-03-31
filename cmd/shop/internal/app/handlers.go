package app

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"xsolla/internal/dom"
)

func (a *App) CreateOrder(ctx context.Context, order Order) (id uuid.UUID, err error) {
	err = a.repo.Tx(ctx, func(repo Repo) error {
		order.Status = dom.OrderStatusNew
		o, err := a.repo.SaveOrder(ctx, order)
		if err != nil {
			return fmt.Errorf("a.repo.SaveOrder: %w", err)
		}

		_, err = a.repo.SaveTask(ctx, Task{
			Order:    *o,
			TaskKind: dom.TaskKindEventAdd,
		})
		if err != nil {
			return fmt.Errorf("a.repo.SaveTask: %w", err)
		}

		id = o.ID
		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (a *App) ListOrders(ctx context.Context, params OrderParams) ([]Order, int, error) {
	return a.repo.ListOrders(ctx, params)
}

func (a *App) ChangeOrderStatus(ctx context.Context, orderID uuid.UUID, status dom.OrderStatus) error {
	// todo validate status
	order, err := a.repo.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("a.repo.GetOrder: %w", err)
	}

	if order.Status == status {
		return ErrSameStatus
	}

	err = a.repo.Tx(ctx, func(repo Repo) error {
		updatedOrder, err := a.repo.UpdateOrder(ctx, Order{
			Status: status,
		})
		if err != nil {
			return fmt.Errorf("a.repo.UpdateOrder: %w", err)
		}

		_, err = a.repo.SaveTask(ctx, Task{
			Order:    *updatedOrder,
			TaskKind: dom.TaskKindEventUpdate,
		})
		if err != nil {
			return fmt.Errorf("a.repo.SaveTask: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *App) GetOrder(ctx context.Context, id uuid.UUID) (*Order, error) {
	order, err := a.repo.GetOrder(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("a.repo.GetOrder: %w", err)
	}

	return order, nil
}
