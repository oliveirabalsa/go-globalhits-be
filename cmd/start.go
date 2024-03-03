package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	// Load environment variables
	err := loadEnv()
	if err != nil {
		fmt.Println("Failed to load environment variables:", err)
		return
	}

	// Create queue
	createQueue()

	// Start API service
	startAPIService()

	// Start worker service
	startWorkerService()
}

func loadEnv() error {
	// Run the command to load the environment variables from the .env file
	cmd := exec.Command("bash", "-c", "set -a && . ./.env && set +a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createQueue() {
	// Run curl command to create the queue
	cmd := exec.Command("curl",
		"-X", "PUT",
		fmt.Sprintf("http://%s:%s/api/queues/%%2F/%s", os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_API_PORT"), os.Getenv("RABBITMQ_QUEUE")),
		"-u", fmt.Sprintf("%s:%s", os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD")),
		"-H", "Content-Type: application/json",
		"-d", `{"auto_delete":false,"durable":true,"arguments":{}}`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to create queue:", err)
	}
}

func startAPIService() {
	// Run go run cmd/api/main.go in a new goroutine
	go func() {
		cmd := exec.Command("go", "run", "cmd/api/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Failed to start API service:", err)
		}
	}()

	// Wait for a short time to allow the API service to start
	time.Sleep(1 * time.Second)
}

func startWorkerService() {
	// Run go run cmd/worker/worker.go in a new goroutine
	go func() {
		cmd := exec.Command("go", "run", "cmd/worker/worker.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Failed to start worker service:", err)
		}
	}()

	// Wait for a short time to allow the worker service to start
	time.Sleep(1 * time.Second)
}
