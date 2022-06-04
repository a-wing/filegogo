SHELL=/bin/bash
OS=
ARCH=
PROFIX=
NAME=filegogo
VERSION=$(shell git describe --tags || git rev-parse --short HEAD || echo "unknown version")
BUILD_TIME=$(shell date +%FT%T%z)
LD_FLAGS='-X "filegogo/version.Version=$(VERSION)" -X "filegogo/version.Date=$(BUILD_TIME)"'
GOBUILD=CGO_ENABLED=0 \
				go build -trimpath -ldflags $(LD_FLAGS)

.PHONY: default
default: data build

.PHONY: build
build:
	GOOS=$(OS) GOARCH=$(ARCH) $(GOBUILD) -o $(NAME) ./main.go

.PHONY: install
install:
	install -Dm755 ${NAME} -t ${PROFIX}/usr/bin/
	install -Dm644 conf/${NAME}.toml -t ${PROFIX}/etc/
	install -Dm644 conf/${NAME}.service -t ${PROFIX}/lib/systemd/system/

.PHONY: webapp
webapp:
	npm run build

.PHONY: data
data: webapp
	cp -r dist server/build

test-e2e: default
	pushd e2e && npm run test && popd

webapp-clean:
	rm -r webapp/build

data-clean:
	rm -r server/build

clean: webapp-clean data-clean
	go clean -cache

