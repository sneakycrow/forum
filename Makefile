build:
	go build -o bin/forum pkg/*.go

start:
	go run pkg/*.go
