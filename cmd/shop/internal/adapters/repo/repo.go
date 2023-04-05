package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"xsolla/cmd/shop/internal/app"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	// metrics, migrations ...
	return &Repo{db: db}
}

func (r *Repo) SaveOrder(ctx context.Context, order app.Order) (*app.Order, error) {
	repoOrder, err := convertToOrder(order)
	if err != nil {
		return nil, fmt.Errorf("convertToOrder: %w", err)
	}

	query := `insert into orders (address, items, status, comment) values ($1, $2, $3, $4) returning *`
	err = r.db.GetContext(ctx, &repoOrder, query, repoOrder.Address, repoOrder.Items, repoOrder.Status, repoOrder.Comment)
	if err != nil {
		return nil, fmt.Errorf("r.db.GetContext: %w", convertError(err))
	}

	o, err := repoOrder.convert()
	if err != nil {
		return nil, fmt.Errorf("repoOrder.convert: %w", convertError(err))
	}

	return o, nil
}

func (r *Repo) UpdateOrder(ctx context.Context, order app.Order) (*app.Order, error) {
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
	ts, err := convertTask(task)
	if err != nil {
		return uuid.Nil, fmt.Errorf("convertTask: %w", err)
	}

	const query = `insert into tasks (kind, order_bytes) values ($1, $2) returning id`
	var id uuid.UUID
	err = r.db.GetContext(ctx, &id, query, ts.Kind, ts.OrderBytes)
	if err != nil {
		return uuid.Nil, fmt.Errorf("db.GetContext: %w", convertError(err))
	}

	return id, nil
}

func (r *Repo) FinishTask(ctx context.Context, id uuid.UUID) error {
	const query = `update tasks 
        set updated_at = now(),
        finished_at = now() 
        where id = $1 returning *`

	err := r.db.GetContext(ctx, &task{}, query, id)
	if err != nil {
		return fmt.Errorf("db.GetContext: %w", err)
	}

	return nil
}

func (r *Repo) ListActualTask(ctx context.Context, limit int) ([]app.Task, error) {
	const query = `select * from tasks where finished_at is null order by created_at asc limit $1`
	res := make([]task, 0, limit)
	err := r.db.SelectContext(ctx, &res, query, limit)
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

func (r *Repo) ListProducts(ctx context.Context, items []app.Item) ([]app.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) SaveItem(context.Context, app.Item) (*app.Item, error) {
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

func (r *Repo) Close() error {
	return r.db.Close()
}
