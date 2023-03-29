package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xsolla/cmd/shop/internal/app"
)

func (a *api) MakeOrder(w http.ResponseWriter, r *http.Request) {
	id, err := a.make(w, r)
	if err != nil {
		errHandler(w, err, statusFromError(err))
		return
	}

	response := responseMakeOrder{
		ID: id.String(),
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		errHandler(w, err, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(bytes)
	if err != nil {
		errHandler(w, err, http.StatusInternalServerError)
		return
	}
}

func (a *api) ListOfOrders(w http.ResponseWriter, r *http.Request) {
	_, _, err := a.list(w, r)
	if err != nil {
		errHandler(w, err, statusFromError(err))
		return
	}

	response := responseListOrders{}
	bytes, err := json.Marshal(response)
	if err != nil {
		errHandler(w, err, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(bytes)
	if err != nil {
		errHandler(w, err, http.StatusInternalServerError)
		return
	}
}

func (a *api) make(_ http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	request, err := validate()
	if err != nil {
		return uuid.Nil, err
	}

	appItems := make([]app.Item, len(request.Items))
	for i := range request.Items {
		appItems[i] = request.Items[i].convert()
	}

	id, err := a.app.CreateOrder(r.Context(), app.Order{
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

func (a *api) list(_ http.ResponseWriter, r *http.Request) ([]app.Order, int, error) {
	vars := mux.Vars(r)

	limit, err := strconv.Atoi(vars["limit"])
	if err != nil || limit < 0 {
		return nil, 0, app.ErrInvalidArgument
	}

	offset, err := strconv.Atoi(vars["offset"])
	if err != nil || offset < 0 {
		return nil, 0, app.ErrInvalidArgument
	}

	res, total, err := a.app.ListOrders(r.Context(), app.OrderParams{
		Limit:  uint16(limit),
		Offset: uint16(offset),
	})

	return res, total, nil
}

func statusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	code := http.StatusInternalServerError
	switch {
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

func validate() (*requestMakeOrder, error) {
	return &requestMakeOrder{}, nil
}
