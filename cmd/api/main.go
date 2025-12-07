package main

import (
	"log"
	"net"
	"os"

	"github.com/Kardbord/k8s-sandbox/internal/api"
	pb "github.com/Kardbord/k8s-sandbox/internal/gen/proto"
	"google.golang.org/grpc"
)

func main() {
	srv, err := api.NewServer()
	if err != nil {
		log.Fatalf("failed to start API server: %v", err)
	}

	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":50051"
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterJobServiceServer(grpcServer, srv)

	log.Printf("API server listening on %s\n", grpcAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
