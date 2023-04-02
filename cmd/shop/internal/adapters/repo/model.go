package repo

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"strings"
	"time"
	"xsolla/cmd/shop/internal/app"
	"xsolla/internal/dom"
)

// todo structs fields and converts
type (
	order struct {
		ID        uuid.UUID       `db:"id" json:"id"`
		Address   []byte          `db:"address" json:"address"`
		Items     []byte          `db:"items" json:"items"`
		Status    dom.OrderStatus `db:"status" json:"status"`
		Comment   string          `db:"comment" json:"comment"`
		CreatedAt time.Time       `db:"created_at" json:"created_at"`
		UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
	}

	task struct {
		ID         uuid.UUID `id:"id"`
		OrderBytes []byte    `db:"order_bytes"`
		Kind       string    `db:"kind"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
		FinishedAt time.Time `db:"finished_at"`
	}

	address struct {
		City     string `db:"city" json:"city"`
		Street   string `db:"street" json:"street"`
		House    string `db:"house" json:"house"`
		Entrance string `db:"entrance" json:"entrance"`
		Flat     string `db:"flat" json:"flat"`
	}

	item struct {
		OrderID uuid.UUID `db:"order_id"`
		Product []byte    `db:"product"`
		Count   uint16    `db:"count"`
		Comment string    `db:"comment"`
	}

	product struct {
		ID   uuid.UUID `db:"id" json:"id"`
		Type string    `db:"type" json:"type"`
		Name string    `db:"name" json:"name"`
	}
)

func convertToOrder(o app.Order) *order {
	return &order{}
}

func convertToAddress(a app.Address) *address {
	return &address{}
}

func convertToProduct(p app.Product) *product {
	return &product{}
}



func prepareItems(items []app.Item) (string, error) {
	valueStrings := make([]string, 0, len(items))
	valueArgs := make([]interface{}, 0, len(items)*3)
	for _, i := range items {
		repoProduct := convertToProduct(i.Product)
		repoProductBytes, err := json.Marshal(repoProduct)
		if err != nil {
			return "", fmt.Errorf("json.Marshal: %w", err)
		}

		valueStrings = append(valueStrings, "($1, $2, $3, $4)")
		valueArgs = append(valueArgs, i.OrderID)
		valueArgs = append(valueArgs, repoProductBytes)
		valueArgs = append(valueArgs, i.Count)
		valueArgs = append(valueArgs, i.Comment)
	}
	return fmt.Sprintf("INSERT INTO items (order_id, product, count, comment) VALUES %s",
		strings.Join(valueStrings, ",")), nil
}
