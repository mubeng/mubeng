all: mubeng

APP_NAME = mubeng
VERSION  = $(shell git describe --always --tags)
LINT_CMD = "golangci-lint run ./... -v --timeout 5m"
LINT_BIN = "https://install.goreleaser.com/github.com/golangci/golangci-lint.sh"
TEST_CMD = go test -short ./...

mubeng: test build

test:
	@if [ "$(VERBOSE)" = "1" ]; then \
		$(TEST_CMD); \
	else \
		output_file=$$(mktemp "$${TMPDIR:-/tmp}/mubeng-test.XXXXXX" 2>&1); \
		status=$$?; \
		if [ "$$status" -ne 0 ]; then \
			printf 'test: failed (exit %s)\n' "$$status"; \
			printf '%s\n' "$$output_file"; \
			exit "$$status"; \
		fi; \
		trap 'rm -f "$$output_file"' 0; \
		if $(TEST_CMD) >"$$output_file" 2>&1; then \
			printf '%s\n' 'test: ok'; \
		else \
			status=$$?; \
			printf 'test: failed (exit %s)\n' "$$status"; \
			cat "$$output_file"; \
			exit "$$status"; \
		fi; \
	fi

test-extra: golangci-lint test

build:
	@echo "Building ${APP_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	@mkdir -p bin/
	@go build -ldflags "-s -w -X github.com/mubeng/mubeng/common.Version=${VERSION}" -o ./bin/${APP_NAME} .


golangci-lint:
	@echo "Run GolangCI-Lint"
	@if [ -x $(command -v golangci-lint) ]; then \
		eval "${LINT_CMD}"; \
	else \
		echo "Download GolangCI-Lint..."; \
		curl -sfL "${LINT_BIN}" | sh; \
		eval "./bin/${LINT_CMD}"; \
	fi;

clean:
	@echo "Removing binaries"
	@rm -rf bin/