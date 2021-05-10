package main

import "sync"

// Store is a memory store
type Store struct {
	lock sync.RWMutex
	data map[string]string
}

// NewStore returns a new store
func NewStore() *Store {
	s := &Store{
		lock: sync.RWMutex{},
		data: make(map[string]string),
	}
	return s
}
