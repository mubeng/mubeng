package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/mubeng/mubeng/common"
	"github.com/mubeng/mubeng/internal/metrics"
	"github.com/mubeng/mubeng/internal/proxygateway"
	"github.com/mubeng/mubeng/pkg/helper/awsurl"
	"github.com/mubeng/mubeng/pkg/mubeng"
)

type requestResult struct {
	response    *http.Response
	err         error
	proxy       string
	retryCount  int
	startTime   time.Time
}

// onRequest handles client request
func (p *Proxy) onRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if p.Options.Sync {
		mutex.Lock()
		defer mutex.Unlock()
	}

	if (req.URL.Scheme != "http") && (req.URL.Scheme != "https") {
		return req, serverErr(req)
	}

	if metricsEnabled {
		metrics.ActiveConnections.Inc()
		defer metrics.ActiveConnections.Dec()
	}

	resChan := make(chan requestResult)

	go func(r *http.Request) {
		log.Debugf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		result := requestResult{startTime: time.Now()}
		i := 0
		for {
			proxy := p.rotateProxy()
			result.proxy = proxy

			retryablehttpClient, err := p.getClient(r, proxy)
			if err != nil {
				result.err = err
				result.retryCount = i
				resChan <- result
				return
			}

			retryablehttpRequest, err := retryablehttp.FromRequest(r)
			if err != nil {
				result.err = err
				result.retryCount = i
				resChan <- result
				return
			}

			resp, err := retryablehttpClient.Do(retryablehttpRequest)
			if err != nil {
				if metricsEnabled {
					metrics.RetriesTotal.WithLabelValues(proxy).Inc()
				}

				if i >= p.Options.MaxErrors && p.Options.MaxErrors >= 0 {
					result.err = err
					result.retryCount = i
					resChan <- result
					return
				}

				if p.Options.RemoveOnErr {
					p.removeProxy(proxy)

					log.Debugf(
						"%s Removing proxy IP from proxy pool [proxies=%q]",
						r.RemoteAddr, fmt.Sprint(p.Options.ProxyManager.Count()),
					)
				}

				if p.Options.RotateOnErr && (i < p.Options.MaxErrors || p.Options.MaxErrors <= 0) {
					remaining := fmt.Sprint(p.Options.MaxErrors - i)
					if p.Options.MaxErrors <= 0 {
						remaining = "∞"
					}

					log.Debugf(
						"%s Retrying (rotated) %s %s [remaining=%q]",
						r.RemoteAddr, r.Method, r.URL, remaining,
					)

					i++
					continue
				} else {
					result.err = err
					result.retryCount = i
					resChan <- result
					return
				}
			}
			defer resp.Body.Close()

			buf, err := io.ReadAll(resp.Body)
			if err != nil {
				result.err = err
				result.retryCount = i
				resChan <- result
				return
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(buf))

			result.response = resp
			result.retryCount = i
			resChan <- result
			return
		}
	}(req)

	var resp *http.Response

	result := <-resChan
	duration := time.Since(result.startTime).Seconds()

	if result.response != nil {
		resp = result.response
		log.Debug(req.RemoteAddr, " ", resp.Status)

		if metricsEnabled {
			statusCode := strconv.Itoa(resp.StatusCode)
			metrics.RequestsTotal.WithLabelValues(req.Method, statusCode, result.proxy).Inc()
			metrics.RequestDuration.WithLabelValues(req.Method, result.proxy).Observe(duration)
			metrics.ProxyRequestsTotal.WithLabelValues(result.proxy, "success").Inc()
		}
	} else if result.err != nil {
		log.Errorf("%s %s", req.RemoteAddr, result.err)
		resp = serverErr(req)

		if metricsEnabled {
			errorType := metrics.ClassifyError(result.err)
			metrics.RequestErrorsTotal.WithLabelValues(errorType, result.proxy).Inc()
			metrics.RequestsTotal.WithLabelValues(req.Method, "502", result.proxy).Inc()
			metrics.RequestDuration.WithLabelValues(req.Method, result.proxy).Observe(duration)
			metrics.ProxyRequestsTotal.WithLabelValues(result.proxy, "error").Inc()
		}
	}

	return req, resp
}

// onConnect handles CONNECT method
func (p *Proxy) onConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	if p.Options.Auth != "" {
		auth := ctx.Req.Header.Get("Proxy-Authorization")
		if auth != "" {
			creds := strings.SplitN(auth, " ", 2)
			if len(creds) != 2 {
				return goproxy.RejectConnect, host
			}

			auth, err := base64.StdEncoding.DecodeString(creds[1])
			if err != nil {
				log.Warnf("%s: Error decoding proxy authorization", ctx.Req.RemoteAddr)
				return goproxy.RejectConnect, host
			}

			if string(auth) != p.Options.Auth {
				log.Errorf("%s: Invalid proxy authorization", ctx.Req.RemoteAddr)
				return goproxy.RejectConnect, host
			}
		} else {
			log.Warnf("%s: Unathorized proxy request to %s", ctx.Req.RemoteAddr, host)
			return goproxy.RejectConnect, host
		}
	}

	return goproxy.MitmConnect, host
}

// onResponse handles backend responses, and removing hop-by-hop headers
func (p *Proxy) onResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	for _, h := range mubeng.HopHeaders {
		resp.Header.Del(h)
	}

	return resp
}

func (p *Proxy) rotateProxy() string {
	var proxy string
	var err error

	if ok >= p.Options.Rotate {
		proxy, err = p.Options.ProxyManager.Rotate(p.Options.Method)
		if err != nil {
			log.Fatalf("Could not rotate proxy IP: %s", err)
		}

		if ok >= p.Options.Rotate {
			ok = 1
		}
	} else {
		ok++
	}

	return proxy
}

func (p *Proxy) removeProxy(target string) {
	err := p.Options.ProxyManager.RemoveProxy(target)
	if err != nil {
		log.Error(err)
		return
	}

	if metricsEnabled {
		metrics.ProxyRemovalsTotal.WithLabelValues(target).Inc()
		metrics.ProxyPoolSize.Set(float64(p.Options.ProxyManager.Count()))
	}
}

func (p *Proxy) getClient(req *http.Request, proxyAddr string) (*retryablehttp.Client, error) {
	tr, err := mubeng.Transport(proxyAddr)
	if err != nil && !errors.Is(err, mubeng.ErrSwitchTransportAWSProtocolScheme) {
		return nil, err
	}

	proxy := &mubeng.Proxy{
		Address:      proxyAddr,
		MaxRedirects: p.Options.MaxRedirects,
		Timeout:      p.Options.Timeout,
		Transport:    tr,
	}

	client, err := proxy.New(req)
	if err != nil {
		return nil, err
	}

	if awsurl.IsURL(proxyAddr) {
		var pg *proxygateway.ProxyGateway

		awsURL, err := awsurl.Parse(proxyAddr)
		if err != nil {
			return nil, err
		}

		_, err = awsURL.Credentials("")
		if err != nil {
			return nil, err
		}

		accessKeyID := awsURL.AccessKeyID
		secretAccessKey := awsURL.SecretAccessKey
		region := awsURL.Region

		baseURL, _, err := proxygateway.GetBaseURL(req.URL.String())
		if err != nil {
			return nil, err
		}

		gatewayKey := getGatewayKey(baseURL, region)

		if p.Gateways[gatewayKey] == nil {
			ctx := context.Background()
			gateway, err := proxygateway.New(ctx, accessKeyID, secretAccessKey, region)
			if err != nil {
				return nil, err
			}

			err = gateway.SetBaseURL(baseURL)
			if err != nil {
				return nil, err
			}

			err = gateway.Start(ctx)
			if err != nil {
				return nil, err
			}

			pg = gateway

			p.mu.Lock()
			p.Gateways[gatewayKey] = pg
			p.mu.Unlock()
		} else {
			pg = p.Gateways[gatewayKey]
		}

		// rewrite request URL to API Gateway endpoint URL
		gatewayEndpoint := pg.GetEndpoint()
		req.URL.Path = filepath.Join("/", proxygateway.StageName, req.URL.Path)
		req.URL.Host = gatewayEndpoint.Host
		req.URL.Scheme = gatewayEndpoint.Scheme
		req.Host = gatewayEndpoint.Host
	}

	if p.Options.Verbose {
		client.Transport = dump.RoundTripper(tr)
	}

	retryablehttpClient := mubeng.ToRetryableHTTPClient(client)
	retryablehttpClient.RetryMax = p.Options.MaxRetries
	retryablehttpClient.RetryWaitMin = client.Timeout
	retryablehttpClient.RetryWaitMax = client.Timeout
	retryablehttpClient.Logger = ReleveledLogo{
		Logger:  log,
		Request: req,
		Verbose: p.Options.Verbose,
	}

	return retryablehttpClient, nil
}

// nonProxy handles non-proxy requests
func nonProxy(w http.ResponseWriter, req *http.Request) {
	if common.Version != "" {
		w.Header().Add("X-Mubeng-Version", common.Version)
	}

	if req.URL.Path == "/cert" {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Disposition", fmt.Sprint("attachment; filename=", "goproxy-cacert.der"))
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(goproxy.GoproxyCa.Certificate[0]); err != nil {
			http.Error(w, "Failed to get proxy certificate authority.", 500)
			log.Errorf("%s %s %s %s", req.RemoteAddr, req.Method, req.URL, err.Error())
		}

		return
	}

	http.Error(w, "This is a mubeng proxy server. Does not respond to non-proxy requests.", 500)
}

func serverErr(req *http.Request) *http.Response {
	return goproxy.NewResponse(req, mime, http.StatusBadGateway, "Proxy server error")
}
