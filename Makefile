BINARY_NAME=iot-switch-service

build:
 GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go router.go mock.go
 GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go router.go mock.go
 GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go router.go mock.go

run: build
 ./${BINARY_NAME}

clean:
 go clean
 rm ${BINARY_NAME}-darwin
 rm ${BINARY_NAME}-linux
 rm ${BINARY_NAME}-windows