.DEFAULT_GOAL := run

.PHONY:fmt vet build
fmt:
		go fmt cmd/hyde_bot/main.go
vet: fmt
		go vet cmd/hyde_bot/main.go
build: vet
		go build cmd/hyde_bot/main.go
run: vet
		go run cmd/hyde_bot/main.go $(ARGS)
 
