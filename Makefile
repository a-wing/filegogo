OS=
ARCH=
PROFIX=
NAME=filegogo
BINDIR=bin
VERSION=$(shell git describe --tags || git rev-parse --short HEAD || echo "unknown version")
BUILD_TIME=$(shell date +%FT%T%z)
LD_FLAGS='-X "filegogo/version.Version=$(VERSION)" -X "filegogo/version.BuildTime=$(BUILD_TIME)"'
GOBUILD=CGO_ENABLED=0 \
				go build -ldflags $(LD_FLAGS)

PLATFORM_LIST = \
								darwin-amd64 \
								linux-386 \
								linux-amd64 \
								linux-armv7 \
								linux-armv8 \
								freebsd-amd64

WINDOWS_ARCH_LIST = \
										windows-386 \
										windows-amd64
default: data
	GOOS=$(OS) GOARCH=$(ARCH) $(GOBUILD) -o $(NAME)

install:
	install -Dm755 ${NAME} -t ${PROFIX}/usr/bin/
	install -Dm644 conf/config.json -t ${PROFIX}/etc/${NAME}/
	install -Dm644 conf/${NAME}.service -t ${PROFIX}/lib/systemd/system/

all: linux-amd64 darwin-amd64 windows-amd64 # Most used

frontend:
	npm run build

data: frontend
	cp -r dist server/

test: data

run:
	go run -tags=dev main.go server

darwin-amd64: data
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-386: data
	GOARCH=386 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-amd64: data
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv7: data
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv8: data
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

freebsd-amd64: data
	GOARCH=amd64 GOOS=freebsd $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

windows-386: data
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

windows-amd64: data
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

releases: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

frontend-clean:
	rm dist

generate-clean:
	data/*_vfsdata.go

clean: frontend-clean generate-clean
	rm $(BINDIR)/*

