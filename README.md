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
  - [Metrics Server](https://github.com/kubernetes-sigs/metrics-server)
- Protobuf
  - `protoc`
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`

---

## Usage

Before trying to deploy the job queue app in Kubernetes, start by verifying
that the app components work as expected as regular user processes.

1. Start Redis and the database by running `make start-db start-redis`.
2. Build the app locally by running `make build`.
3. Start the API and service worker.
    1. `./build/api`
    2. `./build/worker`
4. Start up a client by running `./build/client`.

You should see the client submitting jobs to the API, which in turn submits them
to the database. From there, the worker pulls the job IDs from the Redis FIFO
queue, processes them, and updates the job status in the database. You can
examine the contents of the database with several utility `make` targets.

- `make view-db`
- `make view-db-finished`
- `make view-db-unfinished`

If everything looks like it's working you're ready to move on to deploying
the app in Kubernetes. Stop your running processes with `Ctrl+c` and get
back to a clean starting slate by running `make clean`.

Now we're ready to try it in Kubernetes. Deploying the app sandbox is as
simple as running `make all`. If this succeeds, your app should now be
locally deployed. Check that your pods are running with `kubectl get pods`.
You should see at least one `api-*` and one `worker-*` pod. You can check
the logs of your pods using several utility `make` targets.

- `make view-api-logs`
- `make view-worker-logs`
- `make watch-api-logs`
- `make watch-worker-logs`

If your API and service worker(s) are up and ready, you can start a client
with `./build/client`. You should see your app running as before, but this
time it's in Kubernetes! You should be able to watch Kubernetes scale the
app based on CPU load.

- `make watch-pods`
- `for _ in $(seq 1 100); do ./build/client &>/dev/null & done`

You should be able to see the number of worker pods increase or decrease
based on load.
