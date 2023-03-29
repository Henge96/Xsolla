package dom

// OrderStatus for orders in app.
type OrderStatus uint8

//go:generate stringer -output=stringer.OrderStatus.go -type=OrderStatus -trimprefix=OrderStatus
const (
	_ OrderStatus = iota
	OrderStatusNew
	OrderStatusConfirmed
	OrderStatusCanceled
	OrderStatusCooking
	OrderStatusCooked
	OrderStatusDelivering
	OrderStatusCompleted
)
