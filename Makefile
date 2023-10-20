VERSION=v1.1.0
APPNAME=useful-tools

build-all: build-release build-test

build-release:
	rm -rf ./_bin/windows/x86_64/release/$(APPNAME).exe
	CGO_ENABLED=1 GOOS=windows GOARCH=386 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION)" \
	-o ./_bin/windows/x86_64/release/$(APPNAME).exe ./

build-test:
	rm -rf ./_bin/windows/x86_64/test/$(APPNAME).exe
	CGO_ENABLED=1 GOOS=windows GOARCH=386 go build \
	--trimpath \
	-ldflags "-s -w -X useful-tools/common/config.Version=$(VERSION)" \
	-o ./_bin/windows/x86_64/test/$(APPNAME).exe ./