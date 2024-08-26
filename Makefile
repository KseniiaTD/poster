all:
	@echo use make start WITH-DB={type}
	@echo type: postgres, in-memory

start:

ifeq ($(WITH-DB),in-memory)
	go run server.go -m true
else ifeq ($(WITH-DB),postgres)
	docker compose --env-file=./config/config.env up -d	
	go run server.go
else
	@echo unsupported database type
endif