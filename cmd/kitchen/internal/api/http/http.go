package http

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"net/http"
	"xsolla/cmd/kitchen/internal/app"
)

type (
	application interface {
		ListCooking(ctx context.Context, params app.CookingParams) ([]app.Cooking, int, error)
		ChangeCookingStatus(ctx context.Context, cookingID uuid.UUID, status app.CookingStatus) error
	}

	requestValidator interface {
		Struct(s interface{}) error
	}
)

type api struct {
	app       application
	validator requestValidator
}

func New(app application) http.Handler {
	v := validator.New()

	a := &api{
		app,
		v,
	}

	r := mux.NewRouter()
	// example
	r.Use()

	// todo change to version prefix
	r.HandleFunc("/v1/cookings", a.ListOfCooking).Methods(http.MethodGet)
	r.HandleFunc("/v1/cooking", a.ChangeCookingStatus).Methods(http.MethodPut)

	return r
}
