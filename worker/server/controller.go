package server

import (
	"context"

	"github.com/distributed-cache-grpc/connector"
)

func (s *Server) AddToCache(ctx context.Context, r *connector.Request) (*connector.Response, error) {
	// s.Cache[r.Key] = r
	return &connector.Response{Body: "Key Added Successfuly"}, nil
}

func (s *Server) ReplaceInCache(ctx context.Context, req *connector.Request) (*connector.Response, error) {
	// Remove record from server
	// if req.KeyToDelete != "" {
	// 	fmt.Println("Removing Key From cache: " + req.KeyToDelete)
	// 	delete(s.Cache, req.KeyToDelete)
	// }

	// s.Cache[req.Key] = req
	return &connector.Response{Body: "Key Replaced Successfuly"}, nil
}

func (s *Server) RemoveFromCache(ctx context.Context, req *connector.Request) (*connector.Response, error) {
	// _, ok := s.Cache[req.Key]
	// if !ok {
	// 	return nil, errors.New("Key Not Found")
	// }

	// delete(s.Cache, req.Key)

	return &connector.Response{Body: "Key Removed Successfuly"}, nil
}

func (s *Server) GetFromCache(ctx context.Context, req *connector.Request) (*connector.Response, error) {

	// val, ok := s.Cache[req.Key]
	// if !ok {
	// 	return nil, errors.New("Key Not Found")
	// }

	return &connector.Response{Body: "Key Found Successfuly"}, nil
}
