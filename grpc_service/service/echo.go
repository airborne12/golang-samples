package service

import (
	"context"

	pb "github.com/airborne12/golang-samples/grpc_service/echo-proto"
	"github.com/airborne12/golang-samples/grpc_service/log"
	"github.com/airborne12/golang-samples/grpc_service/metrics"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	glog = log.Get()
)

//EchoService is sample service for service proto
type EchoService struct{}

// Hello implements hello function of EchoServiceServer
func (s *EchoService) Hello(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	glog.Debugf("Echo: get Hello %v", in)
	return &pb.EchoMessage{Body: "Hello from server!"}, nil
}

// Echo implements echo function of EchoServiceServer
func (s *EchoService) Echo(ctx context.Context, in *pb.EchoMessage) (*pb.EchoMessage, error) {
	glog.Debugf("Echo: Received '%v'", in.Body)
	metrics.CustomizedCounterMetric.WithLabelValues(in.Body).Inc()
	return &pb.EchoMessage{Body: "ACK " + in.Body}, nil
}

// Register2PB service to pb
func (s *EchoService) Register2PB(grpcServer *grpc.Server) {
	glog.Debug("Echo: Register2PB")
	pb.RegisterEchoServiceServer(grpcServer, &EchoService{})
}

//RegisterServiceHandlerFromEndpoint register service to rest api
func (s *EchoService) RegisterServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	glog.Debug("Echo: RegisterServiceHandlerFromEndpoint")
	err := pb.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		glog.Errorf("RegisterEchoServiceHandlerFromEndpoint error: %v", err)
		return err
	}
	return nil
}
