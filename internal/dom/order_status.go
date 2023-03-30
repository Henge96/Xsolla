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

func ToOrderStatus(in string) OrderStatus {
	switch in {
	case OrderStatusCompleted.String():
		return OrderStatusCompleted
	case OrderStatusDelivering.String():
		return OrderStatusDelivering
	case OrderStatusCooked.String():
		return OrderStatusCooked
	case OrderStatusCooking.String():
		return OrderStatusCooking
	case OrderStatusCanceled.String():
		return OrderStatusCanceled
	case OrderStatusConfirmed.String():
		return OrderStatusConfirmed
	case OrderStatusNew.String():
		return OrderStatusNew
	default:
		return 0
	}
}
