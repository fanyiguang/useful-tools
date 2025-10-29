VERSION=2.2.0
APPNAME=useful-tools
UPGARDENAME=upgrade

# build all -----------------------------------------------------------

build: build-mac-arm64 build-mac-amd64 build-windows-amd64 zip

# package mac arm64 -----------------------------------------------------------

build-mac-arm64: proc-mac-arm64 build-upgrade-mac-arm64 package-mac-arm64 copy-upgrade-arm64 mv-mac-arm64-pkg
proc-mac-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build \
	-trimpath \
	-ldflags "-s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=release" \
	-o ./bin/darwin/arm64/$(APPNAME) ./cmd/usefultools

build-upgrade-mac-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build \
	-trimpath \
	-ldflags "-s -w" \
	-o ./bin/darwin/arm64/$(UPGARDENAME) ./cmd/upgrade

package-mac-arm64:
	rm -rf ./useful-tools.app
	fyne package -os darwin -release -icon ./resource/icon.png --exe ./bin/darwin/arm64/$(APPNAME) --name useful-tools

copy-upgrade-arm64:
	cp ./bin/darwin/arm64/$(UPGARDENAME) ./useful-tools.app/Contents/MacOS/

mv-mac-arm64-pkg:
	cp -r ./useful-tools.app ./bin/darwin/arm64
	rm -rf ./useful-tools.app

# cross mac amd64 -------------------------------------------------------

build-mac-amd64: cross-mac-amd64-upgrade mkdir-darwin-amd64 mv-darwin-amd64-upgrade cross-mac-amd64-tool mv-mac-amd64-tool package-mac-amd64 mv-upgrade-to-darwin-amd64-pkg mv-mac-amd64-pkg

cross-mac-amd64-tool:
	fyne-cross darwin -arch=amd64 -icon ./resource/icon.png -name $(APPNAME) -app-id com.useful-tools.app -app-build 1 -app-version 2.1.0 -ldflags "-s -w -X 'useful-tools/common/config.Version=$(VERSION)' -X 'useful-tools/common/config.env=release'" ./cmd/usefultools

cross-mac-amd64-upgrade:
	fyne-cross darwin -arch=amd64 -icon ./resource/icon.png -app-id com.upgrade.app -name $(UPGARDENAME) -ldflags "-s -w -X useful-tools/common/config.Version=$(VERSION) -X useful-tools/common/config.env=release" ./cmd/upgrade

mkdir-darwin-amd64:
	mkdir -p ./bin/darwin/amd64

mv-darwin-amd64-upgrade:
	mv ./fyne-cross/bin/darwin-amd64/$(UPGARDENAME) ./bin/darwin/amd64

mv-upgrade-to-darwin-amd64-pkg:
	mv ./bin/darwin/amd64/$(UPGARDENAME) ./useful-tools.app/Contents/MacOS

mv-mac-amd64-tool:
	mv ./fyne-cross/dist/darwin-amd64/useful-tools.app/Contents/MacOS/usefultools ./bin/darwin/amd64/useful-tools

package-mac-amd64:
	rm -rf ./useful-tools.app
	fyne package -os darwin -release -icon ./resource/icon.png --exe ./bin/darwin/amd64/$(APPNAME) --name useful-tools

mv-mac-amd64-pkg:
	cp -r ./useful-tools.app ./bin/darwin/amd64
	rm -rf ./useful-tools.app

# cross windows amd64 -------------------------------------------------------

build-windows-amd64: cross-windows-amd64-tool mkdir-windows-amd64 mv-windows-amd64-tool cross-windows-amd64-upgrade mv-windows-amd64-upgrade

cross-windows-amd64-tool:
	GOFLAGS="-ldflags=-X=useful-tools/common/config.Version=$(VERSION) -X=useful-tools/common/config.env=release" fyne-cross windows -arch=amd64 -icon ./resource/icon.png -name $(APPNAME).exe -app-id com.useful-tools.app -app-build 1 -app-version 2.0.0 ./cmd/usefultools

cross-windows-amd64-upgrade:
	fyne-cross windows -arch=amd64 -icon ./resource/icon.png -app-id com.useful-tools.app -name $(UPGARDENAME).exe ./cmd/upgrade

mkdir-windows-amd64:
	mkdir -p ./bin/windows/amd64

mv-windows-amd64-tool:
	mv ./fyne-cross/bin/windows-amd64/$(APPNAME).exe ./bin/windows/amd64

mv-windows-amd64-upgrade:
	mv ./fyne-cross/bin/windows-amd64/$(UPGARDENAME).exe ./bin/windows/amd64

# zip pkg -----------------------------------------------------------------------------------

zip: zip-windows-amd64 zip-mac-amd64 zip-mac-arm64

zip-windows-amd64:
	cd ./bin/windows/amd64 && zip windows_amd64_useful-tools.zip -r ./$(APPNAME).exe ./$(UPGARDENAME).exe
zip-mac-amd64:
	cd ./bin/darwin/amd64 && zip darwin_amd64_useful-tools.zip -r ./useful-tools.app
zip-mac-arm64:
	cd ./bin/darwin/arm64 && zip darwin_arm64_useful-tools.zip -r ./useful-tools.app

# tidy -----------------------------------------------------------------------------------

tidy:
	go mod tidy
	go mod vendor
