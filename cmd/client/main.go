package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Kardbord/k8s-sandbox/internal/client"
)

const clientID = "faux-client"

func main() {
	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = "localhost:30051"
	}

	c, err := client.New(grpcAddr)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	count := rand.Intn(1000)
	log.Printf("Submitting %d jobs to %s...\n", count, grpcAddr)

	for range count {
		iterations := uint32(rand.Intn(1000000) + 500000) // simulate work

		job, err := c.CreateJob(context.Background(), clientID, iterations)
		if err != nil {
			log.Printf("CreateJob failed: %v", err)
			continue
		}

		log.Printf("Submitted job %s with %d iterations\n", job.JobId, job.Iterations)
		time.Sleep(200 * time.Millisecond)
	}

	log.Println("Done submitting jobs.")
}
