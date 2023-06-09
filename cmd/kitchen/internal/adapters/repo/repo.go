package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"xsolla/cmd/kitchen/internal/app"
)

var _ app.Repo = &Repo{}

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	// metrics, migrations ...
	return &Repo{db: db}
}

func (r *Repo) Close() error {
	return r.db.Close()
}

func (r *Repo) UpdateCookingStatusByOrderID(ctx context.Context, orderID uuid.UUID, status app.CookingStatus) (*app.Cooking, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) SaveOrder(ctx context.Context, order app.Order) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) UpdateOrder(ctx context.Context, order app.Order) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) GetOrder(ctx context.Context, u uuid.UUID) (*app.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) SaveTask(ctx context.Context, task app.Task) (uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) FinishTask(ctx context.Context, u uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) ListActualTask(ctx context.Context, i int) ([]app.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) SaveCooking(ctx context.Context, cooking app.Cooking) (*app.Cooking, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) UpdateCooking(ctx context.Context, cooking app.Cooking) (*app.Cooking, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) GetCooking(ctx context.Context, uuid uuid.UUID) (*app.Cooking, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) ListCooking(ctx context.Context, params app.CookingParams) ([]app.Cooking, int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) Tx(ctx context.Context, f func(app.Repo) error) error {
	opts := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	return txHelper(ctx, r.db, opts, func(tx *sqlx.Tx) error {
		return f(&txRepo{tx: tx})
	})
}

func txHelper(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions, cb func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("db.BeginTx: %w", err)
	}

	err = cb(tx)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			err = fmt.Errorf("%w: %s", err, errRollback)
		}
		return err
	}

	return tx.Commit()
}
