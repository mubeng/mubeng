package common

var (
	// App name
	App = "mubeng"
	// Version of mubeng itself
	Version = ""
	// Author of mubeng
	Author = "@dwisiswant0"
	// Banner of mubeng
	Banner = `
           _   ` + Version + `
 _____ _ _| |_ ___ ___ ___
|     | | | . | -_|   | . |
|_|_|_|___|___|___|_|_|_  |
                      |___|
  by ` + Author
	// Usage of mubeng
	Usage = `
  mubeng [-c|-a :8080] -f file.txt [options...]

Options:
  GENERAL
    -f, --file <FILE>                Proxy file (required)
    -o, --output <FILE>              Write log output to FILE
    -t, --timeout <TIME>             Max. time allowed for connection (default: 30s)
    -u, --update                     Update mubeng to the latest stable version
    -v, --verbose                    Verbose mode
    -V, --version                    Show current mubeng version
        --max-retries <N>            Max. retries for failed HTTP requests (default: 0)

  PROXY CHECKER
    -c, --check                      Perform proxy check
    -g, --goroutine <N>              Max. goroutine to use (default: 50)
        --only-cc <AA>,<BB>          Only for specific country code (comma separated)

  IP ROTATOR
    -a, --address <ADDR>:<PORT>      Run proxy server
    -A, --auth <USER>:<PASS>         Set authorization for proxy server
    -d, --daemon                     Daemonize proxy server
    -m, --method <METHOD>            Rotation method (sequent/random) (default: sequent)
    -r, --rotate <N>                 Rotate proxy IP after N request (default: 1)
        --rotate-on-error            Rotate proxy IP and retry failed HTTP requests
        --remove-on-error            Remove proxy IP from proxy pool on failed HTTP requests
        --max-errors <N>             Max. errors allowed during rotation (default: 3)
                                     Use this with --rotate-on-error
                                     If value is less than 0 (e.g., -1), rotation will
                                     continue indefinitely
        --max-redirs <N>             Max. redirects allowed (default: 10)
    -s, --sync                       Syncrounus mode
    -w, --watch                      Watch proxy file, live-reload from changes
    	--type                       http or socks5 proxy server type ( default: http )

Examples:
  mubeng -f proxies.txt --check --output live.txt
  mubeng -a localhost:8080 -f live.txt -r 10 -w
  mubeng -a localhost:8080 -f live.txt -r 10 -w --type socks5

`
)
