package node

import (
	"github.com/glvd/accipfs/config"
	"testing"
)

var testConfig = config.Default()

func init() {
	testConfig.Path = ""
	//testConfig.Node.
}

func TestNew(t *testing.T) {

	m := New(testConfig)
	//m.Push()
}
