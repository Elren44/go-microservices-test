COMPOSE_FILES = -f docker-compose.yml -f docker-compose.gateway.yml -f docker-compose.auth.yml
DOCKER_COMPOSE = docker-compose $(COMPOSE_FILES)

.PHONY: build run down

build:
	$(DOCKER_COMPOSE) build

run:
	$(DOCKER_COMPOSE) up

down:
	$(DOCKER_COMPOSE) down
