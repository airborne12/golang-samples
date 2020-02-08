package main


import (
	"github.com/airborne12/golang-samples/logmodule/log"
)

func main(){
	logger := log.Get()
	logger.Info("x")
}
