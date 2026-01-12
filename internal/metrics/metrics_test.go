package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRequestsTotal(t *testing.T) {
	RequestsTotal.WithLabelValues("GET", "200", "http://proxy1:8080").Inc()
	RequestsTotal.WithLabelValues("GET", "200", "http://proxy1:8080").Inc()
	RequestsTotal.WithLabelValues("POST", "500", "http://proxy2:8080").Inc()

	count := testutil.ToFloat64(RequestsTotal.WithLabelValues("GET", "200", "http://proxy1:8080"))
	if count != 2 {
		t.Errorf("expected 2 GET 200 requests, got %f", count)
	}

	count = testutil.ToFloat64(RequestsTotal.WithLabelValues("POST", "500", "http://proxy2:8080"))
	if count != 1 {
		t.Errorf("expected 1 POST 500 request, got %f", count)
	}
}

func TestRequestDuration(t *testing.T) {
	RequestDuration.WithLabelValues("GET", "http://proxy1:8080").Observe(0.5)
	RequestDuration.WithLabelValues("GET", "http://proxy1:8080").Observe(1.5)

	count, err := testutil.GatherAndCount(prometheus.DefaultGatherer, "mubeng_request_duration_seconds")
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Error("expected request_duration_seconds metric to be present")
	}
}

func TestRequestErrorsTotal(t *testing.T) {
	RequestErrorsTotal.WithLabelValues(ErrorTypeTimeout, "http://proxy1:8080").Inc()
	RequestErrorsTotal.WithLabelValues(ErrorTypeConnectionRefused, "http://proxy2:8080").Inc()

	count := testutil.ToFloat64(RequestErrorsTotal.WithLabelValues(ErrorTypeTimeout, "http://proxy1:8080"))
	if count != 1 {
		t.Errorf("expected 1 timeout error, got %f", count)
	}
}

func TestRetriesTotal(t *testing.T) {
	RetriesTotal.WithLabelValues("http://proxy1:8080").Inc()
	RetriesTotal.WithLabelValues("http://proxy1:8080").Inc()
	RetriesTotal.WithLabelValues("http://proxy1:8080").Inc()

	count := testutil.ToFloat64(RetriesTotal.WithLabelValues("http://proxy1:8080"))
	if count != 3 {
		t.Errorf("expected 3 retries, got %f", count)
	}
}

func TestProxyPoolSize(t *testing.T) {
	ProxyPoolSize.Set(10)
	count := testutil.ToFloat64(ProxyPoolSize)
	if count != 10 {
		t.Errorf("expected pool size 10, got %f", count)
	}

	ProxyPoolSize.Set(5)
	count = testutil.ToFloat64(ProxyPoolSize)
	if count != 5 {
		t.Errorf("expected pool size 5, got %f", count)
	}
}

func TestActiveConnections(t *testing.T) {
	ActiveConnections.Inc()
	ActiveConnections.Inc()

	count := testutil.ToFloat64(ActiveConnections)
	if count != 2 {
		t.Errorf("expected 2 active connections, got %f", count)
	}

	ActiveConnections.Dec()
	count = testutil.ToFloat64(ActiveConnections)
	if count != 1 {
		t.Errorf("expected 1 active connection, got %f", count)
	}
}

func TestProxyRemovalsTotal(t *testing.T) {
	ProxyRemovalsTotal.WithLabelValues("http://proxy1:8080").Inc()

	count := testutil.ToFloat64(ProxyRemovalsTotal.WithLabelValues("http://proxy1:8080"))
	if count != 1 {
		t.Errorf("expected 1 proxy removal, got %f", count)
	}
}

func TestProxyRequestsTotal(t *testing.T) {
	ProxyRequestsTotal.WithLabelValues("http://proxy1:8080", "success").Inc()
	ProxyRequestsTotal.WithLabelValues("http://proxy1:8080", "success").Inc()
	ProxyRequestsTotal.WithLabelValues("http://proxy1:8080", "error").Inc()

	successCount := testutil.ToFloat64(ProxyRequestsTotal.WithLabelValues("http://proxy1:8080", "success"))
	if successCount != 2 {
		t.Errorf("expected 2 successful proxy requests, got %f", successCount)
	}

	errorCount := testutil.ToFloat64(ProxyRequestsTotal.WithLabelValues("http://proxy1:8080", "error"))
	if errorCount != 1 {
		t.Errorf("expected 1 error proxy request, got %f", errorCount)
	}
}

func TestMetricsServerEndpoint(t *testing.T) {
	server := NewServer(":0")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.httpServer.Handler.ServeHTTP(w, r)
	})

	req := httptest.NewRequest("GET", "/metrics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "mubeng_") {
		t.Error("expected metrics output to contain mubeng_ prefix")
	}
}
