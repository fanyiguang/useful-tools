VERSION=2.1.0
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

build-upgrade:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
	--trimpath \
	-ldflags "-H windowsgui" \
	-o ./bin/windows/amd64/$(UPGARDENAME).exe ./cmd/upgrade

# mac arm64 -----------------------------------------------------------

build-mac-arm64: proc-mac-arm64 build-upgrade-mac-arm64 package-mac-arm64
proc-mac-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build \
	--trimpath \
	-ldflags "-s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=release" \
	-o ./bin/darwin/release/arm64/$(APPNAME) ./cmd/usefultools

build-upgrade-mac-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build \
	--trimpath \
	-ldflags "-s -w" \
	-o ./bin/darwin/arm64/$(UPGARDENAME) ./cmd/upgrade

package-mac-arm64:
	fyne package -os darwin -release -icon ./resource/icon.png --exe ./bin/darwin/release/arm64/$(APPNAME) --name useful-tools

# cross compile -------------------------------------------------------

cross-mac-arm64-tool:
	fyne-cross darwin -arch=arm64 -icon=./resource/icon.png -name=useful-tools -app-id=com.useful-tools.app  ./cmd/usefultools

cross-windows-amd64-tool:
	fyne-cross windows -arch=amd64 -icon ./resource/icon.png -name $(APPNAME).exe -app-id com.useful-tools.app -app-build 1 -app-version 2.0.0 ./cmd/usefultools
cross-windows-amd64-upgrade:
	fyne-cross windows -arch=amd64 -icon ./resource/icon.png -app-id com.useful-tools.app -name $(UPGARDENAME).exe ./cmd/upgrade