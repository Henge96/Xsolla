package repo

import (
	"database/sql"
	"errors"
	"xsolla/cmd/shop/internal/app"
)

func convertError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return app.ErrNotFound
	default:
		return err
	}
}
