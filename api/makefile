.PHONY: client
client:
ifeq ($(OS),Windows_NT)
	@echo "Please generate ../pkg/client/ on macOS or Linux, currently unsupported on windows."
else
# Versions from https://github.com/OpenAPITools/openapi-generator/releases
	@chmod +x ./openapi-generator
	@rm -rf ../pkg/client/
	OPENAPI_GENERATOR_VERSION=4.3.0 ./openapi-generator generate --package-name client -i api.yml -g go -o ../pkg/client/
	rm -rf ../pkg/client/go.mod ../pkg/client/go.sum ../pkg/client/api/
	go fmt ../...
	go test ../pkg/client/
endif
