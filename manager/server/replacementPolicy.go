package server

import "time"

var MAX_CACHE_RECORDS = 12

// LeastRecentlyUsed function use least recently used policy to replace an old cache record with a new one
// It will remove the record from cache index and hit
// remove endpoint on the corresponding worker
func (s *Server) LeastRecentlyUsed() (string, string, int) {
	keyToRemove := ""
	workerId := 0
	addr := ""
	referenceTime := time.Now()

	for key, record := range s.CacheIndex {
		if referenceTime.After(record.LastUsedOn) {
			referenceTime = record.LastUsedOn
			keyToRemove = key
			workerId = record.WorkerId
			addr = record.Addr
		}
	}

	// Remove record from server
	delete(s.CacheIndex, keyToRemove)

	return keyToRemove, addr, workerId
}

// IsLimitReached checks if cache record limit is reached or not
// For Brevity and testing lets keep max cache records tp 12
func (s *Server) IsLimitReached() bool {
	count := 0

	for _, worker := range *s.Workers {
		count = count + worker.KeyCount
	}

	if count >= MAX_CACHE_RECORDS {
		return true
	}

	return false
}
