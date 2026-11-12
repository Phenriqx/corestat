# Requires Go, Node.js/npm, Wails CLI, and the native Wails build dependencies.
# Install the CLI: go install github.com/wailsapp/wails/v2/cmd/wails@v2.10.1
WAILS ?= wails
# Wails defaults to WebKit 4.0; use 4.1 on Linux when only that version exists.
# Explicit WAILS_FLAGS values override this automatic default.
WAILS_FLAGS ?= $(shell \
	if [ "$$(uname -s)" = "Linux" ] && \
		! pkg-config --exists webkit2gtk-4.0 2>/dev/null && \
		pkg-config --exists webkit2gtk-4.1 2>/dev/null; then \
		printf '%s' '-tags webkit2_41'; \
	fi)

.DEFAULT_GOAL := build
.PHONY: build run test check-wails

# Build the frontend and package the desktop application in build/bin.
build: check-wails
	$(WAILS) build $(WAILS_FLAGS)

# Run the desktop application with frontend/backend hot reload.
run: check-wails
	$(WAILS) dev $(WAILS_FLAGS)

# Run all Go tests (packages without tests are accepted).
test:
	go test ./...

check-wails:
	@command -v $(WAILS) >/dev/null 2>&1 || { \
		echo "Wails CLI not found. Install it with:"; \
		echo "  go install github.com/wailsapp/wails/v2/cmd/wails@v2.10.1"; \
		echo 'Then add $$(go env GOPATH)/bin to your PATH.'; \
		exit 1; \
	}
