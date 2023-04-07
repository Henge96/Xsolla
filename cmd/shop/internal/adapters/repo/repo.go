package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (r *Repo) UpdateOrder(ctx context.Context, o app.Order) (*app.Order, error) {
	updateOrder, err := convertToOrder(o)
	if err != nil {
		return nil, fmt.Errorf("convertToOrder: %w", convertError(err))
	}
	const query = `update comments set status = $1, updated_at = now() where id = $2`
	var res order
	err = r.db.GetContext(ctx, &res, query, updateOrder.Status, updateOrder.ID)
	if err != nil {
		return nil, fmt.Errorf("r.db.ExecContext: %w", convertError(err))
	}

	appOrder, err := res.convert()
	if err != nil {
		return nil, fmt.Errorf("res.convert: %w", convertError(err))
	}

	return appOrder, nil
}

func (r *Repo) GetOrder(ctx context.Context, id uuid.UUID) (*app.Order, error) {
	var res order
	const query = `select * from orders where id = $q`
	err := r.db.GetContext(ctx, &res, query, id)
	if err != nil {
		return nil, fmt.Errorf("r.db.GetContext: %w", convertError(err))
	}

	appOrder, err := res.convert()
	if err != nil {
		return nil, fmt.Errorf("res.convert: %w", convertError(err))
	}

	return appOrder, nil
}

func (r *Repo) ListOrders(ctx context.Context, params app.OrderParams) ([]app.Order, int, error) {
	const query = `select * from orders where status = $1 order by created_at desc limit $2 offset $3`

	result := make([]order, 0, int(params.Limit))
	err := r.db.SelectContext(ctx, &result, query, params.Status.String(), params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("db.SelectContext: %w", convertError(err))
	}

	var total int
	const getTotal = `select count(*) as total from orders where status = $1`
	err = r.db.GetContext(ctx, &total, getTotal, params.Status)
	if err != nil {
		return nil, 0, fmt.Errorf("db.GetContext: %w", convertError(err))
	}

	orders := make([]app.Order, 0, len(result))
	for i := range result {
		o, err := result[i].convert()
		if err != nil {
			return nil, 0, fmt.Errorf("result.convert: %w", err)
		}

		orders = append(orders, *o)
	}

	return orders, total, nil
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

func (r *Repo) ListProducts(ctx context.Context, types, names []string) ([]app.Product, error) {
	const query = `select * from orders where type = any($1) and name = any($2)`

	result := make([]product, 0, len(types))
	err := r.db.SelectContext(ctx, &result, query, pq.Array(types), pq.Array(names))
	if err != nil {
		return nil, fmt.Errorf("r.db.SelectContext: %w", convertError(err))
	}

	orders := make([]app.Product, 0, len(result))
	for i := range result {
		o := result[i].convert()
		orders = append(orders, *o)
	}

	return orders, nil
}

func (r *Repo) SaveItem(ctx context.Context, i app.Item) (*app.Item, error) {
	repoItem, err := convertToItem(i)
	if err != nil {
		return nil, fmt.Errorf("convertToItem: %w", convertError(err))
	}

	var res item
	query := `insert into orders (order_id, product, count, comment) values ($1, $2, $3, $4) returning *`
	err = r.db.GetContext(ctx, &res, query, repoItem.OrderID, repoItem.Product, repoItem.Count, repoItem.Comment)
	if err != nil {
		return nil, fmt.Errorf("r.db.GetContext: %w", convertError(err))
	}

	appItem, err := res.convert()
	if err != nil {
		return nil, fmt.Errorf("repoOrder.convert: %w", convertError(err))
	}

	return appItem, nil
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
