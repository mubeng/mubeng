package metrics

import (
	"crypto/tls"
	"errors"
	"net"
	"os"
	"syscall"
	"testing"
)

func TestClassifyError_Nil(t *testing.T) {
	result := ClassifyError(nil)
	if result != "" {
		t.Errorf("expected empty string for nil error, got %s", result)
	}
}

func TestClassifyError_Timeout(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"timeout string", errors.New("connection timeout")},
		{"context deadline exceeded", errors.New("context deadline exceeded")},
		{"i/o timeout", errors.New("i/o timeout")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyError(tt.err)
			if result != ErrorTypeTimeout {
				t.Errorf("expected %s, got %s", ErrorTypeTimeout, result)
			}
		})
	}
}

func TestClassifyError_ConnectionRefused(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"connection refused string", errors.New("connection refused")},
		{"dial tcp connection refused", errors.New("dial tcp 127.0.0.1:8080: connection refused")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyError(tt.err)
			if result != ErrorTypeConnectionRefused {
				t.Errorf("expected %s, got %s", ErrorTypeConnectionRefused, result)
			}
		})
	}
}

func TestClassifyError_ConnectionReset(t *testing.T) {
	err := errors.New("connection reset by peer")
	result := ClassifyError(err)
	if result != ErrorTypeConnectionReset {
		t.Errorf("expected %s, got %s", ErrorTypeConnectionReset, result)
	}
}

func TestClassifyError_DNS(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"no such host", errors.New("no such host")},
		{"dns error string", errors.New("dns lookup failed")},
		{"net.DNSError", &net.DNSError{Err: "no such host", Name: "example.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyError(tt.err)
			if result != ErrorTypeDNS {
				t.Errorf("expected %s, got %s", ErrorTypeDNS, result)
			}
		})
	}
}

func TestClassifyError_TLS(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"tls string", errors.New("tls: handshake failure")},
		{"certificate error", errors.New("x509: certificate signed by unknown authority")},
		{"certificate verify", errors.New("certificate verification failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyError(tt.err)
			if result != ErrorTypeTLS {
				t.Errorf("expected %s, got %s", ErrorTypeTLS, result)
			}
		})
	}
}

func TestClassifyError_TLSCertificateVerificationError(t *testing.T) {
	err := &tls.CertificateVerificationError{Err: errors.New("test")}
	result := ClassifyError(err)
	if result != ErrorTypeTLS {
		t.Errorf("expected %s, got %s", ErrorTypeTLS, result)
	}
}

func TestClassifyError_ProxyAuth(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"proxy authentication", errors.New("proxy authentication required")},
		{"407 error", errors.New("407 proxy authentication required")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyError(tt.err)
			if result != ErrorTypeProxyAuth {
				t.Errorf("expected %s, got %s", ErrorTypeProxyAuth, result)
			}
		})
	}
}

func TestClassifyError_SyscallErrors(t *testing.T) {
	t.Run("ECONNREFUSED", func(t *testing.T) {
		err := &net.OpError{
			Op:  "dial",
			Net: "tcp",
			Err: &os.SyscallError{
				Syscall: "connect",
				Err:     syscall.ECONNREFUSED,
			},
		}
		result := ClassifyError(err)
		if result != ErrorTypeConnectionRefused {
			t.Errorf("expected %s, got %s", ErrorTypeConnectionRefused, result)
		}
	})

	t.Run("ECONNRESET", func(t *testing.T) {
		err := &net.OpError{
			Op:  "read",
			Net: "tcp",
			Err: &os.SyscallError{
				Syscall: "read",
				Err:     syscall.ECONNRESET,
			},
		}
		result := ClassifyError(err)
		if result != ErrorTypeConnectionReset {
			t.Errorf("expected %s, got %s", ErrorTypeConnectionReset, result)
		}
	})
}

func TestClassifyError_Other(t *testing.T) {
	err := errors.New("some random error")
	result := ClassifyError(err)
	if result != ErrorTypeOther {
		t.Errorf("expected %s, got %s", ErrorTypeOther, result)
	}
}
