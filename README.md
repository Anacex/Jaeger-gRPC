Jaeger gRPC Log Collector

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
