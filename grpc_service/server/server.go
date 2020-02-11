package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/google/credstore/client"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	//	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/airborne12/golang-samples/grpc_service/config"
	"github.com/airborne12/golang-samples/grpc_service/metrics"
	"github.com/airborne12/golang-samples/grpc_service/pki"
	"github.com/airborne12/golang-samples/grpc_service/service"

	//	pb "github.com/airborne12/golang-samples/grpc_service/echo-proto"
	"github.com/airborne12/golang-samples/grpc_service/log"
)

var (
	glog         = log.Get()
	gconfig      = config.Get()
	gservicePool = service.Get()
	gcertKey     = pki.Get()
	// ListenAddress is the grpc listen address
	ListenAddress    = &gconfig.ListenAddr
	serverCert       = &gconfig.ServerCert
	serverKey        = &gconfig.ServerKey
	clientCA         = &gconfig.ClientCA
	allowedCN        = &gconfig.AllowedCN
	credStoreAddress = &gconfig.CredStoreAddress
	credStoreCA      = &gconfig.CredStoreCA
	sn               = &gconfig.SN
)

//NewRestMux create rest api handler
func NewRestMux(grpcServer *grpc.Server) (*http.ServeMux, error) {

	// get context, this allows control of the connection
	ctx := context.Background()

	var dopts []grpc.DialOption
	if !gconfig.Insecure {
		// These credentials are for the upstream connection to the GRPC server
		dcreds := credentials.NewTLS(&tls.Config{
			ServerName:   *sn,
			RootCAs:      gcertKey.CertPool,
			Certificates: []tls.Certificate{*gcertKey.KeyPair},
		})
		dopts = []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	} else {
		dopts = []grpc.DialOption{grpc.WithInsecure()}
	}
	// Which multiplexer to register on.
	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard,
		&runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	for _, s := range gservicePool.Services {
		err := s.RegisterServiceHandlerFromEndpoint(ctx, gwmux, *ListenAddress, dopts)
		if err != nil {
			return nil, err
		}
	}
	mux := http.NewServeMux()

	// we can add any non-grpc endpoints here.
	//mux.HandleFunc("/foobar/", simpleHTTPHello)

	// register the gateway mux onto the root path.
	mux.Handle("/", gwmux)

	return mux, nil
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func PromHTTPServe(grpcServer *grpc.Server) {
	metrics.GRPCMetrics.InitializeMetrics(grpcServer)

	httpServer := &http.Server{Handler: promhttp.HandlerFor(metrics.Reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9092)}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			glog.Fatal("Unable to start a http server.")
		}
	}()
}

// ListenAndServe starts grpc server
func ListenAndServe(grpcServer *grpc.Server, otherHandler http.Handler) error {

	lis, err := net.Listen("tcp", *ListenAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	glog.Debugf("config:%v", gconfig)
	if !gconfig.Insecure {
		var h http.Handler
		if otherHandler == nil {
			h = grpcServer
		} else {
			h = grpcHandlerFunc(grpcServer, otherHandler)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{*gcertKey.KeyPair},
			NextProtos:   []string{"h2"},
		}

		if gconfig.AllowedCN != "" {
			tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				for _, chains := range verifiedChains {
					if len(chains) != 0 {
						if *allowedCN == chains[0].Subject.CommonName {
							return nil
						}
					}
				}
				return errors.New("CommonName authentication failed")
			}
		}
		if gcertKey.CertPool != nil {
			tlsConfig.ClientCAs = gcertKey.CertPool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			glog.Warning("no client ca provided for grpc server")
		}

		httpsServer := &http.Server{
			Handler:   h,
			TLSConfig: tlsConfig,
		}

		glog.Infof("serving on %v", *ListenAddress)
		err = httpsServer.Serve(tls.NewListener(lis, httpsServer.TLSConfig))
		return fmt.Errorf("failed to serve: %v", err)
	}
	//do not support rest here
	glog.Warningf("serving INSECURE on %v", *ListenAddress)
	err = grpcServer.Serve(lis)
	return fmt.Errorf("failed to serve: %v", err)
}

// NewServer creates a new GRPC server stub with credstore auth (if requested).
func NewServer() (*grpc.Server, *client.CredstoreClient, error) {
	var grpcServer *grpc.Server
	var cc *client.CredstoreClient

	if *credStoreAddress != "" {
		var err error
		cc, err = client.NewCredstoreClient(context.Background(), *credStoreAddress, *credStoreCA)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to init credstore: %v", err)
		}

		glog.Infof("enabled credstore auth")
		grpcServer = grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				grpc_prometheus.UnaryServerInterceptor,
				client.CredStoreTokenInterceptor(cc.SigningKey()),
				client.CredStoreMethodAuthInterceptor(),
			)))
	} else {
		grpcServer = grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				metrics.GRPCMetrics.UnaryServerInterceptor(),
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
			)))
	}
	for _, s := range gservicePool.Services {
		s.Register2PB(grpcServer)
	}
	reflection.Register(grpcServer)
	grpc_prometheus.Register(grpcServer)

	return grpcServer, cc, nil
}
