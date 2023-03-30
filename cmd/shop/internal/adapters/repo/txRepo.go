package repo

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"xsolla/cmd/shop/internal/app"
)

var _ app.Repo = &txRepo{}

type txRepo struct {
	tx *sqlx.Tx
}

func (t *txRepo) SaveTask(ctx context.Context, t2 app.Task) (uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (t *txRepo) FinishTask(ctx context.Context, uuid uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (t *txRepo) ListActualTask(ctx context.Context, i int) ([]app.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t *txRepo) SaveOrder(ctx context.Context, o app.Order) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (t *txRepo) UpdateOrder(ctx context.Context, o app.Order) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (t *txRepo) GetOrder(ctx context.Context, uuid uuid.UUID) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (t *txRepo) ListOrders(ctx context.Context, params app.OrderParams) ([]app.Order, int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *txRepo) Tx(ctx context.Context, f func(app.Repo) error) error {
	panic("couldn`t start transaction into transaction")
}
