VERSION = $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.0.0)
COMMIT := $(shell git rev-list --abbrev-commit --all --max-count=1)

GOFLAGS = -ldflags="-X 'github.com/epiefe/jswap/cmd.version=$(VERSION) (commit $(COMMIT))'"

.PHONY: build
build:
	go build $(GOFLAGS) -o build/ jswap.go

.PHONY: linux-amd64
linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o build/linux-amd64/ jswap.go

.PHONY: mac-amd64
mac-amd64:
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o build/mac-amd64/ jswap.go

.PHONY: mac-arm64
mac-arm64:
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o build/mac-arm64/ jswap.go

.PHONY: win-amd64
win-amd64:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o build/win-amd64/ jswap.go

.PHONY: win-installer
win-installer: win-amd64
	makensis -Dversion=$(VERSION) installer.nsi

.PHONY: all
all: linux-amd64 mac-amd64 mac-arm64 win-installer

.PHONY: clean
clean:
	rm -r build
