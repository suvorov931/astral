start-postgres:
	 docker compose --env-file config/config.env up -d postgres

start-redis:
	 docker compose --env-file config/config.env up -d redis

start-app:
	docker compose --env-file config/config.env up -d astral-service

down:
	docker compose --env-file config/config.env down

all: start-postgres start-redis start-app