include config/config.env

help:
	@echo use make start WITH-DB={type}
	@echo type: postgres, in-memory

stop-db:
	docker compose down

start-db:
	docker compose --env-file=./config/config.env up -d
	echo "Waiting for postgres to be healthy."
	@until [ `docker inspect -f "{{.State.Health.Status}}" poster-db` = healthy ] ; do sleep 5 ; echo "ready for database..."; done
	@echo "database is ready"
	goose  -dir migrations postgres "user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable host=${POSTGRES_HOST}"  up
 
start:
ifeq ($(WITH-DB),in-memory)
	go run server.go -m true
else ifeq ($(WITH-DB),postgres)
	make start-db
	go run server.go
else
	@echo unsupported database type
endif

start-test-db:
	docker compose -f docker-compose-test.yml down
	docker compose -f docker-compose-test.yml up -d
	echo "Waiting for postgres to be healthy."
	@until [ `docker inspect -f "{{.State.Health.Status}}" test-poster-db` = healthy ] ; do sleep 5 ; echo "ready for database..."; done
	@echo "database is ready"
	goose  -dir migrations postgres "user=pguser password=postgres dbname=test_poster sslmode=disable host=localhost" up

test: start-test-db
	go test ./...
	docker compose -f docker-compose-test.yml down