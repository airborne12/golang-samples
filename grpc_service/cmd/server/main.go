package main

import (
	"github.com/airborne12/golang-samples/grpc_service/server"

	"github.com/airborne12/golang-samples/configmodule/config"
	"github.com/airborne12/golang-samples/logmodule/log"
)

var glog = log.Get()
var gconfig = config.Get()

func main() {
	grpcServer, _, err := server.NewServer()
	if err != nil {
		glog.Fatalf("GRPC server create failed:", err)
	}
	server.ListenAndServe(grpcServer, nil)
}
