package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kardbord/k8s-sandbox/internal/db"
	"github.com/Kardbord/k8s-sandbox/internal/redis"
	"github.com/Kardbord/k8s-sandbox/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("received signal %v; shutting down...", sig)
		cancel()
	}()

	rdb, err := redis.NewClient()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer rdb.Close()

	pool, err := db.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	w := worker.New(pool, rdb)

	if err := w.Run(ctx); err != nil {
		log.Fatalf("worker exited with error: %v", err)
	}
}
