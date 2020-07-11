# Makefile pour to build Gotwarden
PROJECT?=gitlab.com/tutilus/gotwarden
APP := $(shell basename "$(PWD)")
PORT?=3000
RELEASE?=0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Docker info
CONTAINER_IMAGE?=docker.io/tutilus/${APP}
DOCKER_REGISTRY=
GOOS?=linux
GOARCH?=amd64

## clean : Remove old binary
clean:
	rm -f ${APP}

## build: Build binary
build: clean
	CGO_ENABLED=1 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w  \
			-X ${APP}/version.Release=${RELEASE} \
			-X ${APP}/version.Commit=${COMMIT} \
			-X ${APP}/version.BuildTime=${BUILD_TIME}" \
			-o ${APP}

## run: Run application.
run-debug: build
	PORT=${PORT} ./${APP}

run: build
	PORT=${PORT} GIN_MODE=release ./${APP}

## container: Build application into a container.
container:
	docker build -t ${CONTAINER_IMAGE}:${RELEASE} .

## push: Push container into remote registry
push: container
	docker push ${DOCKER_REGISTRY}/${CONTAINER_IMAGE}:${RELEASE}

## serve: Start application in docker
serve: 
	docker stop ${APP}:${RELEASE} || true && docker rm ${APP}:${RELEASE} || true
	docker run --name ${APP} -p ${PORT}:${PORT} --rm \
		-e "PORT=${PORT}" \
		tutilus/${APP}:${RELEASE}

## test: Race tests
test:
	go test ./...

.PHONY: help

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo