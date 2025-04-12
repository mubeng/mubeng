package mubeng

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mubeng/mubeng/pkg/helper/awsurl"
	"h12.io/socks"
)

// Transport to auto-switch transport between HTTP/S or SOCKS v4(A) & v5 proxies.
//
// Depending on the protocol scheme, returning value of [http.Transport] with
// [http.Transport.Dialer] or [http.Transport.Proxy]. If protocol scheme is "aws",
// it will return default [http.Transport].
func Transport(p string) (*http.Transport, error) {
	var proxyURL *url.URL
	var err error

	tr := new(http.Transport)

	if awsurl.IsURL(p) {
		return tr, fmt.Errorf("%w: %w", ErrUnsupportedProxyProtocolScheme, ErrSwitchTransportAWSProtocolScheme)
	} else {
		proxyURL, err = url.Parse(p)
		if err != nil {
			return nil, err
		}
	}

	switch proxyURL.Scheme {
	case "socks4", "socks4a", "socks5":
		// TODO(dwisiswant0): deprecated, update this later.
		// nolint: staticcheck
		tr.Dial = socks.Dial(p)
	case "http", "https":
		tr.Proxy = http.ProxyURL(proxyURL)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProxyProtocolScheme, proxyURL.Scheme)
	}

	tr.DisableKeepAlives = true
	tr.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
		CipherSuites:       getUnsafeCipherSuites(),
	}

	return tr, nil
}
