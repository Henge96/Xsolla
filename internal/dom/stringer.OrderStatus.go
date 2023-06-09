// Code generated by "stringer -output=stringer.OrderStatus.go -type=OrderStatus -trimprefix=OrderStatus"; DO NOT EDIT.

package dom

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OrderStatusNew-1]
	_ = x[OrderStatusConfirmed-2]
	_ = x[OrderStatusCanceled-3]
	_ = x[OrderStatusCooking-4]
	_ = x[OrderStatusCooked-5]
	_ = x[OrderStatusDelivering-6]
	_ = x[OrderStatusCompleted-7]
}

const _OrderStatus_name = "NewConfirmedCanceledCookingCookedDeliveringCompleted"

var _OrderStatus_index = [...]uint8{0, 3, 12, 20, 27, 33, 43, 52}

func (i OrderStatus) String() string {
	i -= 1
	if i >= OrderStatus(len(_OrderStatus_index)-1) {
		return "OrderStatus(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _OrderStatus_name[_OrderStatus_index[i]:_OrderStatus_index[i+1]]
}
