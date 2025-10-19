.PHONY: build 

run:build
	./build/app
build:
	go mod tidy
	go build -o build/app main.go
