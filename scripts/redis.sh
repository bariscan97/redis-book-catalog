#!/bin/bash

source .env
until docker exec redisearch redis-cli -h localhost -p 6379 PING | grep -q PONG; do
  sleep 1
  echo "Waiting for redis to be ready..."
done

echo "REDIS IS READY"

docker exec -it redisearch redis-cli FT.CREATE idx:books ON HASH PREFIX 1 "book:" SCHEMA title TEXT SORTABLE author TEXT SORTABLE category TEXT SORTABLE price NUMERIC SORTABLE created_at NUMERIC SORTABLE


