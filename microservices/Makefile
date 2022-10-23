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

# docker
docker-build:
	docker build -t edbeer/acc:v1 .

docker-run:
	docker run -it --rm --name account edbeer/acc:v1

# buf gen
lint:
	cd proto && buf lint

gen: lint
	cd proto && buf generate

gen-proto-auth:
	protoc --proto_path=proto --go_out=proto/api/ --go_opt=paths=source_relative \
	--go-grpc_out=proto/api/ --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=proto/api/ --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/ --openapiv2_opt=allow_merge=true,merge_file_name=auth \
	proto/account/v1/*.proto

gen-proto-ex:
	protoc --proto_path=proto --go_out=proto/api/ --go_opt=paths=source_relative \
	--go-grpc_out=proto/api/ --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=proto/api/ --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/ --openapiv2_opt=allow_merge=true,merge_file_name=example \
	proto/example/v1/*.proto
