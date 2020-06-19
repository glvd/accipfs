package node

import "encoding/json"

// Type ...
type Type int

// RequestTypeID ...
const (
	// RequestID ...
	RequestID Type = 0x01
	// ResponseID ...
	ResponseID Type = 0x02
)

// Exchange ...
type Exchange struct {
	Type Type
	Data []byte
}

// JSON ...
func (e Exchange) JSON() []byte {
	marshal, _ := json.Marshal(e)
	return marshal
}
