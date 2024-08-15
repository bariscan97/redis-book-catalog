#!/bin/bash

source .env

winpty docker-compose up -d

echo "Postgres container starting..."

until winpty docker exec ${IMAGE_NAME} pg_isready -U $DB_USER; do
  sleep 1
  echo "Waiting for Postgres to be ready..."
done

echo "Postgres is ready!"

winpty docker exec -i ${IMAGE_NAME} psql -U $DB_USER -d postgres -c "CREATE EXTENSION IF NOT EXISTS dblink;"

winpty docker exec -i ${IMAGE_NAME} psql -U $DB_USER -d postgres -c "

DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME') THEN
      PERFORM dblink_exec('dbname=' || current_database(), 'CREATE DATABASE $DB_NAME');
   END IF;
END
\$\$;
"
echo "Database todoapp created"

winpty docker exec -it ${IMAGE_NAME} psql -U $DB_USER -d $DB_NAME -c '

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS books (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    author VARCHAR(100) NOT NULL,
    category VARCHAR(100) NOT NULL,
    price VARCHAR(100) NOT NULL
);
'
echo "Tables  created"

until docker exec redisearch redis-cli -h localhost -p 6379 PING | grep -q PONG; do
  sleep 1
  echo "Waiting for redis to be ready..."
done

echo "REDIS IS READY BABY"

winpty docker exec -it redisearch redis-cli FT.CREATE idx:books ON HASH PREFIX 1 "book:" SCHEMA title TEXT SORTABLE author TEXT SORTABLE category TEXT SORTABLE price NUMERIC SORTABLE


