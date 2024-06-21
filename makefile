build:
	go build -o bin/echoex

run: build
	bin/echoex