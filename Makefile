VERSION=1.2.0
APPNAME=useful-tools
UPGARDENAME=upgrade

build-all: build-useful build-upgrade

build-useful:
	CGO_ENABLED=1 GOOS=windows GOARCH=386 go build \
	--trimpath \
	-ldflags "-H windowsgui -s -w -X useful-tools/common/config.Version=$(VERSION)" \
	-o ./bin/windows/386/$(APPNAME).exe ./

build-upgrade:
	CGO_ENABLED=1 GOOS=windows GOARCH=386 go build \
	--trimpath \
	-ldflags "-H windowsgui" \
	-o ./bin/windows/386/$(UPGARDENAME).exe ./cmd/upgrade

build-useful-2:
	CGO_ENABLED=1 GOOS=windows go build  \
	--trimpath \
	-o ./bin/windows/amd64/$(APPNAME).exe ./cmd/usefultools