TARGETOS=linux
TARGETARCH=amd64

BIN_DIR=bin
CRAWLER_BIN=${BIN_DIR}/webcrawler
CRAWLER_MAIN=main.go
CRAWLER_SOURCES=$(shell find . -name '*.go')

.PHONY: all
all: build

.PHONY: clean
clean:
	rm -rf ${BIN_DIR} 

.PHONY: build
build: ${CRAWLER_BIN}

.PHONY: run
run: ${CRAWLER_MAIN} ${CRAWLER_SOURCES}
	go run -race . ${ARG}

${CRAWLER_BIN}: ${CRAWLER_MAIN} ${CRAWLER_SOURCES}
	go mod tidy	
	go mod download
	go test ./...
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o $@ . 
