package server

import (
	"net/http"
	"sync"

	"github.com/henvic/httpretty"
	"github.com/mbndr/logo"
)

var (
	handler *Proxy
	server  *http.Server
	dump    *httpretty.Logger
	mime    = "text/plain"
	log     *logo.Logger
	ok      = 0

	mutex = sync.Mutex{}
)
