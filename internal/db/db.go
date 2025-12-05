package db

import (
	"context"

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
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		job.JobId, job.ClientId, job.Iterations, int32(job.Status),
		job.CreatedAt.AsTime(), job.UpdatedAt.AsTime())
	return err
}

func GetJob(ctx context.Context, pool *pgxpool.Pool, jobID string) (*pb.Job, error) {
	row := pool.QueryRow(ctx,
		`SELECT client_id, job_id, iterations, status, created_at, updated_at
		 FROM jobs WHERE job_id=$1`, jobID)

	var job pb.Job
	var statusInt int32
	var createdAt, updatedAt any

	if err := row.Scan(&job.ClientId, &job.JobId, &job.Iterations, &statusInt, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	job.Status = pb.JobStatus(statusInt)
	job.CreatedAt = createdAt.(*timestamppb.Timestamp)
	job.UpdatedAt = updatedAt.(*timestamppb.Timestamp)

	return &job, nil
}
