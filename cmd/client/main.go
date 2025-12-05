package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Kardbord/k8s-sandbox/internal/client"
)

const (
	addr     = "localhost:50051"
	clientID = "faux-client"
)

func main() {
	count := rand.Intn(1000)

	c, err := client.New(addr)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	fmt.Printf("Submitting %d jobs to %s...\n", count, addr)

	for range count {
		iterations := uint32(rand.Intn(1000000) + 500000) // simulate work

		job, err := c.CreateJob(context.Background(), clientID, iterations)
		if err != nil {
			log.Printf("CreateJob failed: %v", err)
			continue
		}

		fmt.Printf("Submitted job %s with %d iterations\n", job.JobId, job.Iterations)
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Done submitting jobs.")
}
