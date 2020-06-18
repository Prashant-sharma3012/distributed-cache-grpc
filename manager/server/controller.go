package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/distributed-cache-grpc/connector"
)

func (s *Server) getWorkerAddress() (string, string, int, bool) {
	min := 0
	pos := 0

	// In case cache is full replace the record
	if s.IsLimitReached() {
		keyToRemove, addr, id := s.LeastRecentlyUsed()
		return keyToRemove, addr, id, true
	}

	for indx, worker := range *s.Workers {
		if indx == 0 {
			min = worker.KeyCount
			pos = indx
		} else {
			if worker.KeyCount < min {
				min = worker.KeyCount
				pos = indx
			}
		}
	}

	(*s.Workers)[pos].Lock()
	(*s.Workers)[pos].KeyCount++
	(*s.Workers)[pos].Unlock()

	return "", (*s.Workers)[pos].Addr, (*s.Workers)[pos].Id, false
}

func (s *Server) getWorkerConnection(id int) connector.ConnectorServiceClient {
	var conn connector.ConnectorServiceClient
	for _, worker := range *s.Workers {
		if worker.Id == id {
			conn = worker.Conn
			break
		}
	}

	return conn
}

func (s *Server) AddToCache(w http.ResponseWriter, r *http.Request) {
	var req Req
	var conn connector.ConnectorServiceClient

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check is key already present, error if yes
	_, ok := s.CacheIndex[req.Key]
	if ok {
		http.Error(w, "Key Already Present", http.StatusAlreadyReported)
		return
	}

	keyToRemove, port, id, replaceRecord := s.getWorkerAddress()
	conn = s.getWorkerConnection(id)

	fmt.Println("Using worker" + strconv.Itoa(id) + "Running on port" + port)

	s.CacheIndex[req.Key] = &CacheIndexRecord{
		WorkerId:   id,
		Key:        req.Key,
		Addr:       port,
		CreatedAt:  time.Now(),
		LastUsedOn: time.Now(),
	}

	val, _ := json.Marshal(req.Value)

	reqBody := &connector.Request{
		Key:         req.Key,
		Value:       string(val),
		KeyToDelete: req.KeyToDelete,
	}

	var resFromWorker *connector.Response
	var err1 error
	if replaceRecord {
		fmt.Println("Removing Old Key: " + keyToRemove)
		resFromWorker, err1 = conn.ReplaceInCache(context.Background(), reqBody)
	} else {
		resFromWorker, err1 = conn.AddToCache(context.Background(), reqBody)
	}

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(resFromWorker.Body))
}

func (s *Server) RemoveFromCache(w http.ResponseWriter, r *http.Request) {
	var req Req
	var conn connector.ConnectorServiceClient

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := s.CacheIndex[req.Key]
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	conn = s.getWorkerConnection(s.CacheIndex[req.Key].WorkerId)

	reqBody := &connector.Request{
		Key:         req.Key,
		Value:       "",
		KeyToDelete: req.Key,
	}

	resFromWorker, err1 := conn.RemoveFromCache(context.Background(), reqBody)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	// remove fromcache index
	delete(s.CacheIndex, req.Key)

	w.Write([]byte(resFromWorker.Body))
}

func (s *Server) GetFromCache(w http.ResponseWriter, r *http.Request) {
	var req Req
	var conn connector.ConnectorServiceClient

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := s.CacheIndex[req.Key]
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	s.CacheIndex[req.Key].LastUsedOn = time.Now()

	conn = s.getWorkerConnection(s.CacheIndex[req.Key].WorkerId)
	reqBody := &connector.Request{
		Key:         req.Key,
		Value:       "",
		KeyToDelete: req.KeyToDelete,
	}

	resFromWorker, err1 := conn.GetFromCache(context.Background(), reqBody)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(resFromWorker.Body))
}
