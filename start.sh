#!/bin/bash
echo "Creating queue"
curl -X PUT http://localhost:15672/api/queues/%2F/globalhitss -u globalhitss:globalhitss -H "Content-Type: application/json" -d '{"auto_delete":false,"durable":true,"arguments":{}}'

echo "Starting API service..."
go run cmd/api/main.go &

echo "Starting worker service..."
go run cmd/worker/worker.go
