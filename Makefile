start: build
	docker compose -f ./docker-compose.yml up -d

build:
	DOCKER_BUILDKIT=1 docker build -t restapigo:latest -f ./Dockerfile .

start-dev:
	docker compose -f ./docker-compose.yml -f ./docker-compose.dev.override.yml up -d

status:
	docker compose ps -a

logs:
	docker compose logs -f 

stop:
	docker compose down --remove-orphans

down: stop