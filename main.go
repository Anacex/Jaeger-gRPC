package main

import (
	"context"
	"log"
	"net"

	pb "log-collector/proto"

	"google.golang.org/grpc"            //grpc server
	"google.golang.org/grpc/reflection" //reflection allows clients to discover the gRPC services at runtime, useful for debugging and tools like grpcurl
)

/*
Notes:
gRPC returns immediately (Unary logs are fire, not streamed for now)

Heavy work done in background (Worker pool)

# No goroutines per request

Controlled concurrency
*/
type server struct {
	pb.UnimplementedLogServiceServer // This is required to satisfy the interface, it provides default implementations for all methods, so we only need to implement the ones we care about
	workers                          *WorkerPool
}

func (s *server) SendLog(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
	entry := LogEntry{
		ServiceName: req.ServiceName,
		Message:     req.Message,
	}

	ok := s.workers.Enqueue(entry)
	if !ok {
		log.Println("Failed to enqueue log, queue is full:", entry)
	}

	return &pb.LogResponse{Status: "OK"}, nil

}

func main() {
	handler := NewJaegerHandler()
	workerPool := NewWorkerPool(handler, 1000, 5) //buffer size 1000, 5 workers - these numbers can be tuned based on expected load and system resources
	tp, err := InitTracer()

	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background()) // Ensure that all spans are flushed before the application exits

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLogServiceServer(grpcServer, &server{workers: workerPool}) // Register the server implementation with the gRPC server, this allows the gRPC server to route incoming requests to our SendLog method

	reflection.Register(grpcServer)
	log.Println("gRPC server running on : 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
