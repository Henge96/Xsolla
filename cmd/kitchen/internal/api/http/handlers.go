package http

import (
	"context"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"net/http"
	"xsolla/cmd/kitchen/internal/app"
)

func (a *api) ListOfCooking(w http.ResponseWriter, r *http.Request) {
	// todo handler
	_, _, _ = a.app.ListCooking(r.Context(), app.CookingParams{})
	return
}

func (a *api) ChangeCookingStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := uuid.FromStringOrNil(vars["id"])

	status, err := toAppCookingStatus(vars["status"])
	if err != nil {
		errHandler(w, err, statusCodeFromError(err))

		return
	}

	err = a.app.ChangeCookingStatus(r.Context(), id, status)
	if err != nil {
		errHandler(w, err, statusCodeFromError(err))

		return
	}
	return
}

func errHandler(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

func statusCodeFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, app.ErrSameStatus):
		code = http.StatusAlreadyReported
	case errors.Is(err, app.ErrInvalidArgument):
		code = http.StatusBadRequest
	case errors.Is(err, context.DeadlineExceeded):
		code = http.StatusRequestTimeout
	case errors.Is(err, context.Canceled):
		code = http.StatusRequestTimeout
	}

	return code
}
