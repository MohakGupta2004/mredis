package services

import (
	"sync"
	"time"

	"github.com/xgzlucario/quicklist"
)

type Value struct {
	Value string
	End   time.Time
}

type StoreType struct {
	mu      sync.RWMutex
	data    map[string]Value
	list    map[string]quicklist.QuickList
	waiters map[string][]chan string
}

func NewStore() *StoreType {
	return &StoreType{
		data:    make(map[string]Value),
		list:    make(map[string]quicklist.QuickList),
		waiters: make(map[string][]chan string),
	}
}

// getOrCreateList expects the caller to already hold the store lock.
func (store *StoreType) getOrCreateList(key string) quicklist.QuickList {
	list, exists := store.list[key]
	if !exists {
		list = *quicklist.New()
	}
	return list
}

// getListRange returns the values from index start up to (not including) end.
func getListRange(list quicklist.QuickList, start int, end int) []string {
	values := []string{}
	list.Range(start, end, func(item []byte) (stop bool) {
		values = append(values, string(item))
		return false
	})
	return values
}
