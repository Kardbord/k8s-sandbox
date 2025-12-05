package db

import (
	"context"
	"time"

	pb "github.com/Kardbord/k8s-sandbox/internal/gen/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	dbURL := "postgres://postgres:password@localhost:5432/jobs?sslmode=disable"
	return pgxpool.New(ctx, dbURL)
}

func InsertJob(ctx context.Context, pool *pgxpool.Pool, job *pb.Job) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO jobs (job_id, client_id, iterations, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		job.JobId,
		job.ClientId,
		job.Iterations,
		int32(job.Status),
		job.CreatedAt.AsTime(),
		job.UpdatedAt.AsTime(),
	)
	return err
}

func GetJob(ctx context.Context, pool *pgxpool.Pool, jobID string) (*pb.Job, error) {
	row := pool.QueryRow(ctx,
		`SELECT client_id, job_id, iterations, status, created_at, updated_at
		 FROM jobs WHERE job_id = $1`,
		jobID,
	)

	var (
		clientID   string
		jID        string
		iterations uint32
		statusInt  int32
		createdAt  time.Time
		updatedAt  time.Time
	)

	if err := row.Scan(&clientID, &jID, &iterations, &statusInt, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	return &pb.Job{
		ClientId:   clientID,
		JobId:      jID,
		Iterations: iterations,
		Status:     pb.JobStatus(statusInt),
		CreatedAt:  timestamppb.New(createdAt),
		UpdatedAt:  timestamppb.New(updatedAt),
	}, nil
}

func UpdateJobStatus(ctx context.Context, pool *pgxpool.Pool, jobID string, status pb.JobStatus, updated *timestamppb.Timestamp) error {
	_, err := pool.Exec(ctx,
		`UPDATE jobs
		 SET status = $1, updated_at = $2
		 WHERE job_id = $3`,
		int32(status),
		updated.AsTime(),
		jobID,
	)
	return err
}
