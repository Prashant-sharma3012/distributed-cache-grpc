package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var BaseUrl = "http://localhost:"

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

func (s *Server) AddToCache(w http.ResponseWriter, r *http.Request) {
	var req Req
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
	fmt.Println("Using worker" + strconv.Itoa(id) + "Running on port" + port)

	s.CacheIndex[req.Key] = &CacheIndexRecord{
		WorkerId:   id,
		Key:        req.Key,
		Addr:       port,
		CreatedAt:  time.Now(),
		LastUsedOn: time.Now(),
	}

	workerURL := BaseUrl + port + "/add"
	if replaceRecord {
		fmt.Println("Removing Old Key: " + keyToRemove)
		req.KeyToDelete = keyToRemove
		workerURL = BaseUrl + port + "/replace"
	}

	reqBody, _ := json.Marshal(req)
	resFromWorker, err1 := http.Post(workerURL, "application/json", bytes.NewBuffer(reqBody))
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(resFromWorker.Body)
	w.Write(body)
}

func (s *Server) RemoveFromCache(w http.ResponseWriter, r *http.Request) {
	var req Req
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

	// remove fromcache index
	delete(s.CacheIndex, req.Key)

	port := s.CacheIndex[req.Key].Addr
	workerURL := BaseUrl + port + "/remove"

	reqBody, _ := json.Marshal(req)
	resFromWorker, err1 := http.Post(workerURL, "application/json", bytes.NewBuffer(reqBody))
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(resFromWorker.Body)
	w.Write(body)
}

func (s *Server) GetFromCache(w http.ResponseWriter, r *http.Request) {
	var req Req
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

	port := s.CacheIndex[req.Key].Addr
	workerURL := BaseUrl + port + "/get"

	reqBody, _ := json.Marshal(req)
	resFromWorker, err1 := http.Post(workerURL, "application/json", bytes.NewBuffer(reqBody))
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	body, _ := ioutil.ReadAll(resFromWorker.Body)
	w.Write(body)
}
