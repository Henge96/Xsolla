package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"net/http"
	"xsolla/cmd/kitchen/internal/app"
)

func (a *api) ListOfCooking(w http.ResponseWriter, r *http.Request) {
	params, err := a.validateRequestListOfCooking(r)
	if err != nil {
		errHandler(w, err, statusCodeFromError(err))

		return
	}

	_, _, _ = a.app.ListCooking(r.Context(), *params)

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

func (a *api) validateRequestListOfCooking(r *http.Request) (*app.CookingParams, error) {
	var request requestListOfCooking
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("json.NewDecoder.Decode: %w", err)
	}

	err = a.validator.Struct(&request)
	if err != nil {
		return nil, fmt.Errorf("a.validator.Struct: %w", err)
	}

	status, err := toAppCookingStatus(request.CookingStatus.String())
	if err != nil {
		return nil, fmt.Errorf("toAppCookingStatus: %w", err)
	}

	return &app.CookingParams{
		CookingStatus: status,
		Limit:         request.Limit,
		Offset:        request.Offset,
	}, nil
}
