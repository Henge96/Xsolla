// Code generated by "stringer -output=stringer.CookingStatus.go -type=CookingStatus -trimprefix=CookingStatus"; DO NOT EDIT.

package app

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OrderStatusNew-1]
	_ = x[OrderStatusCompleted-2]
}

const _CookingStatus_name = "OrderStatusNewOrderStatusCompleted"

var _CookingStatus_index = [...]uint8{0, 14, 34}

func (i CookingStatus) String() string {
	i -= 1
	if i >= CookingStatus(len(_CookingStatus_index)-1) {
		return "CookingStatus(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _CookingStatus_name[_CookingStatus_index[i]:_CookingStatus_index[i+1]]
}
