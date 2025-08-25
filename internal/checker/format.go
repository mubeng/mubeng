package checker

import (
	"net/url"
	"strings"

	"github.com/valyala/fasttemplate"
)

// ProxyInfo contains parsed proxy address components and IP information
type ProxyInfo struct {
	Proxy    string // Full proxy URL (protocol://host:port)
	Protocol string // Protocol scheme (http, https, socks4, socks5, etc.)
	Host     string // Host/IP from proxy address
	Port     string // Port from proxy address
	IPInfo   IPInfo // IP information from the check
}

// parseProxyAddr extracts components from a proxy address
func parseProxyAddr(address string) ProxyInfo {
	info := ProxyInfo{
		Proxy: address,
	}

	parsedURL, err := url.Parse(address)
	if err != nil {
		parts := strings.Split(address, ":")
		if len(parts) >= 2 {
			info.Host = strings.Join(parts[:len(parts)-1], ":")
			info.Port = parts[len(parts)-1]
			// Try to guess protocol from common patterns
			if strings.Contains(strings.ToLower(address), "socks") {
				info.Protocol = "socks5" // default assumption
			} else {
				info.Protocol = "http" // default assumption
			}
		}

		return info
	}

	info.Protocol = parsedURL.Scheme
	info.Host = parsedURL.Hostname()

	port := parsedURL.Port()
	if port == "" {
		switch parsedURL.Scheme {
		case "http":
			port = "8080"
		case "https":
			port = "8080"
		case "socks4", "socks4a":
			port = "1080"
		case "socks5":
			port = "1080"
		}
	}
	info.Port = port

	return info
}

// formatOutput formats the output using the provided template
func formatOutput(template string, proxyInfo ProxyInfo) string {
	vars := map[string]any{
		"proxy":    proxyInfo.Proxy,
		"protocol": proxyInfo.Protocol,
		"host":     proxyInfo.Host,
		"port":     proxyInfo.Port,
		"ip":       proxyInfo.IPInfo.IP,
		"country":  proxyInfo.IPInfo.Country,
		"city":     proxyInfo.IPInfo.City,
		"org":      proxyInfo.IPInfo.Org,
		"region":   proxyInfo.IPInfo.Region,
		"timezone": proxyInfo.IPInfo.Timezone,
		"loc":      proxyInfo.IPInfo.Loc,
		"hostname": proxyInfo.IPInfo.Hostname,
		"duration": proxyInfo.IPInfo.Duration.String(),
	}

	t := fasttemplate.New(template, "{{", "}}")
	result := t.ExecuteStringStd(vars)

	return result
}
