package worker

import (
	"fmt"

	"github.com/distributed-cache-grpc/worker/server"
)

func StartWorker(port string, id int) {
	fmt.Println("Starting Worker")
	server.InitServer(port, id)
}
