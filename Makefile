install_oapi_codegen:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

generate_api_types:
	oapi-codegen -package api -generate types api/swagger.yml > api/types.gen.go

generate_api_server:
	oapi-codegen -package api -generate gorilla api/swagger.yml > api/server.gen.go

generate_api_spec:
	oapi-codegen -package api -generate spec api/swagger.yml > api/spec.gen.go