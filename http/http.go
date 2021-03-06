package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	// HTTPRequestCount is a Prometheus Counter that counts the total number of http requests made
	httpRequestCount *prometheus.CounterVec

	// HTTPErrorCount is a Prometheus Counter that counts the total number of errors
	httpErrorCount *prometheus.CounterVec

	// HTTPEgressBytes is a Prometheus Counter that counts the number of bytes sent by responses
	httpEgressBytes *prometheus.CounterVec

	// HTTPIngressBytes is a Prometheus Counter that counts the number of bytes received during requests
	httpIngressBytes *prometheus.CounterVec

	// HTTPDurationSeconds is a Prometheus Historgram that measures the duration of http requests
	httpDurationSeconds *prometheus.HistogramVec
)

func init() {
	httpRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nvgo_http_request_count",
			Help: "The total number of http requests",
		},
		[]string{"service", "tenant", "method", "code"},
	)

	httpErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nvgo_http_request_error_count",
			Help: "The total number of http errors",
		},
		[]string{"service", "tenant", "method", "code"},
	)

	httpEgressBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nvgo_http_egress_bytes",
			Help: "The total number of bytes sent back to the requesting client",
		},
		[]string{"service", "tenant", "method"},
	)

	httpIngressBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nvgo_http_ingress_bytes",
			Help: "The total number of bytes sent by the requesting client",
		},
		[]string{"service", "tenant", "method"},
	)

	httpDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "nvgo_http_request_duration_seconds",
			Help: "The time taken to process a request",
		},
		[]string{"service", "tenant", "method"},
	)

	prometheus.MustRegister(httpRequestCount)
	prometheus.MustRegister(httpErrorCount)
	prometheus.MustRegister(httpEgressBytes)
	prometheus.MustRegister(httpIngressBytes)
	prometheus.MustRegister(httpDurationSeconds)
}

// Client defines the API for an http metrics client
type Client interface {
	TrackRequest(handler http.Handler, w http.ResponseWriter, r *http.Request, opts *TrackingOpts)
}

// NewClient builds and returns a new client
func NewClient() Client {
	return NewClientWithDefaults(&TrackingOpts{
		Service: "none",
		Tenant:  "system",
	})
}

// NewClientWithDefaults builds and returns a new client using the provided default settings
func NewClientWithDefaults(opts *TrackingOpts) Client {
	logrus.Debug("Creating new prometheus client")
	return &client{opts}
}

// TrackingOpts define options for tracking http requests
type TrackingOpts struct {
	Service string
	Tenant  string
}

type client struct {
	defaultOpts *TrackingOpts
}

func (c *client) TrackRequest(handler http.Handler, w http.ResponseWriter, r *http.Request, opts *TrackingOpts) {
	logrus.Debug("Incrementing Prometheus Counters")

	service := opts.Service
	if service == "" {
		service = c.defaultOpts.Service
	}

	tenant := opts.Tenant
	if tenant == "" {
		tenant = c.defaultOpts.Tenant
	}

	// make sure service and tenant are lower case
	service = strings.ToLower(service)
	tenant = strings.ToLower(tenant)

	m := httpsnoop.CaptureMetrics(handler, w, r)

	respCode := "000"
	if m.Code >= 0 {
		respCode = fmt.Sprintf("%d", m.Code)

		if m.Code > http.StatusBadRequest {
			httpErrorCount.WithLabelValues(service, tenant, r.Method, respCode).Inc()
		}
	}

	// Increment request count
	httpRequestCount.WithLabelValues(service, tenant, r.Method, respCode).Inc()
	// Increment Ingress
	httpIngressBytes.WithLabelValues(service, tenant, r.Method).Add(float64(r.ContentLength))
	// Increment Egress
	httpEgressBytes.WithLabelValues(service, tenant, r.Method).Add(float64(m.Written))
	// Observe duration
	httpDurationSeconds.WithLabelValues(service, tenant, r.Method).Observe(m.Duration.Seconds())
}
