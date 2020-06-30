package basis

import (
	"encoding/hex"
	"math"
	"math/big"
)

// INT256 ...
type INT256 struct {
	bits [32]uint8
}

// String ...
func (i INT256) String() string {
	return hex.EncodeToString(i.bits[:])
}

// AsByteArray ...
func (i *INT256) AsByteArray() [32]byte {
	return i.bits
}

// ByteString ...
func (i *INT256) ByteString() string {
	return string(i.bits[:])
}

// BitLen ...
func (i *INT256) BitLen() int {
	var a big.Int
	a.SetBytes(i.bits[:])
	return a.BitLen()
}

// SetBytes ...
func (i *INT256) SetBytes(b []byte) {
	n := copy(i.bits[:], b)
	if n != 20 {
		panic(n)
	}
}

// Bytes ...
func (i INT256) Bytes() []byte {
	return i.bits[:]
}

// Cmp ...
func (i INT256) Cmp(r INT256) int {
	for n := range i.bits {
		if i.bits[n] < r.bits[n] {
			return -1
		} else if i.bits[n] > r.bits[n] {
			return 1
		}
	}
	return 0
}

// SetMax ...
func (i *INT256) SetMax() {
	for n := range i.bits {
		i.bits[n] = math.MaxUint8
	}
}

// Xor ...
func (i *INT256) Xor(a, b *INT256) {
	for n := range i.bits {
		i.bits[n] = a.bits[n] ^ b.bits[n]
	}
}

// IsZero ...
func (i *INT256) IsZero() bool {
	for _, b := range i.bits {
		if b != 0 {
			return false
		}
	}
	return true
}

// INT256FromBytes ...
func INT256FromBytes(b []byte) (ret INT256) {
	ret.SetBytes(b)
	return
}

// INT256FromByteArray ...
func INT256FromByteArray(b [32]byte) (ret INT256) {
	ret.SetBytes(b[:])
	return
}

// INT256FromByteString ...
func INT256FromByteString(s string) (ret INT256) {
	ret.SetBytes([]byte(s))
	return
}

// Distance ...
func Distance(a, b INT256) (ret INT256) {
	ret.Xor(&a, &b)
	return
}
