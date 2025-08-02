package common

import (
	"os"
	"time"

	"github.com/mubeng/mubeng/internal/proxymanager"
)

// Options consists of the configuration required.
type Options struct {
	ProxyManager *proxymanager.ProxyManager
	Result       *os.File
	Timeout      time.Duration

	Address      string
	Auth         string
	CC           string
	Check        bool
	Countries    []string
	Daemon       bool
	File         string
	Goroutine    int
	Method       string
	Output       string
	OutputFormat string
	Rotate       int
	RotateOnErr  bool
	RemoveOnErr  bool
	Sync         bool
	Verbose      bool
	Watch        bool
	MaxErrors    int
	MaxRedirects int
	MaxRetries   int
}
