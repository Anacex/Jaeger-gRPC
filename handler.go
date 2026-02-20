package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

/*for backend abstraction.
We can have multiple implementations of LogHandler,
 for example one that writes to a file,
 another that sends logs to a remote service, etc.
This allows us to change the underlying log processing
logic without affecting the gRPC server or the protobuf definitions.*/

type LogHandler interface {
	Handle(ctx context.Context, entry LogEntry) error
}

type JaegerHandler struct{}

func NewJaegerHandler() LogHandler {
	return &JaegerHandler{}
}

// tracing encapsulated inside handler, not exposed to gRPC server, this allows us to change the tracing implementation without affecting the gRPC server code
func (j *JaegerHandler) Handle(ctx context.Context, entry LogEntry) error {
	tracer := otel.Tracer("log-collector-handler")

	ctx, span := tracer.Start(ctx, "process-log")
	defer span.End()

	span.SetAttributes(
		attribute.String("service.name", entry.ServiceName),
		attribute.String("Log.message", entry.Message),
	)

	return nil
}
