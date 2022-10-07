migrate-create:
	migrate create -ext sql -dir ./migrations -seq ${NAME}

# TODO env Port, db_name
migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/acc_db?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/acc_db?sslmode=disable" down