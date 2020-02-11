package main

import (
	"github.com/airborne12/golang-samples/grpc_service/server"

	"github.com/airborne12/golang-samples/grpc_service/config"
	"github.com/airborne12/golang-samples/grpc_service/log"
)

var glog = log.Get()
var gconfig = config.Get()

func main() {
	grpcServer, _, err := server.NewServer()
	if err != nil {
		glog.Fatalf("GRPC server create failed:", err)
	}
	restHandler, err := server.NewRestMux(grpcServer)
	if err != nil {
		glog.Fatalf("New rest http handler failed:", err)
	}
	server.PromHTTPServe(grpcServer)
	server.ListenAndServe(grpcServer, restHandler)
}
