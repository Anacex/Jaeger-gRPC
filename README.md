Jaeger gRPC Log Collector
Overview

This project implements a production-style gRPC-based Log Collector written in Go.

The service receives log entries from multiple microservices and forwards them to Jaeger using OpenTelemetry (OTLP).

It is designed to:

Handle multiple microservices

Support medium to high load

Run inside Kubernetes

Provide asynchronous, non-blocking log ingestion

Remain backend-pluggable via interface-based design

High-Level Architecture

Microservices → gRPC Log Collector → Worker Pool → OpenTelemetry → Jaeger

Flow

Client sends log via Unary gRPC call.

Collector immediately returns OK.

Log is pushed into an internal buffered queue.

Worker pool processes logs asynchronously.

OpenTelemetry exports spans to Jaeger via OTLP.

Traces are visualized in Jaeger UI.

Key Design Decisions
1. Asynchronous Processing

The gRPC server does not process logs directly.
Logs are pushed into a buffered channel and handled by a worker pool.

Benefits:

Zero-latency impact on clients

Controlled concurrency

Fault isolation

2. Worker Pool Architecture

Configurable buffer size

Configurable number of workers

Non-blocking enqueue

Drop strategy when queue is full

Panic recovery inside workers

This prevents:

Memory explosion

Server crash due to panic

Backend slowness affecting ingestion

3. Backend Abstraction

A LogHandler interface abstracts the export layer.

Current implementation:

JaegerHandler (via OpenTelemetry OTLP)

Future possible handlers:

Loki

Kafka

File storage

External monitoring pipeline

No changes to gRPC server logic are required to swap backend.

Technology Stack

Go 1.24+

gRPC

Protocol Buffers

OpenTelemetry SDK

OTLP HTTP Exporter

Jaeger (all-in-one container)

Docker

Ports
Application

50051 → gRPC Server

Jaeger

16686 → Web UI
4318 → OTLP HTTP endpoint

Jaeger UI:
http://localhost:16686

Environment Variables

The OTLP endpoint is configurable:

OTLP_ENDPOINT

Default:
localhost:4318

Example for Kubernetes:
jaeger.monitoring.svc.cluster.local:4318

Running Locally

Start Jaeger:

docker run -d --name jaeger
-p 16686:16686
-p 4318:4318
jaegertracing/all-in-one:latest

Run the service:

go run .

Send test request:

grpcurl -plaintext
-d '{"service_name":"test-service","message":"hello world"}'
localhost:50051
logcollector.LogService/SendLog

Open Jaeger:

http://localhost:16686

Select service:
log-collector

Click "Find Traces".

Current Status

Implemented:

Unary gRPC ingestion

Buffered internal queue

Worker pool with panic recovery

Non-blocking enqueue

Drop strategy when overloaded

Jaeger OTLP integration

Configurable OTLP endpoint

Next Roadmap

Streaming RPC support

Prometheus metrics for queue monitoring

Retry with exponential backoff

mTLS support

Multi-stage Dockerfile

Kubernetes deployment manifests

Horizontal scaling support

Intended Deployment

Designed for:

Kubernetes-based microservices

Centralized monitoring cluster

Medium to high traffic environments

Repository Structure

main.go → gRPC server
tracer.go → OpenTelemetry setup
handler.go → LogHandler interface + Jaeger handler
worker.go → Worker pool implementation
logmodel.go → Internal log model
proto/ → Protobuf definitions
