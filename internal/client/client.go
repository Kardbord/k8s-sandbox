package client

import (
	"context"
	"fmt"

	pb "github.com/Kardbord/k8s-sandbox/internal/gen/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient pb.JobServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server: %w", err)
	}

	return &Client{
		grpcClient: pb.NewJobServiceClient(conn),
	}, nil
}

func (c *Client) CreateJob(ctx context.Context, clientID string, iterations uint32) (*pb.Job, error) {
	req := &pb.CreateJobRequest{
		ClientId:   clientID,
		Iterations: iterations,
	}
	resp, err := c.grpcClient.CreateJob(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Job, nil
}

func (c *Client) GetJobStatus(ctx context.Context, jobID string) (*pb.Job, error) {
	req := &pb.GetJobStatusRequest{
		JobId: jobID,
	}
	resp, err := c.grpcClient.GetJobStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Job, nil
}
