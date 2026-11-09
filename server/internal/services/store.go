package services

import (
	"sync"
	"time"
)

type Value struct {
	Value string
	End   time.Time
}
type StoreType struct {
	mu   sync.RWMutex
	data map[string]Value
}

func NewStore() *StoreType {
	return &StoreType{
		data: make(map[string]Value),
	}
}
