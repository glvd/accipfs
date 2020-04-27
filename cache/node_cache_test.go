package cache

import (
	"fmt"
	"github.com/glvd/accipfs/core"
	"testing"
	"time"
)

func TestFault(t *testing.T) {
	n := &core.Node{LastTime: time.Now()}
	for {
		remain, fa := faultTimeCheck(n, 5)
		fmt.Println("is fault", fa, "remain", remain)
		time.Sleep(1 * time.Second)
	}
}
