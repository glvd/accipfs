package core

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersion_String(t *testing.T) {
	v := Version{
		'v', 1, 3, 5,
	}
	t.Log(strings.Compare(v.String(), "v1.3.5"))
}

func TestParseVersion(t *testing.T) {
	version, err := ParseVersion("v1.3.5")
	if err != nil {
		t.Fatal(err)
	}
	b := [4]byte{
		'v', 1, 3, 5,
	}

	t.Log(bytes.Compare(version[:], b[:]))

	version1, err := ParseVersion("v1.3")
	if err != nil {
		t.Fatal(err)
	}
	b1 := [4]byte{
		'v', 1, 3,
	}

	t.Log(bytes.Compare(version1[:], b1[:]))

}
