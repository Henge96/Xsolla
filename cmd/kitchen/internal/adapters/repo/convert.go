package repo

import (
	"encoding/json"
	"fmt"
	"xsolla/cmd/kitchen/internal/app"
	"xsolla/internal/dom"
)

func convertToOrder(o app.Order) (*order, error) {
	return &order{}, nil
}

func convertTask(t app.Task) (*task, error) {
	repoOrder, err := convertToOrder(t.Order)
	if err != nil {
		return nil, fmt.Errorf("convertToOrder: %w", err)
	}

	orderBytes, err := json.Marshal(repoOrder)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return &task{
		ID:         t.ID,
		OrderBytes: orderBytes,
		Kind:       t.TaskKind.String(),
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
		FinishedAt: t.FinishedAt,
	}, nil
}

func (t *task) convert() (*app.Task, error) {
	var o order
	err := json.Unmarshal(t.OrderBytes, &o)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	appOrder, err := o.convert()
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	taskKind := dom.DomTaskKind(o.Status)
	if taskKind == 0 {
		panic(fmt.Sprintf("unknown status in task: %v", o))
	}

	return &app.Task{
		ID:         t.ID,
		Order:      *appOrder,
		TaskKind:   taskKind,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
		FinishedAt: t.FinishedAt,
	}, nil
}

func (o *order) convert() (*app.Order, error) {
	items := make([]item, 0)
	err := json.Unmarshal(o.Items.Bytes, &items)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", convertError(err))
	}

	status := dom.ToOrderStatus(o.Status)
	if status == 0 {
		panic(fmt.Sprintf("unknown status in order: %v", o))
	}

	appItems := make([]app.Item, 0, len(items))
	for i := range items {
		appItem, err := items[i].convert()
		if err != nil {
			return nil, fmt.Errorf("items.convert: %w", err)
		}

		appItems = append(appItems, *appItem)
	}

	return &app.Order{
		ID:        o.ID,
		Items:     appItems,
		Status:    status,
		Comment:   o.Comment,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}, nil
}

func (i *item) convert() (*app.Item, error) {
	var p product
	err := json.Unmarshal(i.Product.Bytes, &p)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", convertError(err))
	}

	return &app.Item{
		ID:      i.ID,
		OrderID: i.OrderID,
		Product: *p.convert(),
		Count:   i.Count,
		Comment: i.Comment,
	}, nil
}

func (p *product) convert() *app.Product {
	return &app.Product{
		ID:   p.ID,
		Type: p.Type,
		Name: p.Name,
	}
}
