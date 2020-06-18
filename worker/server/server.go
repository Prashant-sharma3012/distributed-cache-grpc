package server

import (
	"log"
	"net"

	"github.com/distributed-cache-grpc/connector"
	grpc "google.golang.org/grpc"
)

type Req struct {
	Key         string
	Value       interface{}
	KeyToDelete string
}

type Server struct {
	Id    int
	Port  string
	Cache map[string]interface{}
}

func InitServer(port string, id int) {
	lis, _ := net.Listen("tcp", port)

	s := Server{
		Id:    id,
		Port:  port,
		Cache: make(map[string]interface{}),
	}

	grpcServer := grpc.NewServer()

	connector.RegisterConnectorServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc %v", err)
	}

	log.Printf("Initialized Worker, listening on %s" + port)
}
