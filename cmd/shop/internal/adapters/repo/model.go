package repo

import "github.com/gofrs/uuid/v5"

// todo structs fields and converts
type (
	order struct {
		ID uuid.UUID `db:"id"`
	}

	task struct {
		ID uuid.UUID `id:"id"`
	}
)
