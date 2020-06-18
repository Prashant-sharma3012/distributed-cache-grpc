package main

import (
	"flag"

	"github.com/distributed-cache/manager"
)

func main() {

	serverPortPtr := flag.Int("serverPort", 42, "an int")
	numOfWorkersPtr := flag.Int("numOfWorkers", 42, "an int")
	workerPortStartAtPtr := flag.Int("workerPortStartAt", 42, "an int")

	flag.Parse()

	manager.StartCacheServer(*serverPortPtr, *numOfWorkersPtr, *workerPortStartAtPtr)
}
