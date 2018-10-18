package main

import "sync"

// Counter adds a value and returns a new value
type Counter interface {
	Add(addRequest) addResponse
}

type countService struct {
	v  int
	mu sync.Mutex
}

func (c *countService) Add(v addRequest) addResponse {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v += v.V
	return addResponse{c.v, ""}
}
