package accipfs

import (
	"fmt"
	"os/exec"
)

// DefaultPath ...
var DefaultPath = "."

// CheckPort ...
func CheckPort(port int) error {
	checkStatement := fmt.Sprintf("netstat -anp | grep -q %d ", port)
	output, err := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	if err != nil {
		return nil
	}

	if len(output) > 0 {
		return fmt.Errorf("port %d already occupied", port)
	}

	return nil
}
