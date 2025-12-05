package main

import (
	"github.com/Kardbord/k8s-sandbox/internal/api"
	pb "github.com/Kardbord/k8s-sandbox/internal/gen/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	srv, err := api.NewServer()
	if err != nil {
		log.Fatalf("failed to start API server: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterJobServiceServer(grpcServer, srv)

	log.Println("API server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
