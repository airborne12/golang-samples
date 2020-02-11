package metrics

import (
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Reg Create a metrics registry.
	Reg = prometheus.NewRegistry()

	// GRPCMetrics Create some standard server metrics.
	GRPCMetrics = grpc_prometheus.NewServerMetrics()

	// CustomizedCounterMetric Create a customized counter metric.
	CustomizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
)

func init() {
	// Register standard server metrics and customized metrics to registry.
	Reg.MustRegister(GRPCMetrics, CustomizedCounterMetric)
	//CustomizedCounterMetric.WithLabelValues("Test")
}
