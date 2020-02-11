package main

import (
	"context"
	"os"

	"github.com/airborne12/golang-samples/grpc_service/client"
	"github.com/airborne12/golang-samples/grpc_service/config"
	pb "github.com/airborne12/golang-samples/grpc_service/echo-proto"
	"github.com/airborne12/golang-samples/grpc_service/log"
)

var (
	glog          = log.Get()
	gconfig       = config.Get()
	serverAddress = gconfig.ListenAddr
	serverCert    = gconfig.ServerCert
	serverKey     = gconfig.ServerKey
	clientCA      = gconfig.ClientCA
	sn            = gconfig.SN
	defaultName   = "world"
)

func main() {
	conn, err := client.NewGRPCConn(serverAddress, clientCA, serverCert, serverKey, sn)
	if err != nil {
		glog.Fatalf("NewGRPCConn failed:%v", err)
	}
	defer conn.Close()
	c := pb.NewEchoServiceClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r1, err := c.Hello(context.Background(), &pb.EchoMessage{Body: name})
	if err != nil {
		glog.Fatalf("could not greet: %v", err)
	}
	glog.Debug(r1.Body)
	r2, err := c.Echo(context.Background(), &pb.EchoMessage{Body: name})
	if err != nil {
		glog.Fatalf("could not greet: %v", err)
	}
	glog.Debug(r2.Body)

}
