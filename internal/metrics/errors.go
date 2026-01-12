package metrics

import (
	"crypto/tls"
	"net"
	"os"
	"strings"
	"syscall"
)

func ClassifyError(err error) string {
	if err == nil {
		return ""
	}

	errStr := strings.ToLower(err.Error())

	if os.IsTimeout(err) || strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		return ErrorTypeTimeout
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return ErrorTypeTimeout
	}

	if opErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
			switch sysErr.Err {
			case syscall.ECONNREFUSED:
				return ErrorTypeConnectionRefused
			case syscall.ECONNRESET:
				return ErrorTypeConnectionReset
			}
		}
	}

	if strings.Contains(errStr, "connection refused") {
		return ErrorTypeConnectionRefused
	}

	if strings.Contains(errStr, "connection reset") {
		return ErrorTypeConnectionReset
	}

	if _, ok := err.(*net.DNSError); ok {
		return ErrorTypeDNS
	}
	if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "dns") {
		return ErrorTypeDNS
	}

	if _, ok := err.(*tls.CertificateVerificationError); ok {
		return ErrorTypeTLS
	}
	if strings.Contains(errStr, "tls") || strings.Contains(errStr, "certificate") || strings.Contains(errStr, "x509") {
		return ErrorTypeTLS
	}

	if strings.Contains(errStr, "proxy authentication") || strings.Contains(errStr, "407") {
		return ErrorTypeProxyAuth
	}

	return ErrorTypeOther
}
