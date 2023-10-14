package fscli

import (
	"log"
	"os"
)

//lint:ignore U1000 This is used in debug
func debugLogger() *log.Logger {
	f, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return log.New(f, "", log.LstdFlags)
}
