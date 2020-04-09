package service

import alog "github.com/glvd/accipfs/log"

const module = "service"

func log(msg string, v ...interface{}) {
	alog.Module(module).Infow(msg, v...)
}
