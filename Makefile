.DEFAULT_GOAL := run bot

.PHONY:fmt vet build
fmt bot:
		go fmt cmd/hyde_bot/main.go
vet bot: fmt
		go vet cmd/hyde_bot/main.go
clean:
		rm -rf dist/
build bot: clean vet
		go build -o dist/ -v cmd/hyde_bot/main.go
run bot: vet
		go run cmd/hyde_bot/main.go $(ARGS)
 
