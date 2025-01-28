package mubeng

// HopHeaders are meaningful only for a single transport-level connection, and are not stored by caches or forwarded by proxies.
var HopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Proxy-Connection",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

var ServerType string
