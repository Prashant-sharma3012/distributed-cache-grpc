package server

import (
	"net/http"
	"sync"
	"time"
)

type Req struct {
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	KeyToDelete string      `json:"keyToDelete"`
}

type Worker struct {
	sync.Mutex
	Id       int
	KeyCount int
	Addr     string
}

type CacheIndexRecord struct {
	WorkerId   int
	Key        string
	Addr       string
	CreatedAt  time.Time
	LastUsedOn time.Time
}

type Server struct {
	Srv        *http.Server
	Handler    *http.ServeMux
	Workers    *[]Worker
	CacheIndex map[string]*CacheIndexRecord
}

func InitServer(port string) *Server {
	handler := http.NewServeMux()

	return &Server{
		Srv: &http.Server{
			Addr:    port,
			Handler: handler,
		},
		Handler:    handler,
		CacheIndex: make(map[string]*CacheIndexRecord),
	}
}
