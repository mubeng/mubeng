package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/logrusorgru/aurora"
	"github.com/mubeng/mubeng/common"
	"github.com/mubeng/mubeng/pkg/helper"
	"github.com/mubeng/mubeng/pkg/mubeng"
	"github.com/sourcegraph/conc/pool"
)

// Do checks proxy from list.
//
// Displays proxies that have died if verbose mode is enabled,
// or save live proxies into user defined files.
func Do(opt *common.Options) {
	p := pool.New().WithMaxGoroutines(opt.Goroutine)
	c := retryablehttp.NewClient()
	c.RetryMax = opt.MaxRetries
	c.Logger = nil

	for _, proxy := range opt.ProxyManager.Proxies {
		address := helper.EvalFunc(proxy)

		p.Go(func() {
			addr, err := check(c, address, opt.Timeout)
			if len(opt.Countries) > 0 && !isMatchCC(opt.Countries, addr.Country) {
				return
			}

			if err != nil {
				if opt.Verbose && opt.OutputFormat == "" {
					fmt.Printf("[%s] %s\n", aurora.Red("DIED"), address)
				}
			} else {
				resultOutput := address

				if opt.OutputFormat != "" {
					proxyInfo := parseProxyAddr(address)
					proxyInfo.IPInfo = addr
					output := formatOutput(opt.OutputFormat, proxyInfo)
					resultOutput = output
					fmt.Println(output)
				} else {
					fmt.Printf(
						"[%s] [%s] [%s] %s (%s)\n",
						aurora.Green("LIVE"), aurora.Magenta(addr.Country),
						aurora.Cyan(addr.IP), address, aurora.Yellow(addr.Duration),
					)
				}

				if opt.Output != "" {
					fmt.Fprintf(opt.Result, "%s\n", resultOutput)
				}
			}
		})
	}

	p.Wait()
}

func isMatchCC(cc []string, code string) bool {
	if code == "" {
		return false
	}

	for _, c := range cc {
		if code == strings.ToUpper(strings.TrimSpace(c)) {
			return true
		}
	}

	return false
}

func check(c *retryablehttp.Client, address string, timeout time.Duration) (IPInfo, error) {
	req, err := retryablehttp.NewRequest("GET", endpoint, nil)
	if err != nil {
		return ipinfo, err
	}
	req.Header.Add("Connection", "close")

	tr, err := mubeng.Transport(address)
	if err != nil {
		return ipinfo, err
	}

	proxy := &mubeng.Proxy{
		Address:   address,
		Transport: tr,
	}

	client, err := proxy.New(req.Request)
	if err != nil {
		return ipinfo, err
	}

	c.HTTPClient = client
	c.HTTPClient.Timeout = timeout

	start := time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return ipinfo, err
	}
	duration := time.Since(start)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ipinfo, err
	}

	err = json.Unmarshal(body, &ipinfo)
	if err != nil {
		return ipinfo, err
	}
	ipinfo.Duration = duration.Truncate(time.Millisecond)

	defer resp.Body.Close()
	defer tr.CloseIdleConnections()

	return ipinfo, nil
}
