package http

import (
	"xsolla/cmd/shop/internal/app"
	"xsolla/internal/dom"
)

type (
	responseMakeOrder struct {
		ID string `json:"id"`
	}

	requestMakeOrder struct {
		Items   []item  `json:"items"`
		Address address `json:"address"`
		Comment string  `json:"comment"`
	}

	requestListOrders struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	// todo add fields to response
	responseListOrders struct {
		ID string `json:"id"`
	}

	requestUpdateOrderStatus struct {
		ID     string          `json:"id"`
		Status dom.OrderStatus `json:"status"`
	}

	item struct {
		Product product `json:"product"`
		Count   uint16  `json:"count"`
		Comment string  `json:"comment"`
	}

	address struct {
		City     string `json:"city"`
		Street   string `json:"street"`
		House    string `json:"house"`
		Entrance string `json:"entrance"`
		Flat     string `json:"flat"`
	}

	product struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}
)

// todo convert logic
func (p product) convert() app.Product {
	return app.Product{}
}

// todo convert logic
func (i item) convert() app.Item {
	return app.Item{
		Product: i.Product.convert(),
	}
}
