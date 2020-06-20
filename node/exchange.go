package node

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
)

// Type ...
type Type uint8

const (
	// ErrorType ...
	ErrorType Type = 0x00
	// RequestID ...
	Request Type = 0x01
	// ResponseID ...
	Response Type = 0x02
)

// TypeDetail ...
type TypeDetail uint16

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

// Unpack ...
func (e Exchange) Unpack(reader io.Reader) (err error) {
	var v []interface{}
	v = append(v, &e.Version, &e.Type, &e.Session, &e.Length)
	for i := range v {
		err = binary.Read(reader, binary.BigEndian, v[i])
		if err != nil {
			return err
		}
	}
	e.Data = make([]byte, e.Length)
	return binary.Read(reader, binary.BigEndian, &e.Data)
}

// ScanExchange ...
func ScanExchange(conn net.Conn) (*Exchange, error) {
	var ex Exchange
	scan := dataScan(conn)
	r := bytes.NewReader(scan.Bytes())
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
			if len(data) > 16 {
				length := int64(0)
				err := binary.Read(bytes.NewReader(data[8:8]), binary.BigEndian, &length)
				if err != nil {
					return 0, nil, err
				}
				length += 16
				if int(length) <= len(data) {
					return int(length), data[:int(length)], nil
				}
			}
		}
		return
	})
	return scanner
}
