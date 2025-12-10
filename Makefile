.PHONY: build
build:
	go build -o server/myserver server/main.go
	go build -o client/myclient client/main.go

.PHONY: run
run:
	./client/myclient
