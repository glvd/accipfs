package accipfs_test

import (
	"github.com/glvd/accipfs/config"
	"testing"
)

func TestRun(t *testing.T) {
	config.WorkDir = "D:\\workspace\\pvt"
	config.Initialize()

}
