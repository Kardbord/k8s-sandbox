# k8s-sandbox

A lightweight Kubernetes sandbox for testing **gRPC jobs**, worker scaling,
and stateful services like **Postgres** and **Redis**.

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
│ ├── api # gRPC API server
│ ├── client # Faux client to submit jobs
│ └── worker # Worker pod to consume jobs
├── deployments
│ ├── api
│ ├── worker
│ ├── db # Postgres deployments and config
│ └── redis # Redis deployments and config
├── internal
│ ├── client # gRPC client wrapper
│ ├── db # Postgres helper functions
│ ├── redis # Redis helper functions
│ └── gen # Generated protobuf code
├── proto # Protobuf definitions
├── Makefile
└── README.md
```

---

## Features

- Submit jobs with random CPU workloads via gRPC
- Store job metadata in Postgres
- Queue jobs in Redis
- Worker pods consume jobs and update status
- Kubernetes-native design: deploy API, worker, Redis, and Postgres pods
- Horizontal Pod Autoscaling support for worker pods

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

### 1. Build the project

```bash
make clean all
```

### 2. TODO: Deploy db, redis, api, worker, and client

### 3. TODO: Watch it work
