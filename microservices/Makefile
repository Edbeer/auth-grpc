# migrations
migrate-create:
	migrate create -ext sql -dir ./migrations -seq ${NAME}

# TODO env Port, db_name
migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/acc_db?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/acc_db?sslmode=disable" down

# cmd
run:
	go run cmd/api/main.go

# tests
test:
	go test -cover -race ./...

# buf gen
lint:
	cd proto && buf lint

gen: lint
	cd proto && buf generate

# docker
docker-build:
	docker build -t edbeer/acc:v1 .

docker-run:
	docker run -it --rm --name account edbeer/acc:v1
