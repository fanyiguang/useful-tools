VERSION=1.2.0
APPNAME=useful-tools
UPGARDENAME=upgrade

build-all: build-useful build-upgrade

build-useful:
	CGO_ENABLED=1 GOOS=windows GOARCH=386 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION)" \
	-o ./_bin/windows/386/$(APPNAME).exe ./

build-upgrade:
	CGO_ENABLED=1 GOOS=windows GOARCH=386 go build \
	--trimpath \
	-ldflags "-H windowsgui" \
	-o ./_bin/windows/386/$(UPGARDENAME).exe ./cmd/upgrade