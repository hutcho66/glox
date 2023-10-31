BINARY_NAME=glox

run:
	go run src/cmd/glox.go

clean:
	go clean
	rm -f ${BINARY_NAME}

build: clean
	go build -o ${BINARY_NAME} src/cmd/glox.go


install:
	go install src/cmd/glox.go

test:
	go test ./src/pkg/interpreter/interpreter_test.go -coverprofile=./cover.out -coverpkg=./src/pkg/...

coverage: test
	go tool cover -html cover.out -o cover.html