// Code generated by "stringer -output=stringer.AcknowledgeKind.go -type=AcknowledgeKind -trimprefix=AcknowledgeKind"; DO NOT EDIT.

package dom

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AcknowledgeKindAck-1]
	_ = x[AcknowledgeKindNack-2]
}

const _AcknowledgeKind_name = "AckNack"

var _AcknowledgeKind_index = [...]uint8{0, 3, 7}

func (i AcknowledgeKind) String() string {
	i -= 1
	if i >= AcknowledgeKind(len(_AcknowledgeKind_index)-1) {
		return "AcknowledgeKind(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _AcknowledgeKind_name[_AcknowledgeKind_index[i]:_AcknowledgeKind_index[i+1]]
}
