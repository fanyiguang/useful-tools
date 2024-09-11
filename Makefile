VERSION=2.1.0
APPNAME=useful-tools
UPGARDENAME=upgrade

build-all: build-useful-release build-upgrade

build-useful-test:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=test" \
	-o ./bin/windows/test/amd64/$(APPNAME).exe ./cmd/usefultools

build-useful-release:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=release" \
	-o ./bin/windows/release/amd64/$(APPNAME).exe ./cmd/usefultools

build-upgrade:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
	--trimpath \
	-ldflags "-H windowsgui" \
	-o ./bin/windows/amd64/$(UPGARDENAME).exe ./cmd/upgrade

build-useful-2:
	CGO_ENABLED=1 GOOS=windows go build  \
	--trimpath \
	-o ./bin/windows/amd64/$(APPNAME).exe ./cmd/usefultools