package node

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

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

// Pack ...
func (e Exchange) Pack(writer io.Writer) (err error) {
	var v []interface{}
	v = append(v, &e.Version, &e.Type, &e.Session, &e.Length, &e.Data)
	for i := range v {
		err = binary.Write(writer, binary.BigEndian, v[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
