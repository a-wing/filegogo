SHELL=/bin/bash
OS=
ARCH=
PROFIX=
NAME=filegogo
VERSION=$(shell git describe --tags || git rev-parse --short HEAD || echo "unknown version")
COMMIT=$(shell git rev-parse HEAD || echo "unknown commit")
BUILD_TIME=$(shell date +%FT%T%z)
LD_FLAGS='\
				 -X "filegogo/version.Version=$(VERSION)" \
				 -X "filegogo/version.Commit=$(COMMIT)" \
				 -X "filegogo/version.Date=$(BUILD_TIME)" \
'

GOBUILD=CGO_ENABLED=0 \
				go build -tags release -trimpath -ldflags $(LD_FLAGS)

.PHONY: default
default: webapp build

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

test-e2e: default
	npm run test:e2e

webapp-clean:
	rm -r server/build

clean: webapp-clean
	go clean -cache

