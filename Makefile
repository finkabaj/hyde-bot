.DEFAULT_GOAL := run_bot_rmcmd

fmt_bot:
		go fmt cmd/hyde-bot/hyde_bot.go
vet_bot: fmt_bot
		go vet cmd/hyde-bot/hyde_bot.go
clean:
		rm -rf dist/
build_bot: clean vet_bot
		go build -o dist/ -v cmd/hyde-bot/hyde_bot.go
run_bot_rmcmd: vet_bot
		go run cmd/hyde-bot/hyde_bot.go --rmcmd=true $(ARGS)
run_bot: vet_bot
		go run cmd/hyde-bot/hyde_bot.go $(ARGS)

fmt_api:
		go fmt cmd/api/api.go
vet_api: fmt_api
		go vet cmd/api/api.go
build_api: clean vet_api
		go build -o dist/ -v cmd/api/api.go
run_api: vet_api
		go run cmd/api/api.go $(ARGS)

test: 
	go test ./... -v
