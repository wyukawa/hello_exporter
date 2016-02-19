package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

const (
	namespace = "hello" // For Prometheus metrics.
)

var (
	listeningAddress = flag.String("telemetry.address", ":9200", "Address on which to expose metrics.")
	metricsEndpoint  = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
)

type Exporter struct {
	helloCounter prometheus.Counter
}

func NewExporter() *Exporter {
	return &Exporter{
		helloCounter: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_hello_total",
			Help:      "Number of hello.",
		}),
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.helloCounter.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.helloCounter.Inc()
	e.helloCounter.Collect(ch)
}

func main() {
	flag.Parse()

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	log.Printf("Starting Server: %s", *listeningAddress)
	http.Handle(*metricsEndpoint, prometheus.Handler())
	log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
