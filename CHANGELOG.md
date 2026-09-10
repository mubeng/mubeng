# Changelog

All notable changes to this project should be documented in this file.

### Unreleased

- Patch
  - Forward retry and error-handling options to daemon services

### v0.11.0

- Minor
  - Add HTTPS proxy protocol support
  - Add helper function for proxy formatting

### v0.10.0

- Minor
  - Add proxy file `--watch`er with live-reload proxy manager

### v0.9.3

- Patch
  - Add error check before write response writer

### v0.9.2

- Patch
  - Fix missing some arguments for daemon mode

### v0.9.1

- Patch
  - Fix non-proxy handler
  - Add CA certificate endpoint

### v0.9.0

- Minor
  - Add only show specific country code using `--only-cc` flag

### v0.8.0

- Minor
  - Add proxy server authentication

### v0.7.0

- Minor
  - Add `-s/--sync` flag for rotation guarantor
- Patch
  - Remove connection checking to perform proxy checks (close #65)
  - Add handler for latest version check

### v0.6.1

- Minor
  - Add rotation method for proxy IP rotator
- Patch
  - Add support for SOCKSv4(A)

### v0.5.2

- Patch
  - Authorization support for SOCKSv5

### v0.5.1

- Patch
  - Fix starting index for ProxyManager

### v0.5.0

- Minor
  - Add proxy manager for proxy rotator.

### v0.4.5

- Patch
  - Fix unit test non-nill tls.Config

### v0.4.4

- Patch
  - Enable InsecureSkipVerify

### v0.4.3

- Patch
  - Fix unable to get latest for go-get binary version

### v0.4.2

- Minor
  - Add updater package.
- Patch
  - Deactivated connection checking & get latest release info on every parsing options.
  - Add version option to show current mubeng version.

### v0.3.0

- Patch
  - Fix unable to run proxy server in daemon at Windows.

### v0.2.1

- Patch
  - Fix break updater logic.

### v0.2.0

- Minor
  - Upgrade dependencies.

### v0.1.2

- Patch
  - Add linkerd flag to inject version on Goreleaser action.

### v0.1.1

- Minor
  - Add `-u/--update` option.
- Patch
  - Make verbose mode works in daemon.

### v0.0.1

- Initial release.