VERSION=2.0.0
APPNAME=useful-tools
UPGARDENAME=upgrade

build-all: build-win64-release build-upgrade

build-win64-test:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=test" \
	-o ./bin/windows/test/amd64/$(APPNAME).exe ./cmd/usefultools

build-win64-release:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=release" \
	-o ./bin/windows/release/amd64/$(APPNAME).exe ./cmd/usefultools

build-mac-amd64-test:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=test" \
	-o ./bin/darwin/test/amd64/$(APPNAME) ./cmd/usefultools

build-mac-amd64-release:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=release" \
	-o ./bin/darwin/release/amd64/$(APPNAME) ./cmd/usefultools

build-mac-arm64-test:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=test" \
	-o ./bin/darwin/test/arm64/$(APPNAME) ./cmd/usefultools

build-mac-arm64-release:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=release" \
	-o ./bin/darwin/release/arm64/$(APPNAME) ./cmd/usefultools

build-upgrade:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
	--trimpath \
	-ldflags "-H windowsgui" \
	-o ./bin/windows/amd64/$(UPGARDENAME).exe ./cmd/upgrade