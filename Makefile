OUTPUT = main 
VERSION = 0.1
SERVICE_NAME = books-api

build-local:
	go build -o $(OUTPUT) ./cmd/$(SERVICE_NAME)/main.go

run: build-local
	@echo ">> Running application ..."
	RRS_PORT=7565 \
	MONGO_DB_NAME=test \
	MONGO_URL=mongodb://localhost:27017 \
	JWT_ACCESS_SECRET=T52N6pRxNZDW45UR \
	./$(OUTPUT)