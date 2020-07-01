package hash

import (
	"bytes"
	"testing"
)

func TestSum(t *testing.T) {
	sum, err := Sum(map[string]string{
		"vala": "a",
		"valb": "b",
		"valc": "c",
	}, &Options{
		Hash:    nil,
		TagName: "",
		ZeroNil: false,
	})
	if err != nil {
		return
	}
	t.Logf("%x\n", sum)
	sum1, err := Sum(map[string]string{
		"valb": "b",
		"vala": "a",
		"valc": "c",
	}, &Options{
		Hash:    nil,
		TagName: "",
		ZeroNil: false,
	})
	if err != nil {
		return
	}
	t.Logf("%x\n", sum1)
	if bytes.Compare(sum, sum1) != 0 {
		t.Fatalf("sum:%x,sum1:%x", sum, sum1)
	}
}

func TestSumString(t *testing.T) {
	sum, err := Sum("hello world", &Options{
		Hash:    nil,
		TagName: "",
		ZeroNil: false,
	})
	if err != nil {
		return
	}
	t.Logf("%x\n", sum)
	sum1, err := Sum("hello world", &Options{
		Hash:    nil,
		TagName: "",
		ZeroNil: false,
	})
	if err != nil {
		return
	}
	t.Logf("%x\n", sum1)
	if bytes.Compare(sum, sum1) != 0 {
		t.Fatalf("sum:%x,sum1:%x", sum, sum1)
	}
}
