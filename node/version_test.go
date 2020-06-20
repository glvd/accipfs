package node

import (
	"strings"
	"testing"
)

func TestVersion_String(t *testing.T) {
	v := Version{
		'v', 1, 3, 5,
	}
	t.Log(strings.Compare(v.String(), "v1.3.5"))
}
