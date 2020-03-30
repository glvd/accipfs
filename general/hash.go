package general

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"
)

// EncodeHash ...
func EncodeHash(ss ...string) (hash string, err error) {
	buffer := bytes.NewBuffer(nil)
	encoder := hex.NewEncoder(buffer)
	i, e := encoder.Write([]byte(strings.Join(ss, "_")))
	if e != nil || i == 0 {
		return "", errors.New("encode failed")
	}
	return buffer.String(), nil
}

// DecodeHash ...
func DecodeHash(hash string) (ss []string, err error) {
	read := strings.NewReader(hash)
	decoder := hex.NewDecoder(read)
	dst := make([]byte, hex.DecodedLen(len(hash)))
	i, e := decoder.Read(dst)
	if e != nil || i == 0 {
		return nil, errors.New("decode failed")
	}
	ss = strings.Split(string(dst), "_")
	return ss, nil
}
