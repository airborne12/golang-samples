package service

import (
	"context"
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

)

var gPool *Pool

func init() {
	gPool = NewPool()
}

//Service called by grpc or rest api
type Service interface {
	Register2PB(grpcServer *grpc.Server)
	RegisterServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

//Pool is container of all services
type Pool struct {
	Services []Service
}

//NewPool register all services
func NewPool() *Pool {
	services := []Service{
		&EchoService{},
	}
	return &Pool{
		Services: services,
	}
}

//Get returns a Pool object
func Get() *Pool {
	return gPool
}
