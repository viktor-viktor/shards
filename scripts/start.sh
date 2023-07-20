#!/bin/bash
echo "Shutting down."
docker-compose down

echo "Starting mongo."
docker-compose up -d mongo

# Wait for MongoDB to be ready before executing the script inside the container
# TODO poll mongo if it's ready
sleep 5

echo "Creating collections."
docker exec -it mongodb_local mongosh -u user -p pwd --authenticationDatabase admin --eval 'use local;' --eval 'db.createCollection("events");' --eval 'db.createCollection("workers");'

#mongosh -u user -p pwd --authenticationDatabase admin
echo "Collections created. Starting server."
docker-compose up --build -d app