package main

import (
	"github.com/sadbanzai/sysmon/internal/collector"
	pb "github.com/sadbanzai/sysmon/internal/pb"
	"github.com/sadbanzai/sysmon/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	limit := 10

	collectorLA := collector.Collector{
		Name:     "LoadAverage",
		FuncName: collector.LoadAverage,
		Data:     make([]collector.Items, limit),
		Limit:    limit,
		Enabled:  false,
	}
	collectorCS := collector.Collector{
		Name:     "CpuStat",
		FuncName: collector.CpuStat,
		Data:     make([]collector.Items, limit),
		Limit:    limit,
		Enabled:  true,
	}

	myServer := server.Server{
		CollectorLoadAverage: collectorLA,
		CollectorCpuStat:     collectorCS,
	}

	myServer.Start()

	pb.RegisterSysmonServer(grpcServer, &myServer)

	reflection.Register(grpcServer)

	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
