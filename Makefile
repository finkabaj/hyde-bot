.DEFAULT_GOAL := run_bot

.PHONY:fmt vet build
fmt_bot:
		go fmt cmd/hyde-bot/hyde_bot.go
vet_bot: fmt
		go vet cmd/hyde-bot/hyde_bot.go
clean:
		rm -rf dist/
build_bot: clean vet
		go build -o dist/ -v cmd/hyde-bot/hyde_bot.go
run_bot: vet_bot
		go run cmd/hyde-bot/hyde_bot.go $(ARGS)
 
