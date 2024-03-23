SHELL         = /bin/sh

APP_NAME      = image-processor
VERSION      := $(shell git describe --always --tags)
GIT_COMMIT    = $(shell git rev-parse HEAD)
GIT_DIRTY     = $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE    = $(shell date '+%Y-%m-%d-%H:%M:%S')
CGO_ENABLED   = 1
GOARCH		  = amd64
GOOS		  = $(shell uname -s)

.PHONY: default
default: help

.PHONY: help
help:
	@echo 'Management commands for ${APP_NAME}:'
	@echo
	@echo 'Usage:'
	@echo '    make name                  Get APP Name.'
	@echo '    make lint				  Run static linter on a compiled project.'
	@echo '    make test                  Run tests on a compiled project.'
	@echo '    make sec			  		  Run security checks on a compiled project.'
	@echo '    make coverage			  See coverage detail.'
	@echo '    make run ARGS=             Run with supplied arguments.'
	@echo '    make build                 Compile the project.'
	@echo '    make clean                 Clean the directory tree.'
	@echo '    make prepare       		  Prepare your branch before make PR.'
	@echo

	.PHONY: get-app-name
get-app-name:
	@echo ${APP_NAME}

.PHONY: lint
lint:
	@echo "Check linter with staticcheck"
	staticcheck ./...

.PHONY: test
test:
	@echo "Testing ${APP_NAME} ${VERSION}"
	go test -race -coverprofile=coverage.out `go list ./... | grep -v util`
	go tool cover -func coverage.out

.PHONY: sec
sec:
	@echo "Check security issues with gosec"
	gosec -fmt=junit-xml -out=junit.xml -stdout -verbose=text -tests -exclude-dir=internal/pkg/util ./...

.PHONY:
coverage: test
	@echo "See coverage ${APP_NAME} ${VERSION}"
	go tool cover -html=coverage.out

.PHONY: run
run: build
	@echo "Running ${APP_NAME} ${VERSION}"
	bin/${APP_NAME} ${ARGS}

.PHONY: build
build:
	@echo "Building ${APP_NAME} ${VERSION}"
	go build -ldflags "-w -s -X github.com/NurfitraPujo/image-processor/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/NurfitraPujo/image-processor/version.Version=${VERSION} -X github.com/NurfitraPujo/image-processor/version.Environment=${APP_ENV} -X github.com/NurfitraPujo/image-processor/version.BuildDate=${BUILD_DATE}" -o bin/${APP_NAME} -trimpath .

.PHONY: clean
clean:
	@echo "Removing ${APP_NAME} ${VERSION}"
	@test ! -e bin/${APP_NAME} || rm bin/${APP_NAME}

.PHONY: prepare
prepare: lint sec test build
	@echo "Your works ready to reviewed! Go to make the PR."
