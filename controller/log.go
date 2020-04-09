package controller

import alog "github.com/glvd/accipfs/log"

const module = "controller"

func logI(msg string, v ...interface{}) {
	alog.Module(module).Infow(msg, v...)
}
