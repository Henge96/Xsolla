package http

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"net/http"
	"xsolla/cmd/shop/internal/app"
	"xsolla/internal/dom"
)

type application interface {
	CreateOrder(ctx context.Context, order app.Order) (uuid.UUID, error)
	ListOrders(ctx context.Context, params app.OrderParams) ([]app.Order, int, error)
	ChangeOrderStatus(ctx context.Context, orderID uuid.UUID, status dom.OrderStatus) error
	GetOrder(ctx context.Context, id uuid.UUID) (*app.Order, error)
}

type api struct {
	app application
}

func New(app application) *mux.Router {
	a := &api{
		app,
	}

	r := mux.NewRouter()
	r.Use(exampleMiddleware)

	// todo change to version prefix
	r.HandleFunc("/v1/order", a.MakeOrder).Methods(http.MethodPost)
	r.HandleFunc("/v1/orders", a.ListOfOrders).Methods(http.MethodGet)
	r.HandleFunc("/v1/order", a.ListOfOrders).Methods(http.MethodGet)
	r.HandleFunc("/v1/order", a.ChangeOrderStatus).Methods(http.MethodPut)

	return r
}
