package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Model data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	inpput := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: inpput.Name,
		Data: inpput.Data,
	}

	err := l.Model.LogEntry.Insert(logEntry)
	if err != nil {
		resp := &logs.LogResponse{Result: "faled"}
		return resp, err
	}

	// return respose
	resp := &logs.LogResponse{Result: "logged"}
	return resp, nil
}

func (app *Config) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	server := grpc.NewServer()
	logs.RegisterLogServiceServer(server, &LogServer{Model: app.Models})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
