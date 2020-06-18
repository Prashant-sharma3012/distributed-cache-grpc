package manager

import (
	"fmt"
	"log"
	"strconv"

	"github.com/distributed-cache-grpc/manager/api"
	"github.com/distributed-cache-grpc/manager/server"
	worker "github.com/distributed-cache-grpc/worker"
)

func StartCacheServer(serverPortPtr int, numOfWorkers int, workerPortStartAtPtr int) {
	fmt.Println("Starting Manager")
	serverPortStr := strconv.Itoa(serverPortPtr)

	s := server.InitServer(":" + serverPortStr)

	fmt.Println("Initializing End Points")
	api.InitRoutes(s)

	fmt.Println("Start Workers")
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

	s.Workers = &workers

	fmt.Println("Initialized Server, listening on post 3000")
	log.Fatal(s.Srv.ListenAndServe())
}
