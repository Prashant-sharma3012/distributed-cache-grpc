To run locally
go run main.go -serverPort=3000 -numOfWorkers=3 -workerPortStartAt=4000

To build Proto
protoc --go_out=plugins=grpc:connector connector/connector.proto