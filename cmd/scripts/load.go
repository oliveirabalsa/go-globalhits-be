package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
)

func makeRequest(url string, client model.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	// Marshal client struct to JSON
	body, err := json.Marshal(client)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Optionally, you can read the response body here
	// responseBody, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(responseBody))
}

func main() {
	// Define the URL to test
	url := "http://localhost:8082/api/v1/clients"
	// Define the number of concurrent requests and total requests
	concurrency := 100
	totalRequests := 10000

	// Create a wait group to wait for all requests to finish
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	// Start timing
	start := time.Now()

	// Start concurrent requests
	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < totalRequests/concurrency; j++ {
				client := model.Client{
					Name:        "Name",
					LastName:    "Last Name",
					Address:     "Address",
					Contact:     "Contact",
					CPF:         fmt.Sprintf("%03d.%03d.%03d-00", i, j, j),
					DateOfBirth: "12/12/1912",
				}
				makeRequest(url, client, &wg)
			}
		}()
	}

	// Wait for all requests to finish
	wg.Wait()

	// Calculate and print the total time taken
	elapsed := time.Since(start)
	fmt.Printf("Total time taken: %s\n", elapsed)
}
