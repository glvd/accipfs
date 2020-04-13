package controller

import (
	"fmt"
	alog "github.com/glvd/accipfs/log"
)

const module = "controller"

func logI(msg string, v ...interface{}) {
	alog.Module(module).Infow(msg, v...)
}

func logE(msg string, v ...interface{}) {
	alog.Module(module).Errorw(msg, v...)
}

func output(v ...interface{}) {
	fmt.Printf("[%s]:%+v\n", module, v)
}
