package api

import (
	"github.com/distributed-cache-grpc/manager/server"
)

type Api struct {
	Srv *server.Server
}

var API *Api

func InitRoutes(s *server.Server) {
	API = &Api{
		Srv: s,
	}

	API.Srv.Handler.HandleFunc("/add", s.AddToCache)
	API.Srv.Handler.HandleFunc("/remove", s.RemoveFromCache)
	API.Srv.Handler.HandleFunc("/get", s.GetFromCache)
}
