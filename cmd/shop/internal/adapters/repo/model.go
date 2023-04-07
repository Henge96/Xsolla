package repo

import (
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgtype"
	"time"
)

// todo structs fields and converts
type (
	order struct {
		ID        uuid.UUID    `db:"id" json:"id"`
		Address   pgtype.JSONB `db:"address" json:"address"`
		Items     pgtype.JSONB `db:"items" json:"items"`
		Status    string       `db:"status" json:"status"`
		Comment   string       `db:"comment" json:"comment"`
		CreatedAt time.Time    `db:"created_at" json:"created_at"`
		UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
	}

	task struct {
		ID         uuid.UUID `id:"id"`
		OrderBytes []byte    `db:"order_bytes"`
		Kind       string    `db:"kind"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
		FinishedAt time.Time `db:"finished_at"`
	}

	address struct {
		City     string `db:"city" json:"city"`
		Street   string `db:"street" json:"street"`
		House    string `db:"house" json:"house"`
		Entrance string `db:"entrance" json:"entrance"`
		Flat     string `db:"flat" json:"flat"`
	}

	item struct {
		ID      uuid.UUID    `db:"id"`
		OrderID uuid.UUID    `db:"order_id"`
		Product pgtype.JSONB `db:"product"`
		Count   uint16       `db:"count"`
		Comment string       `db:"comment"`
	}

	product struct {
		ID   uuid.UUID `db:"id" json:"id"`
		Type string    `db:"type" json:"type"`
		Name string    `db:"name" json:"name"`
	}
)
