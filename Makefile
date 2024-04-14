VERSION = 0.0.1
GOFLAGS = -ldflags="-X 'github.com/epiefe/jswap/cmd.version=$(VERSION)'"

.PHONY: build
build:
	go build $(GOFLAGS) -o build/ jswap.go

linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o build/linux-amd64/ jswap.go

mac-amd64:
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o build/mac-amd64/ jswap.go

mac-arm64:
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o build/mac-arm64/ jswap.go

win-amd64:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o build/win-amd64/ jswap.go

win-installer: win-amd64
	makensis installer.nsi

all: linux-amd64 mac-amd64 mac-arm64 win-installer

clean:
	rm -r build
