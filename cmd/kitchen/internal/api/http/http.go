package http

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"net/http"
	"xsolla/cmd/kitchen/internal/app"
)

type application interface {
	ListCooking(ctx context.Context, params app.CookingParams) ([]app.Cooking, int, error)
	ChangeCookingStatus(ctx context.Context, cookingID uuid.UUID, status app.CookingStatus) error
}

type api struct {
	app application
}

func New(app application) http.Handler {
	a := &api{
		app,
	}

	r := mux.NewRouter()
	// example
	r.Use()

	// todo change to version prefix
	r.HandleFunc("/v1/cookings", a.ListOfCooking).Methods(http.MethodGet)
	r.HandleFunc("/v1/cooking", a.ChangeCookingStatus).Methods(http.MethodPut)

	return r
}
