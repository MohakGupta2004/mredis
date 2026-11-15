package services

import (
	"strconv"
	"time"

	"github.com/MohakGupta2004/mredis/internal/resp"
)

func handleRPush(args []string, store *StoreType) string {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := args[0]
	list := store.getOrCreateList(key)

	for _, val := range args[1:] {
		list.RPush(val)
	}
	if len(store.waiters[key]) > 0 {
		waiter := store.waiters[key][0]
		store.waiters[key] = store.waiters[key][1:]
		values := []string{key}
		values = append(values, getListRange(list, 0, list.Size())...)
		waiter <- resp.WriteArray(values)
	}
	store.list[key] = list
	return resp.WriteInteger(list.Size())
}

func handleLPush(args []string, store *StoreType) string {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := args[0]
	list := store.getOrCreateList(key)

	for _, val := range args[1:] {
		list.LPush(val)
	}
	store.list[key] = list
	return resp.WriteInteger(list.Size())
}

func handleLRange(args []string, store *StoreType) string {
	key := args[0]
	start, _ := strconv.Atoi(args[1])
	end, _ := strconv.Atoi(args[2])
	store.mu.RLock()
	list := store.list[key]
	store.mu.RUnlock()
	size := list.Size()
	if size == 0 {
		return resp.WriteArrayHeader(0)
	}

	if start < 0 && (-start) < size {
		start = size + start
	} else if start < 0 {
		start = 0
	}
	if end < 0 && (-end) < size {
		end = size + end
	} else if end < 0 {
		end = size
	}
	resultSize := end - start + 1
	if resultSize > size {
		resultSize = size - start
	}

	res := resp.WriteArrayHeader(resultSize)
	for _, value := range getListRange(list, start, end+1) {
		res += resp.WriteBulkString(value)
	}
	return res
}

func handleLLen(args []string, store *StoreType) string {
	key := args[0]
	store.mu.RLock()
	list, exists := store.list[key]
	if !exists {
		return resp.WriteInteger(0)
	}
	store.mu.RUnlock()

	return resp.WriteInteger(list.Size())
}

func handleLPop(args []string, store *StoreType) string {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := args[0]
	list, exists := store.list[key]
	if !exists {
		return resp.WriteNullBulkString()
	}

	if len(args) > 1 {
		count, _ := strconv.Atoi(args[1])
		res := resp.WriteArrayHeader(count)
		for range count {
			elem, exists := list.LPop()
			if !exists {
				res += resp.WriteNullBulkString()
			}
			res += resp.WriteBulkString(elem)
		}
		return res
	}

	elem, exists := list.LPop()
	if !exists {
		return resp.WriteNullBulkString()
	}
	store.list[key] = list
	return resp.WriteBulkString(elem)
}

func handleBLPop(args []string, store *StoreType) string {
	store.mu.Lock()
	key := args[0]
	list := store.list[key]
	if list.Size() > 0 {
		returnValue, _ := list.LPop()
		store.list[key] = list
		return returnValue
	}

	timeout := args[1]
	result := make(chan string, 1)
	store.waiters[key] = append(store.waiters[key], result)
	store.mu.Unlock()
	if timeout == "0" {
		return <-result
	}

	parsedTimeout, _ := strconv.ParseFloat(timeout, 64)
	select {
	case value := <-result:
		return value
	case <-time.After(time.Duration(parsedTimeout * float64(time.Second))):
		store.mu.Lock()
		store.waiters[key] = store.waiters[key][1:]
		store.mu.Unlock()
		return resp.WriteNullArray()
	}
}
