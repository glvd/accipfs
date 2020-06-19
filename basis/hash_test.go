package basis

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	hash, _ := EncodeHash("a", "b", "c")
	t.Log(fmt.Sprintf("%s", hash))

	decodeHash, e := DecodeHash(hash)
	if e != nil {
		t.Fatal(e)
	}
	for i, h := range decodeHash {
		t.Log(i, h)
		switch {
		case i == 0 && h == "a":
		case i == 1 && h == "b":
		case i == 2 && h == "c":
		default:
			t.Fatal(i, h)
		}
	}
	t.Log("success done")
}
