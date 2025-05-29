build:
	@go build -o bin/app main.go
	
run: build
	@./bin/app $(ARGS)

test:
	@go test -v ./...

.PHONY: build