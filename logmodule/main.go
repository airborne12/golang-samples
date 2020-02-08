package main

import (
	"github.com/airborne12/golang-samples/logmodule/log"
)

var logger = log.Get()

func bcd() {
	logger.Debug("x")
}
func main() {

	logger.Critical("jkjk")
	logger.Error("z")

	logger.Warning("z")

	logger.Info("y")
	bcd()
}
