BINARY=genstruct_api
MAIN_FILE=main.go

.PHONY: build build_64linux test clean

build: 
	go fmt ./...
	golint ./...
	go vet ./...
	go build -o bin/${BINARY} ${MAIN_FILE}

build_64linux:
	go fmt ./...
	golint ./...
	go vet ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/${BINARY} ${MAIN_FILE}

test:
	go test -v ./... -cover

clean:
	@if [ -f bin/${BINARY} ] ; then rm bin/${BINARY} ; fi
