package repo

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"xsolla/cmd/kitchen/internal/app"
)

var _ app.Repo = &txRepo{}

type txRepo struct {
	tx *sqlx.Tx
}

func (tx *txRepo) SaveTask(ctx context.Context, task app.Task) (uuid.UUID, error) {
	ts, err := convertTask(task)
	if err != nil {
		return uuid.Nil, fmt.Errorf("convertTask: %w", err)
	}

	const query = `insert into tasks (kind, order_bytes) values ($1, $2) returning id`
	var id uuid.UUID
	err = tx.tx.GetContext(ctx, &id, query, ts.Kind, ts.OrderBytes)
	if err != nil {
		return uuid.Nil, fmt.Errorf("db.GetContext: %w", convertError(err))
	}

	return id, nil
}

func (tx *txRepo) FinishTask(ctx context.Context, id uuid.UUID) error {
	const query = `update tasks 
        set updated_at = now(),
        finished_at = now() 
        where id = $1 returning *`

	err := tx.tx.GetContext(ctx, &task{}, query, id)
	if err != nil {
		return fmt.Errorf("db.GetContext: %w", err)
	}

	return nil

}

func (tx *txRepo) ListActualTask(ctx context.Context, limit int) ([]app.Task, error) {
	const query = `select * from tasks where finished_at is null order by created_at asc limit $1 for update`
	res := make([]task, 0, limit)
	err := tx.tx.SelectContext(ctx, &res, query, limit)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", convertError(err))
	}

	tasks := make([]app.Task, 0, len(res))
	for i := range res {
		t, err := res[i].convert()
		if err != nil {
			return nil, fmt.Errorf("res.convert: %w", err)
		}

		tasks = append(tasks, *t)
	}

	return tasks, nil

}

func (tx *txRepo) Tx(ctx context.Context, f func(app.Repo) error) error {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) SaveCooking(ctx context.Context, cooking app.Cooking) (*app.Cooking, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) UpdateCooking(ctx context.Context, cooking app.Cooking) (*app.Cooking, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) UpdateCookingStatusByOrderID(ctx context.Context, orderID uuid.UUID, status app.CookingStatus) (*app.Cooking, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) GetCooking(ctx context.Context, uuid uuid.UUID) (*app.Cooking, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) ListCooking(ctx context.Context, params app.CookingParams) ([]app.Cooking, int, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) SaveOrder(ctx context.Context, o app.Order) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) UpdateOrder(ctx context.Context, o app.Order) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) GetOrder(ctx context.Context, uuid uuid.UUID) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

