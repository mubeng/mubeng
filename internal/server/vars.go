package server

import (
	"net/http"
	"sync"

	"github.com/henvic/httpretty"
	"github.com/mbndr/logo"
	"github.com/mubeng/mubeng/internal/metrics"
)

var (
	handler        *Proxy
	server         *http.Server
	metricsServer  *metrics.Server
	metricsEnabled bool
	dump           *httpretty.Logger
	mime           = "text/plain"
	log            *logo.Logger
	ok             = 1

	mutex = sync.Mutex{}
)
