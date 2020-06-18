package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/distributed-cache-grpc/connector"
)

func (s *Server) AddToCache(ctx context.Context, r *connector.Request) (*connector.Response, error) {
	key := r.GetKey()
	s.Cache[key] = r
	return &connector.Response{Body: "Key Added Successfuly"}, nil
}

func (s *Server) ReplaceInCache(ctx context.Context, req *connector.Request) (*connector.Response, error) {
	// Remove record from server
	keyToDelete := req.GetKeyToDelete()
	key := req.GetKey()

	if keyToDelete != "" {
		fmt.Println("Removing Key From cache: " + keyToDelete)
		delete(s.Cache, keyToDelete)
	}

	s.Cache[key] = req
	return &connector.Response{Body: "Key Replaced Successfuly"}, nil
}

func (s *Server) RemoveFromCache(ctx context.Context, req *connector.Request) (*connector.Response, error) {
	key := req.GetKey()

	_, ok := s.Cache[key]
	if !ok {
		return nil, errors.New("Key Not Found")
	}

	delete(s.Cache, key)

	return &connector.Response{Body: "Key Removed Successfuly"}, nil
}

func (s *Server) GetFromCache(ctx context.Context, req *connector.Request) (*connector.Response, error) {
	key := req.GetKey()

	val, ok := s.Cache[key]
	if !ok {
		return nil, errors.New("Key Not Found")
	}

	res, _ := json.Marshal(val)

	return &connector.Response{Body: string(res)}, nil
}
