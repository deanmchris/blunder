BINARY_NAME=blunder-8.5.5

build:
	GOARCH=amd64 GOAMD64=v1 go build -o ${BINARY_NAME}-default blunder/main.go
	GOARCH=amd64 GOAMD64=v2 go build -o ${BINARY_NAME}-popcnt blunder/main.go
	GOARCH=amd64 GOAMD64=v3 go build -o ${BINARY_NAME}-avx2 blunder/main.go
	GOARCH=amd64 GOAMD64=v4 go build -o ${BINARY_NAME}-avx512 blunder/main.go

build-windows:
	set GOARCH=amd64&& set GOAMD64=v1&& go build -o ${BINARY_NAME}-default.exe blunder/main.go
	set GOARCH=amd64&& set GOAMD64=v2&& go build -o ${BINARY_NAME}-popcnt.exe blunder/main.go
	set GOARCH=amd64&& set GOAMD64=v3&& go build -o ${BINARY_NAME}-avx2.exe blunder/main.go
	set GOARCH=amd64&& set GOAMD64=v4&& go build -o ${BINARY_NAME}-avx512.exe blunder/main.go

build-all:
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows.exe blunder/main.go
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin blunder/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux blunder/main.go

build-all-windows:
	set GOARCH=amd64&& set GOOS=windows&& go build -o ${BINARY_NAME}-windows.exe blunder/main.go
	set GOARCH=amd64&& set GOOS=darwin&& go build -o ${BINARY_NAME}-darwin blunder/main.go
	set GOARCH=amd64&& set GOOS=linux&& go build -o ${BINARY_NAME}-linux blunder/main.go

clean-all:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

clean-all-windows:
	go clean
	del /f ${BINARY_NAME}-darwin.exe
	del /f ${BINARY_NAME}-linux.exe
	del /f ${BINARY_NAME}-windows.exe

clean-build:
	go clean
	rm ${BINARY_NAME}-default
	rm ${BINARY_NAME}-popcnt
	rm ${BINARY_NAME}-avx2
	rm ${BINARY_NAME}-avx512

clean-build-windows:
	go clean
	del /f ${BINARY_NAME}-default.exe
	del /f ${BINARY_NAME}-popcnt.exe
	del /f ${BINARY_NAME}-avx2.exe
	del /f ${BINARY_NAME}-avx512.exe