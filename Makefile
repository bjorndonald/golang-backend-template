project_name = golang-backend-template
image_name = golang-backend-template:latest

include .env
export

POSTGRES_IMAGE_NAME = postgres:latest
POSTGRES_CONTAINER_NAME = postgres-db

run-local:
	go fmt ./backend/... && gosec ./backend/... && air main.go

docs-generate:
	swag init

requirements:
	go mod tidy
	cd frontend
	npm install

clean-packages:
	go clean -modcache

build-backend:
	docker build -t $(image_name) ./backend

run-frontend:
	cd frontend; npm run dev

stop-postgres:
	docker stop $(POSTGRES_CONTAINER_NAME)
	docker rm $(POSTGRES_CONTAINER_NAME)

start-postgres:
	docker run --name $(POSTGRES_CONTAINER_NAME) --env-file .env -d -p 5432:5432 $(POSTGRES_IMAGE_NAME)

start:
	make start-postgres
	docker run -d -p 8000:8000 --env-file .env --link $(POSTGRES_CONTAINER_NAME):postgres $(image_name)
	cd ./frontend
	bun install
	npm run dev

build-no-cache:
	docker build --no-cache -t $(image_name) .

service-stop:
	docker compose down

service-start:
	make service-stop
	docker compose up

reset-data:
	rm -r postgres-data

integration-test:
	make start
	go clean -testcache && go test -v ./integration-test/...
