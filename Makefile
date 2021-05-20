OUTPUT = main 
SERVICE_NAME = books-api

.PHONY: test
test:
	go test ./...

build-local:
	go build -o $(OUTPUT) ./cmd/$(SERVICE_NAME)/main.go

run: build-local
	@echo ">> Running application ..."
	MONGO_PORT=9000 \
	MONGO_DB_NAME=test \
	MONGO_URL=mongodb://localhost:27017 \
	JWT_ACCESS_SECRET=T52N6pRxNZDW45UR \
	./$(OUTPUT)