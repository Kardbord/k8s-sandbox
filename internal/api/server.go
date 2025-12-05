package api

import (
	"context"
	"fmt"
	"log"

	"github.com/Kardbord/k8s-sandbox/internal/db"
	pb "github.com/Kardbord/k8s-sandbox/internal/gen/proto"
	redisutil "github.com/Kardbord/k8s-sandbox/internal/redis"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedJobServiceServer
	redisClient *redis.Client
	dbPool      *pgxpool.Pool
}

func NewServer() (*Server, error) {
	ctx := context.Background()

	rdb, err := redisutil.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	pool, err := db.NewPostgresPool(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	return &Server{
		redisClient: rdb,
		dbPool:      pool,
	}, nil
}

func (s *Server) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	jobID := uuid.NewString()
	job := &pb.Job{
		ClientId:   req.ClientId,
		JobId:      jobID,
		Iterations: req.Iterations,
		Status:     pb.JobStatus_JOB_STATUS_PENDING,
		CreatedAt:  timestamppb.Now(),
		UpdatedAt:  timestamppb.Now(),
	}

	if err := db.InsertJob(ctx, s.dbPool, job); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert job: %v", err)
	}

	if err := redisutil.PushJob(ctx, s.redisClient, job.JobId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to push job to Redis: %v", err)
	}

	log.Printf("Created job %s for client %s (%d iterations)", job.JobId, job.ClientId, job.Iterations)
	return &pb.CreateJobResponse{Job: job}, nil
}

func (s *Server) GetJobStatus(ctx context.Context, req *pb.GetJobStatusRequest) (*pb.GetJobStatusResponse, error) {
	job, err := db.GetJob(ctx, s.dbPool, req.JobId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "job not found: %v", err)
	}
	return &pb.GetJobStatusResponse{Job: job}, nil
}
