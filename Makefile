all:
	docker compose up -d --build

down:
	docker compose down

up:
	docker compose up

build:
	docker compose build

exec:
	docker compose run app sh

swag:
	swag init -g cmd/app/main.go
