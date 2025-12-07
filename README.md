# k8s-sandbox

A Kubernetes sandbox for testing gRPC jobs, worker scaling,
and stateful services like Postgres and Redis.

This project is a self-contained environment to experiment with:

- Go-based gRPC APIs and workers
- Job queueing with Redis
- Persistent storage with Postgres
- Kubernetes deployments, services, and autoscaling

---

## Project Structure

```txt
k8s-sandbox/
├── cmd
│ ├── api       # Stateless gRPC API server
│ ├── client    # Faux client to submit jobs
│ └── worker    # Stateless worker to consume jobs
├── deployments
│ ├── api       # API K8s deployments and config
│ ├── worker    # Worker K8s deployments and config
│ ├── db        # Postgres deployments and config
│ └── redis     # Redis deployments and config
├── internal
│ ├── client    # gRPC client wrapper
│ ├── db        # Postgres helper functions
│ ├── redis     # Redis helper functions
│ └── gen       # Generated protobuf code
├── proto       # Protobuf definitions
├── Makefile    # Makefile for managing builds and environment
└── README.md
```

---

## Prerequisites

- Go 1.25+
- Docker
- Kubernetes
  - `kubectl` CLI
- Protobuf
  - `protoc`
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`

---

## Getting Started

TODO
