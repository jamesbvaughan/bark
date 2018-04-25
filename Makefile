all: linux

linux:
	go build -o bark cmd/bark/main.go
