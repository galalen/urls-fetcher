BINARY_NAME=urls-fetcher

build:
	go build -o ${BINARY_NAME} main.go

build-mac:
	GOOS=darwin go build -o ${BINARY_NAME} main.go

build-linux:
	GOOS=linux go build -o ${BINARY_NAME} main.go

build-windows:
	GOOS=windows go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}