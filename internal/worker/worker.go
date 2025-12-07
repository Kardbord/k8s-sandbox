package worker

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/Kardbord/k8s-sandbox/internal/db"
	pb "github.com/Kardbord/k8s-sandbox/internal/gen/proto"
	rdsutil "github.com/Kardbord/k8s-sandbox/internal/redis"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	blpopTimeout = 5 * time.Second
)

type Worker struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func New(pool *pgxpool.Pool, rdb *redis.Client) *Worker {
	return &Worker{
		db:  pool,
		rdb: rdb,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	log.Println("worker started; waiting for jobs...")

	for {
		if ctx.Err() != nil {
			log.Println("context canceled; worker shutting down")
			return nil
		}

		// Block waiting for job ID
		jobID, err := rdsutil.PopJob(ctx, w.rdb, blpopTimeout)
		if err != nil {
			if err == redis.Nil {
				continue // timeout, queue empty
			}
			log.Printf("PopJob error: %v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if err := w.handleJob(ctx, jobID); err != nil {
			log.Printf("handleJob(%s) error: %v", jobID, err)
		}
	}
}

func (w *Worker) handleJob(ctx context.Context, jobID string) error {
	job, err := db.GetJob(ctx, w.db, jobID)
	if err != nil {
		return err
	}

	if err := db.UpdateJobStatus(ctx, w.db, jobID, pb.JobStatus_JOB_STATUS_IN_PROGRESS, timestamppb.Now()); err != nil {
		return err
	}

	iterations := job.Iterations
	log.Printf("processing job %s (%d iterations)", jobID, iterations)

	var result float64
	for i := range iterations {
		result = math.Sqrt(float64(i))
	}
	_ = result

	if err := db.UpdateJobStatus(
		ctx, w.db, jobID, pb.JobStatus_JOB_STATUS_DONE, timestamppb.Now(),
	); err != nil {
		return err
	}

	log.Printf("job %s completed", jobID)
	return nil
}
