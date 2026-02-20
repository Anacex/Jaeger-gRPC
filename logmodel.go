package main

//We decouple protobuf struct from internal processing model.
//Never tie core logic directly to transport layer.

type LogEntry struct {
	ServiceName string
	Message     string
}
