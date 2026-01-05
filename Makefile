include .env
env:
	 @echo $(MODE)

dev:
	sh -c 'set -a; . ./.env; set +a; gow run cmd/api/main.go'

url=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)

migration-up:
	migrate -database "$(url)" -path ./migration/ up $(version)
	
migration-down:
	migrate -database "$(url)" -path ./migration/ down $(version)
	
migration-create:
	migrate create -ext sql -dir ./migration/ -seq $(name)

migration-force:
	migrate -database "$(url)" -path ./migration/ force $(version)

migration-version:
	migrate -database "$(url)" -path ./migration/ version
swag:
	swag init --parseDependency --parseInternal -g cmd/api/main.go --output ./docs

# Testing
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-coverage-report:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Mock generation
mock-gen:
	mockery

mock-clean:
	rm -rf internal/mocks
