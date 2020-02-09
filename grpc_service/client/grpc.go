package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewGRPCConn is a helper wrapper around grpc.Dial.
func NewGRPCConn(
	address string,
	serverCAFileName string,
	clientCertFileName string,
	clientKeyFileName string,
) (*grpc.ClientConn, error) {
	if serverCAFileName == "" {
		return grpc.Dial(address,
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	}

	caCert, err := ioutil.ReadFile(serverCAFileName)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cfg := &tls.Config{
		RootCAs: caCertPool,
	}

	if clientCertFileName != "" && clientKeyFileName != "" {
		peerCert, err := tls.LoadX509KeyPair(clientCertFileName, clientKeyFileName)
		if err != nil {
			return nil, err
		}
		cfg.Certificates = []tls.Certificate{peerCert}
	}

	return grpc.Dial(address,
		grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
}
