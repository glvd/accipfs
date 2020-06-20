package node

import "encoding/json"

// Type ...
type Type int16

// RequestTypeID ...
const (
	// RequestID ...
	RequestID Type = 0x01
	// ResponseID ...
	ResponseID Type = 0x02
)

// Status ...
type Status int

const (
	// StatusOK ...
	StatusOK = 0x00
	// StatusFailed ...
	StatusFailed = 0x01
)

// Exchange ...
type Exchange struct {
	Version Version
	Type    Type
	Session int16
	Length  int64

	Status Status
	Data   []byte
}

// JSON ...
func (e Exchange) JSON() []byte {
	marshal, _ := json.Marshal(e)
	return marshal
}
