# Makefile

app_init:
	docker-compose up

redis_index:
	./scripts/redis.sh

.PHONY: app_init redis_index