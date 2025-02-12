include .env

live/rest-api:
	set -a; source .env; set +a; \
		go run cmd/app/main.go 

generate/types:
	tygo generate

db/up:
	source .env && docker compose --env-file .env -p cookbook up db -d
db/migration/up:
	docker compose --env-file .env -p cookbook up migrate
db/drop:
	docker compose --env-file .env -p cookbook up drop
db/down:
	docker compose --env-file .env -p cookbook down
db/migration/create:
	docker run -v $(PWD)/migrations:/migrations migrate/migrate -path=/migrations/ -verbose create -ext sql -dir /migrations $(name)

lint/rest-api:
	golangci-lint run

test/rest-api:
	set -a; source .env.test; set +a; \
		gotestsum 

test/db/up:
	source .env.test && docker compose --env-file .env.test -p cookbook-test up db -d
test/db/migration/up:
	source .env.test; docker compose --env-file .env.test -p cookbook-test up migrate
test/db/drop:
	source .env.test; docker compose --env-file .env.test -p cookbook-test up drop
test/db/down:
	source .env.test; docker compose --env-file .env.test -p cookbook-test down
