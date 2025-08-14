start-postgres:
	 docker compose --env-file config/config.env up -d postgres

start-redis:
	 docker compose --env-file config/config.env up -d redis

down:
	docker compose --env-file config/config.env down

all: start-postgres start-redis