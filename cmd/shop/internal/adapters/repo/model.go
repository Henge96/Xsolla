package repo

import "github.com/gofrs/uuid/v5"

type (
	order struct {
		ID uuid.UUID `db:"id"`
	}

	task struct {
		ID uuid.UUID `id:"id"`
	}
)
