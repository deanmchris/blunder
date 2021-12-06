BINARY_NAME=blunder

build:
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows.exe blunder/main.go
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin blunder/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux blunder/main.go

build-windows:
	set GOARCH=amd64&& set GOOS=windows&& go build -o ${BINARY_NAME}-windows.exe blunder/main.go
	set GOARCH=amd64&& set GOOS=darwin&& go build -o ${BINARY_NAME}-darwin blunder/main.go
	set GOARCH=amd64&& set GOOS=linux&& go build -o ${BINARY_NAME}-linux blunder/main.go

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows