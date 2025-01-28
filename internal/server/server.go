package server

import (
	"context"
	"errors"
	"github.com/elazarl/goproxy"
	"github.com/henvic/httpretty"
	"github.com/mbndr/logo"
	"github.com/mubeng/mubeng/common"
	"github.com/mubeng/mubeng/internal/proxygateway"
	"github.com/things-go/go-socks5"
	netProxy "golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

// Run proxy server with a user defined listener.
//
// An active log have 2 receivers, especially stdout and into file if opt.Output isn't empty.
// Then close the proxy server if it receives a signal that interrupts the program.
func Run(opt *common.Options) {
	var recs []*logo.Receiver

	cli := logo.NewReceiver(os.Stderr, "")
	cli.Color = true
	cli.Level = logo.DEBUG
	recs = append(recs, cli)

	file, err := logo.Open(opt.Output)
	if err == nil {
		out := logo.NewReceiver(file, "")
		out.Format = "%s: %s"
		recs = append(recs, out)
	}

	dump = &httpretty.Logger{
		RequestHeader:  true,
		ResponseHeader: true,
		Colors:         true,
	}
	log = logo.NewLogger(recs...)
	handler = &Proxy{}
	handler.Options = opt
	if opt.Watch {
		watcher, err := opt.ProxyManager.Watch()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		go watch(watcher)
	}
	switch opt.Type {
	case "http":
		RunHTTPProxyServer(opt)
	case "socks5":
		RunSocks5ProxyServer(opt)
	}

}

func RunHTTPProxyServer(opt *common.Options) {
	handler.HTTPProxy = goproxy.NewProxyHttpServer()
	handler.HTTPProxy.OnRequest().DoFunc(handler.onRequest)
	handler.HTTPProxy.OnRequest().HandleConnectFunc(handler.onConnect)
	handler.HTTPProxy.OnResponse().DoFunc(handler.onResponse)
	handler.HTTPProxy.NonproxyHandler = http.HandlerFunc(nonProxy)
	handler.Gateways = make(map[string]*proxygateway.ProxyGateway)

	server = &http.Server{
		Addr:    opt.Address,
		Handler: handler.HTTPProxy,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go interrupt(stop)

	log.Infof("%d proxies loaded", opt.ProxyManager.Count())

	log.Infof("[PID: %d] Starting proxy server on %s", os.Getpid(), opt.Address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func RunSocks5ProxyServer(opt *common.Options) {
	socksOpts := []socks5.Option{
		socks5.WithLogger(log),
		socks5.WithDial(func(ctx context.Context, net_, addr string) (net.Conn, error) {
			proxy := handler.rotateProxy()
			parse, err := url.Parse(proxy)
			if err != nil {
				return nil, err
			}
			if parse.Host == "" {
				return nil, errors.New("socks5 address invalid")
			}
			password, _ := parse.User.Password()
			auth := &netProxy.Auth{
				User:     parse.User.Username(),
				Password: password,
			}
			dialer, err := netProxy.SOCKS5("tcp", parse.Host, auth, netProxy.Direct)
			if err != nil {
				return nil, err
			}
			log.Debugf("proxy server: %v target addr: %v", parse.Host, addr)
			return dialer.Dial(net_, addr)
		}),
	}
	if opt.Auth != "" {
		splitN := strings.SplitN(opt.Auth, ":", 2)
		credentials := socks5.StaticCredentials(map[string]string{splitN[0]: splitN[1]})
		socksOpts = append(socksOpts, socks5.WithCredential(credentials))
	}
	// get socks5 server
	socks5Server := socks5.NewServer(socksOpts...)
	log.Infof("%d proxies loaded", opt.ProxyManager.Count())
	log.Infof("[PID: %d] Starting proxy server on %s", os.Getpid(), opt.Address)
	if opt.Auth != "" {
		log.Infof("socks url: socks5://%v@%v", opt.Auth, opt.Address)
	} else {
		log.Infof("socks url: socks5://%v", opt.Address)
	}
	// start socks5 server
	if err := socks5Server.ListenAndServe("tcp", opt.Address); err != nil {
		log.Fatal(err)
	}
}
