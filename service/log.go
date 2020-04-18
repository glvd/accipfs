package service

import (
	"fmt"
	alog "github.com/glvd/accipfs/log"
)

const module = "service"

// LogOutput ...
var LogOutput = true

func logI(msg string, v ...interface{}) {
	alog.Module(module).Infow(msg, v...)
}
func logD(msg string, v ...interface{}) {
	alog.Module(module).Debugw(msg, v...)
}
func logE(msg string, v ...interface{}) {
	alog.Module(module).Errorw(msg, v...)
}

func output(v ...interface{}) {
	fmt.Printf("[%s]:%+v\n", module, v)
}
