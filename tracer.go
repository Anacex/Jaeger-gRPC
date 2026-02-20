package main

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func InitTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	endpoint := os.Getenv("OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318"
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(), //disable TLS for local testing
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // Use batch processing for better performance
		sdktrace.WithResource(resource.NewWithAttributes( // This is the resource information that will appear in jaeger UI
			semconv.SchemaURL,
			semconv.ServiceName("log-collector"),
		)),
	)

	otel.SetTracerProvider(tp) //make the tracer provider globally available, otel.Tracer(...)

	log.Println("Tracer initialized")

	return tp, nil
}
