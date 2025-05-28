package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sangketkit01/logger-service/data"
	"github.com/sangketkit01/logger-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context,req *logs.LogRequest) (*logs.LogResponse, error){
	input := req.GetLogEntry()
	
	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil{
		resp := &logs.LogResponse{Result: "failed"}
		return resp, err
	}

	// return response
	resp := &logs.LogResponse{
		Result: "logged",
	}

	return resp, nil

}

func (app *Config) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil{
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{Models: app.Models})

	log.Printf("gRPC server started on port: %s", gRpcPort)

	if err := server.Serve(listen) ; err != nil{
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}