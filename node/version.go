package node

import "fmt"

// Version ...
type Version [4]byte

// String ...
func (v Version) String() string {
	return fmt.Sprintf("%s%d.%d.%d", string(v[0]), v[1], v[2], v[3])
}
