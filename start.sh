#!/bin/bash

set -a
source .env
set +a

queueExists=$(curl -s -o /dev/null -w "%{http_code}" -X GET "http://${RABBITMQ_HOST}:${RABBITMQ_API_PORT}/api/queues/%2F/${RABBITMQ_QUEUE}" -u "${RABBITMQ_USER}:${RABBITMQ_PASSWORD}")
if [ $queueExists -eq 404 ]; then
  echo "Creating queue"
  curl -X PUT "http://${RABBITMQ_HOST}:${RABBITMQ_API_PORT}/api/queues/%2F/${RABBITMQ_QUEUE}" \
    -u "${RABBITMQ_USER}:${RABBITMQ_PASSWORD}" \
    -H "Content-Type: application/json" \
    -d '{"auto_delete":false,"durable":true,"arguments":{}}'
else
  echo "Queue already exists"
fi

echo "Starting API service..."
go run cmd/api/main.go &

echo "Starting worker service..."
go run cmd/worker/worker.go
