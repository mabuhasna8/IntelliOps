package main

import (
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	pb "github.com/mabuhasna8/IntelliOps/proto/automation/v1"

	"github.com/mabuhasna8/IntelliOps/internal/agent"
	httpapi "github.com/mabuhasna8/IntelliOps/internal/http"
	"github.com/mabuhasna8/IntelliOps/internal/run"
	"github.com/mabuhasna8/IntelliOps/internal/workflow"
)

func main() {
	// In-memory stores for skeleton.
	wfStore := workflow.NewStore()
	runStore := run.NewStore()

	// HTTP REST
	httpServer := httpapi.NewServer(wfStore, runStore)
	go func() {
		addr := ":8080"
		log.Printf("HTTP API listening on %s", addr)
		if err := http.ListenAndServe(addr, httpServer.Router()); err != nil {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// gRPC AgentService
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	agentSvc := agent.NewService()
	pb.RegisterAgentServiceServer(grpcServer, agentSvc)

	log.Println("gRPC AgentService listening on :9090")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc server error: %v", err)
	}
}

