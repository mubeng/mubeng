package server

import (
	"errors"
	"fmt"
	"github.com/mubeng/mubeng/common"
	"github.com/mubeng/mubeng/internal/proxymanager"
	"github.com/mubeng/mubeng/pkg/mubeng"
	"math/rand/v2"
	"testing"
)

func TestRotateProxyIsEmptyString(t *testing.T) {
	proxies := []string{
		"http://1.1.1.1:6580/",
		"http://2.2.2.2:8274/",
		"http://3.3.3.3:5455/",
	}
	rotateAfter := rand.IntN(10)
	if rotateAfter <= 0 {
		rotateAfter = 1
	}

	p := &Proxy{
		Options: &common.Options{
			ProxyManager: &proxymanager.ProxyManager{},
			Rotate:       rotateAfter,
			Method:       "sequent",
		},
	}
	for _, proxy := range proxies {
		_, err := mubeng.Transport(proxy)
		if err == nil || errors.Is(err, mubeng.ErrSwitchTransportAWSProtocolScheme) {
			p.Options.ProxyManager.Proxies = append(p.Options.ProxyManager.Proxies, proxy)
		}
	}

	for _, proxy := range proxies {
		t.Run(fmt.Sprintf("proxy-%s", proxy), func(t *testing.T) {
			for i := 0; i < rotateAfter; i += 1 {
				proxy := p.rotateProxy()
				if proxy == "" {
					t.Error("Expected a non-empty proxy address, got empty string")
					return
				}
				t.Log(proxy)
			}
		})
	}
}
