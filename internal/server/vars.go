package server

import (
	"net/http"
	"sync"

	"github.com/henvic/httpretty"
	"github.com/mbndr/logo"
)

var (
	handler      *Proxy
	server       *http.Server
	dump         *httpretty.Logger
	mime         = "text/plain"
	log          *logo.Logger
	ok           int64
	currentProxy string
	mutex        = sync.Mutex{}
)
