.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/entry entry/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/lambda-url-func lambda-url-func/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
