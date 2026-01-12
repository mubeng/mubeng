package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "mubeng"

	LabelMethod    = "method"
	LabelStatus    = "status"
	LabelProxy     = "proxy"
	LabelErrorType = "error_type"
	LabelOutcome   = "outcome"
	LabelRetried   = "retried"

	OutcomeSuccess = "success"
	OutcomeFailure = "failure"

	ErrorTypeTimeout           = "timeout"
	ErrorTypeConnectionRefused = "connection_refused"
	ErrorTypeConnectionReset   = "connection_reset"
	ErrorTypeDNS               = "dns_error"
	ErrorTypeTLS               = "tls_error"
	ErrorTypeProxyAuth         = "proxy_auth_failed"
	ErrorTypeOther             = "other"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "requests_total",
			Help:      "Total number of requests processed (final outcome)",
		},
		[]string{LabelMethod, LabelStatus, LabelProxy, LabelRetried},
	)

	ProxyAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "proxy_attempts_total",
			Help:      "Total proxy attempts including retries (each attempt counted)",
		},
		[]string{LabelProxy, LabelOutcome},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "request_duration_seconds",
			Help:      "Successful request latency distribution in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{LabelMethod, LabelProxy},
	)

	RequestErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "request_errors_total",
			Help:      "Total number of request errors by type",
		},
		[]string{LabelErrorType, LabelProxy},
	)

	RetriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "retries_total",
			Help:      "Total number of retry attempts",
		},
		[]string{LabelProxy},
	)

	ProxyPoolSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "proxy_pool_size",
			Help:      "Current number of proxies in the pool",
		},
	)

	ProxyRemovalsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "proxy_removals_total",
			Help:      "Total number of proxies removed due to failures",
		},
		[]string{LabelProxy},
	)

	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "active_connections",
			Help:      "Current number of active connections",
		},
	)

	ProxyRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "proxy_requests_total",
			Help:      "Total requests per proxy with success/fail breakdown",
		},
		[]string{LabelProxy, LabelStatus},
	)
)
