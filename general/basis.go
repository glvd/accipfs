package general

import (
	"os"
)

// CurrentDir ...
func CurrentDir() string {
	dir, err := os.Getwd()
	if err == nil {
		return dir
	}
	return "."
}
