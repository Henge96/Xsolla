package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xsolla/cmd/shop/internal/app"
)

func (a *api) MakeOrder(w http.ResponseWriter, r *http.Request) {
	request, err := validateRequestMakeOrder(r)
	if err != nil {
		errHandler(w, err, statusFromError(err))
		return
	}

	id, err := a.make(r.Context(), *request)
	if err != nil {
		errHandler(w, err, statusFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(responseMakeOrder{
		ID: id.String(),
	})
	if err != nil {
		errHandler(w, err, statusFromError(err))
		return
	}
}

func (a *api) ListOfOrders(w http.ResponseWriter, r *http.Request) {
	request, err := validateRequestListOfOrders(r)
	if err != nil {
		errHandler(w, err, statusFromError(err))
		return
	}

	_, _, err = a.app.ListOrders(r.Context(), app.OrderParams{
		Limit:  uint16(request.Limit),
		Offset: uint16(request.Offset),
	})

	// todo convert to struct
	err = json.NewEncoder(w).Encode(responseListOrders{})
	if err != nil {
		errHandler(w, err, statusFromError(err))
		return
	}
}

func (a *api) ChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	request, err := validateRequestChangeOrderStatus(r)
	if err != nil {
		errHandler(w, err, statusFromError(err))
		return
	}

	err = a.app.ChangeOrderStatus(r.Context(), uuid.FromStringOrNil(request.ID), request.Status)
	if err != nil {
		errHandler(w, err, statusFromError(err))
		return
	}

	return
}

func statusFromError(err error) int {
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

func errHandler(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

func (a *api) make(ctx context.Context, request requestMakeOrder) (uuid.UUID, error) {
	appItems := make([]app.Item, len(request.Items))
	for i := range request.Items {
		appItems[i] = request.Items[i].convert()
	}

	id, err := a.app.CreateOrder(ctx, app.Order{
		Address: app.Address{
			City:     request.Address.City,
			Street:   request.Address.Street,
			House:    request.Address.House,
			Entrance: request.Address.Entrance,
			Flat:     request.Address.Flat,
		},
		Items:   appItems,
		Comment: request.Comment,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// todo validate logic
func validateRequestMakeOrder(_ *http.Request) (*requestMakeOrder, error) {
	return &requestMakeOrder{}, nil
}

// todo validate logic
func validateRequestChangeOrderStatus(_ *http.Request) (*requestUpdateOrderStatus, error) {
	return &requestUpdateOrderStatus{}, nil
}

func validateRequestListOfOrders(r *http.Request) (*requestListOrders, error) {
	vars := mux.Vars(r)

	limit, err := strconv.Atoi(vars["limit"])
	if err != nil || limit < 0 {
		return nil, fmt.Errorf("wrong limit: %w", app.ErrInvalidArgument)
	}

	offset, err := strconv.Atoi(vars["offset"])
	if err != nil || offset < 0 {
		return nil, fmt.Errorf("wrong offset: %w", app.ErrInvalidArgument)
	}

	// todo
	return &requestListOrders{
		Limit:  limit,
		Offset: offset,
	}, nil
}
