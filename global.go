package accipfs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DefaultPath ...
var DefaultPath = "."

var _env []string

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

// Environ ...
func Environ() []string {
	if _env == nil {
		return os.Environ()
	}
	return _env
}

// RegisterPathEnv ...
func RegisterPathEnv(paths ...string) {
	path := ""
	if len(paths) > 1 {
		path = strings.Join(paths, string(os.PathListSeparator))
	} else if len(paths) == 1 {
		path = paths[0]
	} else {
		return
	}

	if err := os.Setenv("PATH", strings.Join([]string{os.Getenv("PATH"), path}, string(os.PathListSeparator))); err != nil {
		panic(err)
	}

	_env = os.Environ()
}
