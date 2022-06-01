.PHONY: swag app

swag:
	~/go/bin/swag init --parseDependency --parseInternal

app:
	docker-compose up -d --no-deps --build app
