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

	responseListOrders struct {
		ID string `json:"id"`
	}

	requestUpdateOrderStatus struct {
		ID     string          `json:"id"`
		Status dom.OrderStatus `json:"status"`
	}

	item struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Count uint16 `json:"count"`
	}

	address struct {
		City     string `json:"city"`
		Street   string `json:"street"`
		House    string `json:"house"`
		Entrance string `json:"entrance"`
		Flat     string `json:"flat"`
	}
)

func (i item) convert() app.Item {
	return app.Item{}
}
