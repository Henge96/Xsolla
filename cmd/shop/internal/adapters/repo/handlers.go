package repo

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"xsolla/cmd/shop/internal/app"
)

var _ app.Repo = &Repo{}

func (r *Repo) SaveOrder(ctx context.Context, order app.Order) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) UpdateOrder(ctx context.Context, order app.Order) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) GetOrder(ctx context.Context, uuid uuid.UUID) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) ListOrders(ctx context.Context, params app.OrderParams) ([]app.Order, int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) SaveTask(ctx context.Context, task app.Task) (uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) FinishTask(ctx context.Context, uuid uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) ListActualTask(ctx context.Context, i int) ([]app.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) Close() error {
	return r.db.Close()
}
