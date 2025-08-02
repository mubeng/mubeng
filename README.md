<h1 align="center">
  <br>
  <a href="#"><img src="https://user-images.githubusercontent.com/25837540/107883163-dbf70380-6f1f-11eb-856f-e69e733313e5.png" width="400px" alt="mubeng"></a>
</h1>

<h4 align="center">An incredibly fast proxy checker & IP rotator with ease.</h4>

<p align="center">
	<a href="https://golang.org"><img src="https://img.shields.io/badge/made%20with-Go-brightgreen"></a>
	<a href="https://goreportcard.com/report/github.com/mubeng/mubeng"><img src="https://goreportcard.com/badge/github.com/mubeng/mubeng"></a>
	<a href="https://github.com/mubeng/mubeng/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-yellowgreen"></a>
	<a href="#"><img src="https://img.shields.io/badge/platform-osx%2Flinux%2Fwindows-green"></a>
	<a href="https://github.com/mubeng/mubeng/releases"><img src="https://img.shields.io/github/release/mubeng/mubeng"></a>
	<a href="https://github.com/mubeng/mubeng/issues"><img src="https://img.shields.io/github/issues/mubeng/mubeng"></a>
</p>

<p align="center">
  <a href="https://github.com/mubeng/mubeng/blob/master/.github/CONTRIBUTING.md">Contributing</a> •
  <a href="https://github.com/mubeng/mubeng/blob/master/CHANGELOG.md">What's new</a> •
  <a href="https://pkg.go.dev/github.com/mubeng/mubeng/pkg/mubeng">Documentation</a> •
  <a href="https://github.com/mubeng/mubeng/issues/new/choose">Report Issues</a>
</p>

---

- [Features](#features)
- [Why mubeng?](#why-mubeng)
- [Installation](#installation)
  - [Binary](#binary)
  - [Docker](#docker)
  - [Source](#source)
- [Usage](#usage)
  - [Basic](#basic)
  - [Options](#options)
  	- [Notes](#notes)
  - [Install SSL Certificate](#install-ssl-certificate)
  - [Examples](#examples)
    - [Proxy checker](#proxy-checker)
    - [Proxy IP rotator](#proxy-ip-rotator)
      - [Burp Suite Upstream Proxy](#burp-suite-upstream-proxy)
      - [OWASP ZAP Proxy Chain](#owasp-zap-proxy-chain)
    - [Proxy format](#proxy-format)
    	- [Templating](#templating)
    	- [Amazon API Gateway](#amazon-api-gateway)
- [Limitations](#limitations)
	- [Known Bugs](#known-bugs)
- [Contributors](#contributors)
- [Pronunciation](#pronunciation)
- [Changes](#changes)
- [License](#license)

---

# Features

- **Proxy IP rotator**: Rotates your IP address for every specified request.
- **Proxy checker**: Verifies the availability of your proxy IPs.
- **Supports all HTTP/S methods**.
- **Compatible with HTTP, SOCKS v4(A), v5, and Amazon API Gateway** proxy protocols.
- **Passes all parameters and URIs**.
- **User-friendly**: Simply run it against your proxy file and select the desired action.
- **Cross-platform**: Runs seamlessly on Windows, Linux, Mac, and even Raspberry Pi.

# Why mubeng?

It's fairly simple, there is no need for additional configuration.

`mubeng` has 2 core functionality:

### 1. Run proxy server as proxy IP rotation

This is useful to avoid different kinds of IP ban, i.e. bruteforce protection, API rate-limiting or WAF blocking based on IP. We also leave it entirely up to user to use proxy pool resources from anywhere.

### 2. Perform proxy checks

So, you don't need any extra proxy checking tools out there if you want to check your proxy pool.

# Installation

## Binary

Simply, download a pre-built binary from [releases page](https://github.com/mubeng/mubeng/releases) and run!

## Docker

Pull the [Docker](https://docs.docker.com/get-docker/) image by running:

```bash
▶ docker pull ghcr.io/mubeng/mubeng:latest
```

## Source

Using [Go](https://golang.org/doc/install) compiler:

```bash
▶ go install -v github.com/mubeng/mubeng@latest
```

### — or

Manual building executable from source code:

```bash
▶ git clone https://github.com/mubeng/mubeng
▶ cd mubeng
▶ make build
▶ (sudo) install ./bin/mubeng /usr/local/bin
```

# Usage

For usage, it's always required to provide your proxy list, whether it is used to check or as a proxy pool for your proxy IP rotation.

<center>
  <a href="#"><img alt="kitabisa mubeng" src="https://github.com/user-attachments/assets/3c19e328-cfd7-43f7-bf83-b3996671fc67" width="80%"></a>
</center>

## Basic

```bash
▶ mubeng [-c|-a :8080] -f file.txt [options...]
```

## Options

Here are all the options it supports.

```bash
▶ mubeng -h
```

| **Flag**                      	| **Description**                                              	|
|-------------------------------  |-------------------------------------------------------------- |
| -f, --file `<FILE>`             | Proxy file.                                                   |
| -a, --address `<ADDR>:<PORT>`   | Run proxy server.                                             |
| -A, --auth `<USER>:<PASS>`      | Set authorization for proxy server.                           |
| -d, --daemon                    | Daemonize proxy server.                                       |
| -c, --check                     | To perform proxy live check.                                  |
| -g, --goroutine `<N>`           | Max. goroutine to use (default: 50).                          |
|     --only-cc `<AA>,<BB>`       | Only show specific country code (comma separated).            |
| -t, --timeout                   | Max. time allowed for proxy server/check (default: 30s).      |
| -r, --rotate `<AFTER>`          | Rotate proxy IP for every `AFTER` request (default: 1).       |
|     --rotate-on-error           | Rotate proxy IP and retry failed HTTP requests.               |
|     --remove-on-error           | Remove proxy IP from proxy pool on failed HTTP requests.      |
|     --max-errors `<N>`          | Max. errors allowed during rotation (default: 3).             |
|                                 | Use this with `--rotate-on-error`.                            |
|                                 | If value is less than 0 (e.g., -1), rotation will             |
|                                 | continue indefinitely.                                        |
|     --max-redirs `<N>`          | Max. redirects allowed (default: 10).                         |
|     --max-retries `<N>`         | Max. retries for failed HTTP requests (default: 0).           |
| -m, --method `<METHOD>`         | Rotation method (sequent/random) (default: sequent).          |
| -s, --sync                      | Sync will wait for the previous request to complete.          |
| -v, --verbose                   | Dump HTTP request/responses or show died proxy on check.      |
| -o, --output `<FILE>`           | Save output from proxy server or live check.                  |
| -u, --update                    | Update mubeng to the latest stable version.                   |
| -w, --watch                     | Watch proxy file, live-reload from changes.                   |
| -V, --version                   | Show current mubeng version.                                  |

<table>
	<td>
		<h4>NOTES:</h4>
		<ul>
			<li>Rotations are counted for all requests, even if the request fails.
				<!-- <ul>
					<li>The rotation is incremental starting at the beginning of the list.</li>
					<li>Rotation means random, <b>NOT</b> choosing a proxy after/increment from proxy pool. We do not set up conditions if a proxy has been used. So, there is no guarantee if your request reaches the <i>N</i> value <code>(-r/--rotate)</code> your IP proxy will rotate.</li>
				</ul> -->
			</li>
			<li>The proxy server runs asynchronously by default, so it doesn't guarantee that your requests after <i>N</i> (which is <i>N+1</i> and so on) will rotate the proxy IP, instead use the <code>-s/--sync</code> flag to wait for requests to the previous proxy to complete.</li>
			<li>Daemon mode <code>(-d/--daemon)</code> will install mubeng as a service on the (Linux/OSX) system/setting up callback (Windows).
				<ul>
					<li>Hence you can control service with <code>journalctl</code>, <code>service</code> or <code>net</code> (for Windows) command to start/stop proxy server.</li>
					<li>Whenever you activate the daemon mode, it works by forcibly stop and uninstalling the existing mubeng service, then re-install and starting it up in daemon.</li>
				</ul>
			</li>
			<li>Verbose mode <code>(-v/--verbose)</code> and timeout <code>(-t/--timeout)</code> apply to both proxy check and proxy IP rotation actions.</li>
			<li>HTTP traffic requests and responses is displayed when verbose mode <code>(-v/--verbose)</code> is enabled, but
				<ul>
					<li>We <b>DO NOT</b> explicitly display the request/response body, and</li>
					<li>All cookie values in headers will be redacted automatically.</li>
				</ul>
			</li>
			<li>If you use output option <code>(-o/--output)</code> to run proxy IP rotator, request/response headers are <b>NOT</b> written to the log file.</li>
			<li>A timeout option <code>(-t/--timeout)</code> value is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "5s", "300ms", "-1.5h" or "2h45m".
				<ul>
					<li>Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", and "h".</li>
				</ul>
			</li>
			<li>The max. errors <code>(--max-errors)</code> ≠ max. retries <code>(--max-retries)</code>.</li>
			<ul>
					<li>The max. retries <code>(--max-retries)</code> apply to retrying a failed HTTP request with the same proxy.</li>
					<li>Meanwhile, max. errors <code>(--max-errors)</code> refer to the total failed HTTP requests from the proxies used.</li>
					<li>For example, once the max. retries <code>(--max-retries)</code> are reached <i>(if supplied)</i>, the error is counted, and the proxy IP will rotate for the same failed HTTP request until the max. errors <code>(--max-errors)</code> are reached, whereas max. retries <code>(--max-retries)</code> <b>DO NOT</b> trigger a proxy rotation.</li>
				</ul>
			<li>If the value of max. errors <code>(--max-errors)</code> is less than <b>0</b> (e.g. <b>-1</b>), rotation will continue indefinitely.</li>
		</ul>
	</td>
</table>

## Install SSL Certificate

mubeng uses built-in certificate authority by [GoProxy](https://github.com/elazarl/goproxy). With mubeng proxy server running, the generated certificate can be exported by visiting `http://mubeng/cert` in a browser.

Installation steps for CA certificate is [similar to other](https://portswigger.net/burp/documentation/desktop/external-browser-config/certificate) proxy tools.

## Examples

For example, you've proxy pool `(proxies.txt)` as:

<table>
	<td>
		<pre>http://127.0.0.1:8080
https://127.0.0.1:443
socks4://127.0.0.1:4145
socks5://127.0.0.1:2121
...
...</pre>
	</td>
</table>

> Because we use auto-switch transport, `mubeng` can accept multiple proxy protocol schemes at once.<br>
> Please refer to [documentation](https://pkg.go.dev/github.com/mubeng/mubeng/pkg/mubeng#Transport) for this package.

### Proxy checker

Pass `--check` flag in command to perform proxy checks:

```bash
▶ mubeng -f proxies.txt --check --only-cc AU,US,UK --output live.txt
```

The above case also uses `--output` flag to save a live proxy of specific country code with `--only-cc` flag (`ISO-3166` alpha-2) into file _(live.txt)_ from checking result.

<p align="center">
  <img src="https://github.com/user-attachments/assets/97b9e5ad-64bd-403a-bd3a-1bc4eb7fb3f1" height="400"><br />
  <i>(Figure: Checking proxies mubeng with max. 5s timeout)</i>
</p>

#### Custom Output Format

You can customize the output format of the proxy checker using the `--output-format` flag with [fasttemplate](https://github.com/valyala/fasttemplate) syntax:

```bash
▶ mubeng -f proxies.txt --check --output-format "{{proxy}} | {{country}} | {{duration}}"
```

**Available template variables:**

| Variable | Description | Example |
|----------|-------------|---------|
| `{{proxy}}` | Full proxy URL | `http://192.168.1.1:8080` |
| `{{protocol}}` | Proxy protocol scheme | `http`, `https`, `socks5` |
| `{{host}}` | Proxy host/IP address | `192.168.1.1` |
| `{{port}}` | Proxy port | `8080` |
| `{{ip}}` | External IP address | `203.0.113.1` |
| `{{country}}` | Country code | `US`, `UK`, `AU` |
| `{{city}}` | City name | `New York`, `London` |
| `{{org}}` | Organization/ISP | `Google Inc.` |
| `{{region}}` | Region/State | `California`, `England` |
| `{{timezone}}` | Timezone | `America/New_York` |
| `{{loc}}` | Latitude,Longitude | `40.7128,-74.0060` |
| `{{hostname}}` | Hostname | `example.com` |
| `{{duration}}` | Response time | `245ms` |

**Examples:**

```bash
# JSON-like format
▶ mubeng -f proxies.txt --check --output-format '{"proxy":"{{proxy}}","country":"{{country}}","duration":"{{duration}}"}'

# CSV format
▶ mubeng -f proxies.txt --check --output-format "{{proxy}},{{country}},{{city}},{{duration}}"

# Custom detailed format
▶ mubeng -f proxies.txt --check --output-format "[{{country}}] {{proxy}} ({{org}}) - {{duration}}"
```

When using `--output-format`, the formatted output is applied to both console display and file output (when using `-o`/`--output` flag).

### Proxy IP rotator

Furthermore, if you wish to do proxy IP rotator from proxies that are still alive earlier from the results of checking `(live.txt)` _(or if you have your own list)_, you must use `-a` _(--address)_ flag instead to run proxy server:

```bash
▶ mubeng -a localhost:8089 -f live.txt -r 10 -m random
```

The `-r` _(--rotate)_ flag works to rotate your IP for every _N_ request value you provide `(10)`, and the `-m` _(--method)_ flag will rotate the proxy sequential/randomly.

<p align="center">
  <img src="https://github.com/user-attachments/assets/b86a5bf4-41e6-43ac-8077-9f82442e9577">
  <i>(Figure: Running mubeng as proxy IP rotator with verbose mode)</i>
</p>

### [Burp Suite](https://portswigger.net/burp/documentation/desktop/getting-started/installing-burp) Upstream Proxy

In case you want to use `mubeng` _(proxy IP rotator)_ as an upstream proxy in Burp Suite, acting in-between Burp Suite and mubeng to the internet, so you don't need any additional extensions in Burp Suite for that. To demonstrate this:

<p align="center">
  <img src="https://user-images.githubusercontent.com/25837540/107985702-24d0ba00-6ffd-11eb-9489-c19e52c921f5.jpg">
  <i>(Figure: Settings Burp Suite Upstream Proxy to mubeng)</i>
</p>

In your Burp Suite instance, select **Project options** menu, and click **Connections** tab. In the **Upstream Proxy Servers** section, check **Override user options** then press **Add** button to add your upstream proxy rule. After that, fill required columns _(Destination host, Proxy host & Proxy port)_ with correct details. Click **OK** to save settings.

### [OWASP ZAP](https://www.zaproxy.org/download/) Proxy Chain

It acts the same way when you using an upstream proxy. OWASP ZAP allows you to connect to another proxy for outgoing connections in OWASP ZAP session. To chain it with a mubeng proxy server:

<p align="center">
	<img src="https://user-images.githubusercontent.com/25837540/108060995-41670380-708a-11eb-83ad-c781421af473.png">
	<i>(Figure: Settings proxy chain connection in OWASP ZAP to mubeng)</i>
</p>


Select **Tools** in the menu bar in your ZAP session window, then select the **Options** _(shortcut: Ctrl+Alt+O)_ submenu, and go to **Connection** section. In that window, scroll to **Use proxy chain** part then check **Use an outgoing proxy server**. After that, fill required columns _(Address/Domain Name & Port)_ with correct details. Click **OK** to save settings.

### Proxy format

Currently mubeng supports HTTP(S) & SOCKSv4(A)/v5 protocol, see [examples](#examples) above. But, not limited by that we also support proxy string substitution and helper functions for your proxy pool.

#### Templating

If you have an authenticated proxy, you definitely don't want to write credentials constantly to the proxy pool file. **mubeng** can evaluate environment variable with `{{VARIABLE}}` writing style.

For example:

1. String substitute

```console
$ export USERNAME="FOO"
$ export PASSWORD="BAR"
$ echo "http://{{USERNAME}}:{{PASSWORD}}@192.168.0.1:31337" > list.txt
$ mubeng -f list.txt -a :8080
```

2. Helper function

Available functions currently supported:

- `uint32`, and
- `uint32n N`.

Those following above functions are thread-safe pseudo-randomness.

As an example of its use, we will be utilizing stream isolation over Tor SOCKS. With this method, you just need one Tor instance and each request can use a different stream with a different exit node, but that doesn't guarantee that your _ass_ will be rotated. Thus, we have to create unique `USER:PASS` pair to isolate streams for every connection. In order to pass pseudo-random proxy authorization, use `uint32` or `uint32n` function on your proxy pool, like:

```console
$ echo "socks5://{{uint32}}:{{uint32}}@127.0.0.1:9050" > list.txt
$ while :; do mubeng -f list.txt -c 2>/dev/null; done
[LIVE] [XX] [23.**.177.2] socks5://2123347975:3094119616@127.0.0.1:9050
[LIVE] [XX] [199.**.253.156] socks5://1646373938:2740927425@127.0.0.1:9050
[LIVE] [XX] [185.**.101.137] socks5://814036283:1382144874@127.0.0.1:9050
[LIVE] [XX] [185.**.83.83] socks5://2895805939:2276057153@127.0.0.1:9050
[LIVE] [XX] [103.**.167.10] socks5://408584795:1244204083@127.0.0.1:9050
[LIVE] [XX] [198.**.84.99] socks5://3015151335:251835794@127.0.0.1:9050
[LIVE] [XX] [179.**.159.197] socks5://3952852758:324998250@127.0.0.1:9050
^C
```

#### Amazon API Gateway

The mubeng proxy rotator also supports integration with Amazon API Gateway. This allows you to route traffic through multiple AWS regions for enhanced redundancy and geographic distribution.

Format for AWS proxy strings:

```
aws://AWS_ACCESS_KEY_ID:AWS_SECRET_ACCESS_KEY@REGION
```

> [!TIP]
> Since it uses a custom parser, the AWS secret access key (or any other parts) can be quoted for better readability. Example: `aws://AKIAIOSFODNN7EXAMPLE:"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"@us-west-1`.
> 
> This quoting feature only works for the **`aws`** protocol scheme.

To get started, you'll need to:

1. Export your AWS credentials as environment variables

```bash
export AWS_ACCESS_KEY_ID="AKIAIOSFODNN7EXAMPLE"
export AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" 
```

2. Create a proxy list file containing AWS regions

Generate proxy entries for multiple AWS regions:

```bash
echo "aws://{{AWS_ACCESS_KEY_ID}}:{{AWS_SECRET_ACCESS_KEY}}@"{us,eu}"-"{east,west}"-"{1,2} | tr ' ' '\n' > list.txt
```

This will create entries for regions like:

* `aws://{{AWS_ACCESS_KEY_ID}}:{{AWS_SECRET_ACCESS_KEY}}@us-east-1`
* `aws://{{AWS_ACCESS_KEY_ID}}:{{AWS_SECRET_ACCESS_KEY}}@us-east-2`
* `aws://{{AWS_ACCESS_KEY_ID}}:{{AWS_SECRET_ACCESS_KEY}}@us-west-1`
* ...

3. Start mubeng proxy server pointing to your AWS proxy list:

```bash
$ mubeng -f list.txt -a :8080
```

This setup enables mubeng to automatically rotate traffic through multiple AWS regions via API Gateway. When running the proxy server, mubeng will dynamically substitute your AWS credentials from environment variables using the [templating](#templating) feature described above.

> [!NOTE]
> Ensure your AWS credentials have the appropriate permissions to access API Gateway in the specified regions.

# Limitations

Currently IP rotation runs the proxy server only as an HTTP protocol, not a SOCKSv4(A)/v5 protocol, even though the resource you have is SOCKSv4(A)/v5. In other words, the SOCKSv4(A)/v5 resource that you provide is used properly because it uses auto-switch transport on the client, but this proxy server **DOES NOT** switch to anything other than HTTP protocol.

# Contributors

[![contributions](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/mubeng/mubeng/issues)

This project exists thanks to all the people who contribute. To learn how to setup a development environment and for contribution guidelines, see [CONTRIBUTING.md](https://github.com/mubeng/mubeng/blob/master/.github/CONTRIBUTING.md).

<a href="https://github.com/mubeng/mubeng/graphs/contributors">
	<img src=".github/CONTRIBUTORS.svg">
</a>

# Pronunciation

[`jv_ID`](https://www.localeplanet.com/java/jv-ID/index.html) • **/mo͞oˌbēNG/** — mubeng-mubeng nganti mumet. (ꦩꦸꦧꦺꦁ​ꦔꦤ꧀ꦠꦶ​ꦩꦸꦩꦺꦠ꧀)

# Changes

For changes, see [CHANGELOG.md](https://github.com/mubeng/mubeng/blob/master/CHANGELOG.md).

# License

This program is free software: you can redistribute it and/or modify it under the terms of the [Apache license](https://github.com/mubeng/mubeng/blob/master/LICENSE). mubeng and any contributions are copyright © by Dwi Siswanto 2021-2025.