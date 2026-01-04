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
