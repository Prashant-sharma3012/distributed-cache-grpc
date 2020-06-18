package manager

import (
	"fmt"
	"log"
	"strconv"

	"github.com/distributed-cache-grpc/connector"
	"github.com/distributed-cache-grpc/manager/api"
	"github.com/distributed-cache-grpc/manager/server"
	worker "github.com/distributed-cache-grpc/worker"
	"google.golang.org/grpc"
)

func StartCacheServer(serverPortPtr int, numOfWorkers int, workerPortStartAtPtr int) {
	log.Println("Starting Manager")
	serverPortStr := strconv.Itoa(serverPortPtr)

	s := server.InitServer(":" + serverPortStr)

	log.Println("Initializing End Points")
	api.InitRoutes(s)

	log.Println("Starting workers")
	var workers []server.Worker

	for i := 0; i < numOfWorkers; i++ {
		workerPortStr := strconv.Itoa(workerPortStartAtPtr + i)
		fmt.Println("Starting Worker at " + workerPortStr)

		workers = append(workers, server.Worker{
			Id:       i,
			KeyCount: 0,
			Addr:     workerPortStr,
		})

		go func() { worker.StartWorker(":"+workerPortStr, i) }()
	}

	log.Println("Connecting to all workers")

	// Once all workers are started dial into them and makeconnections
	var conn *grpc.ClientConn
	var err error
	for i := 0; i < numOfWorkers; i++ {
		conn, err = grpc.Dial(":9000", grpc.WithInsecure())
		if err != nil {
			log.Println("Error Connecting to worker running on port %s", workers[i].Addr)
		}

		workers[i].Conn = connector.NewConnectorServiceClient(conn)
	}

	s.Workers = &workers

	fmt.Println("Initialized Server, listening on post 3000")
	log.Fatal(s.Srv.ListenAndServe())
}
