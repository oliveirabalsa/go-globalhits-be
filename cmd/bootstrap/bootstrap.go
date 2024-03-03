package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := loadEnv()
	if err != nil {
		fmt.Println("Failed to load environment variables:", err)
		return
	}

	createQueue()

	startAPIService()

	startWorkerService()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	<-signalCh

	fmt.Println("Shutting down...")
	os.Exit(0)
}

func loadEnv() error {
	cmd := exec.Command("bash", "-c", "set -a && . ./.env && set +a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createQueue() {
	fmt.Println("Creating queue...")
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
	go func() {
		cmd := exec.Command("go", "run", "cmd/api/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Failed to start API service:", err)
		}
	}()

	time.Sleep(1 * time.Second)
}

func startWorkerService() {
	go func() {
		cmd := exec.Command("go", "run", "cmd/worker/worker.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Failed to start worker service:", err)
		}
	}()

	time.Sleep(1 * time.Second)
}
