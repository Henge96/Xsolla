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

func (tx *txRepo) ListProducts(ctx context.Context, items []app.Item) ([]app.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) SaveTask(ctx context.Context, t2 app.Task) (uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) FinishTask(ctx context.Context, uuid uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) ListActualTask(ctx context.Context, i int) ([]app.Task, error) {
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

func (tx *txRepo) ListOrders(ctx context.Context, params app.OrderParams) ([]app.Order, int, error) {
	//TODO implement me
	panic("implement me")
}

func (tx *txRepo) Tx(ctx context.Context, f func(app.Repo) error) error {
	panic("couldn`t start transaction into transaction")
}
