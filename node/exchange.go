package node

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/glvd/accipfs/core"
	"io"
	"net"
	"time"
)

// Type ...
type Type uint16

// TypeDetail ...
type TypeDetail uint16

// Status ...
type Status int32

// Exchange ...
type Exchange struct {
	Version    core.Version
	Length     uint64
	Session    uint32
	Type       Type
	TypeDetail TypeDetail
	Status     Status
	Data       []byte
}

const (
	// TypeError ...
	TypeError Type = 0x00
	// TypeRequest ...
	TypeRequest Type = 0x01
	// TypeResponse ...
	TypeResponse Type = 0x02
)
const (
	// TypeDetailPing ...
	TypeDetailPing TypeDetail = 0x00
	// TypeDetailID ...
	TypeDetailID TypeDetail = 0x01
	// TypeDetailNodes ...
	TypeDetailNodes TypeDetail = 0x02
)
const (
	// StatusOK ...
	StatusOK = 0x00
	// StatusFailed ...
	StatusFailed = 0x01
)

// NewRequestExchange ...
func NewRequestExchange(detail TypeDetail) *Exchange {
	return newExchange(TypeRequest, detail)
}

// NewResponseExchange ...
func NewResponseExchange(detail TypeDetail) *Exchange {
	return newExchange(TypeResponse, detail)
}

// newExchange ...
func newExchange(t Type, detail TypeDetail) *Exchange {
	return &Exchange{
		Version:    core.Version{'v', 0, 0, 1},
		Length:     0,
		Session:    0,
		Type:       t,
		TypeDetail: detail,
		Status:     0,
		Data:       nil,
	}
}

// SetData ...
func (e *Exchange) SetData(data []byte) {
	e.Data = data
	e.Length = uint64(len(data))
}

// JSON ...
func (e Exchange) JSON() []byte {
	marshal, _ := json.Marshal(e)
	return marshal
}

// Pack ...
func (e Exchange) Pack(writer io.Writer) (err error) {
	var v []interface{}
	v = append(v, &e.Version, &e.Length, &e.Session, &e.Type, &e.TypeDetail, &e.Status, &e.Data)
	for i := range v {
		err = binary.Write(writer, binary.BigEndian, v[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Unpack ...
func (e *Exchange) Unpack(reader io.Reader) (err error) {
	var v []interface{}
	v = append(v, &e.Version, &e.Length, &e.Session, &e.Type, &e.TypeDetail, &e.Status)
	for i := range v {
		err = binary.Read(reader, binary.BigEndian, v[i])
		if err != nil {
			return err
		}
	}
	if e.Length != 0 {
		e.Data = make([]byte, e.Length)
		return binary.Read(reader, binary.BigEndian, &e.Data)
	}
	return nil
}

// ScanExchange ...
func ScanExchange(scanner *bufio.Scanner) (*Exchange, error) {
	var ex Exchange
	r := bytes.NewReader(scanner.Bytes())
	err := ex.Unpack(r)
	if err != nil {
		return nil, err
	}
	return &ex, nil
}

func dataScan(conn net.Conn) *bufio.Scanner {
	scanner := bufio.NewScanner(conn)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && data[0] == 'v' {
			if len(data) > 12 {
				length := uint64(0)
				err := binary.Read(bytes.NewReader(data[4:12]), binary.BigEndian, &length)
				if err != nil {
					return 0, nil, err
				}
				length += 24
				if int(length) <= len(data) {
					return int(length), data[:int(length)], nil
				}
			}
		}
		return
	})
	return scanner
}

func defaultQueueOption() *QueueOption {
	return &QueueOption{
		Callback: false,
		Timeout:  30 * time.Second,
	}
}
