.PHONY: default build docker-image test stop clean-images clean

BINARY = quarkcm

VERSION=
BUILD=

PKG            = github.com/CentaurusInfra/quarkcm
GOCMD          = go
BUILD_DATE     = `date +%FT%T%z`
GOFLAGS       ?= $(GOFLAGS:)
LDFLAGS       := "-X '$(PKG)/cmd.buildDate=$(BUILD_DATE)'"

default: build test

build:
	"$(GOCMD)" build ${GOFLAGS} -ldflags ${LDFLAGS} -o "${BINARY}"

docker-image:
	@docker build -t "${BINARY}" .

test:
	"$(GOCMD)" test -race -v ./...

stop:
	@docker stop "${BINARY}"

clean-images: stop
	@docker rmi "${BUILDER}" "${BINARY}"

clean:
	"$(GOCMD)" clean -i
