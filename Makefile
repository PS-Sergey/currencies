install_migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

install_oapi_codegen:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

generate_api_types:
	oapi-codegen -package api -generate types api/swagger.yml > api/types.gen.go

generate_api_server:
	oapi-codegen -package api -generate gorilla api/swagger.yml > api/server.gen.go

generate_api_spec:
	oapi-codegen -package api -generate spec api/swagger.yml > api/spec.gen.go

migration_up:
	migrate -path migrations/postgres/ -database "postgresql://postgres:postgres@localhost:5432/currency?sslmode=disable" -verbose up

migration_down:
	migrate -path migrations/postgres/ -database "postgresql://postgres:postgres@localhost:5432/currency?sslmode=disable" -verbose down