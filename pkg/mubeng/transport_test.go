package mubeng

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"

	"github.com/go-test/deep"
	"h12.io/socks"
)

func TestTransport(t *testing.T) {
	type args struct {
		p string
	}

	failProxy := "gopher://localhost:70"
	httpProxy := "http://localhost:80"
	httpsProxy := "https://localhost:443"
	socks4Proxy := "socks4://localhost:5678"
	socks5Proxy := "socks5://localhost:3128"
	awsProxy := "aws://AKIAIOSFODNN7EXAMPLE:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY@ap-southeast-1"

	httpURL, err := url.Parse(httpProxy)
	if err != nil {
		t.Fatal(err)
	}

	httpsURL, err := url.Parse(httpsProxy)
	if err != nil {
		t.Fatal(err)
	}

	_, err = url.Parse(socks5Proxy)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		args    args
		wantTr  *http.Transport
		wantErr bool
	}{
		{
			name: "Switch-transport to HTTP",
			args: args{p: httpProxy},
			wantTr: &http.Transport{
				Proxy:             http.ProxyURL(httpURL),
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					MinVersion:         tls.VersionTLS10,
					CipherSuites:       getUnsafeCipherSuites(),
				},
			},
			wantErr: false,
		},
		{
			name: "Switch-transport to HTTPS",
			args: args{p: httpProxy},
			wantTr: &http.Transport{
				Proxy:             http.ProxyURL(httpsURL),
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					MinVersion:         tls.VersionTLS10,
					CipherSuites:       getUnsafeCipherSuites(),
				},
			},
			wantErr: false,
		},
		{
			name: "Switch-transport to SOCKSv5",
			args: args{p: socks5Proxy},
			wantTr: &http.Transport{
				Dial:              socks.Dial(socks5Proxy),
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					MinVersion:         tls.VersionTLS10,
					CipherSuites:       getUnsafeCipherSuites(),
				},
			},
			wantErr: false,
		},
		{
			name: "Switch-transport to SOCKSv4",
			args: args{p: socks4Proxy},
			wantTr: &http.Transport{
				Dial:              socks.Dial(socks4Proxy),
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					MinVersion:         tls.VersionTLS10,
					CipherSuites:       getUnsafeCipherSuites(),
				},
			},
			wantErr: false,
		},
		{
			name:    "Unsupported proxy proxy protocol scheme",
			args:    args{p: failProxy},
			wantTr:  nil,
			wantErr: true,
		},
		{
			name:    "AWS proxy protocol scheme",
			args:    args{p: awsProxy},
			wantTr:  new(http.Transport),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTr, err := Transport(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(gotTr, tt.wantTr); diff != nil {
				t.Error(diff)
			}
		})
	}
}
