package services

import (
	"strconv"
	"time"

	"github.com/MohakGupta2004/mredis/internal/resp"
	"github.com/MohakGupta2004/mredis/internal/utils"
)

type listCommand struct {
	args  []string
	store *StoreType
}

func newListCommand(args []string, store *StoreType) *listCommand {
	return &listCommand{
		args:  args,
		store: store,
	}
}

func (c *listCommand) rpush() string {
	c.store.mu.Lock()
	defer c.store.mu.Unlock()
	key := utils.GetKey(c.args)
	list := c.store.getOrCreateList(key)

	for _, val := range utils.GetValues(c.args) {
		list.RPush(val)
	}
	if len(c.store.waiters[key]) > 0 {
		waiter := c.store.waiters[key][0]
		c.store.waiters[key] = c.store.waiters[key][1:]
		values := []string{key}
		values = append(values, getListRange(list, 0, list.Size())...)
		waiter <- resp.WriteArray(values)
	}
	c.store.list[key] = list
	return resp.WriteInteger(list.Size())
}

func (c *listCommand) lpush() string {
	c.store.mu.Lock()
	defer c.store.mu.Unlock()
	key := utils.GetKey(c.args)
	list := c.store.getOrCreateList(key)

	for _, val := range utils.GetValues(c.args) {
		list.LPush(val)
	}
	c.store.list[key] = list
	return resp.WriteInteger(list.Size())
}

func (c *listCommand) lrange() string {
	key := utils.GetKey(c.args)
	start, _ := strconv.Atoi(utils.GetValue(c.args, 0))
	end, _ := strconv.Atoi(utils.GetValue(c.args, 1))
	c.store.mu.RLock()
	list := c.store.list[key]
	c.store.mu.RUnlock()
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

func (c *listCommand) llen() string {
	key := utils.GetKey(c.args)
	c.store.mu.RLock()
	list, exists := c.store.list[key]
	c.store.mu.RUnlock()
	if !exists {
		return resp.WriteInteger(0)
	}

	return resp.WriteInteger(list.Size())
}

func (c *listCommand) lpop() string {
	c.store.mu.Lock()
	defer c.store.mu.Unlock()
	key := utils.GetKey(c.args)
	list, exists := c.store.list[key]
	if !exists {
		return resp.WriteNullBulkString()
	}

	if len(c.args) > 1 {
		count, _ := strconv.Atoi(utils.GetValue(c.args, 0))
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
	c.store.list[key] = list
	return resp.WriteBulkString(elem)
}

func (c *listCommand) blpop() string {
	c.store.mu.Lock()
	key := utils.GetKey(c.args)
	list := c.store.list[key]
	if list.Size() > 0 {
		value, _ := list.LPop()
		c.store.list[key] = list
		c.store.mu.Unlock()
		return resp.WriteArray([]string{key, value})
	}

	timeout := utils.GetValue(c.args, 0)
	result := make(chan string, 1)
	c.store.waiters[key] = append(c.store.waiters[key], result)
	c.store.mu.Unlock()
	if timeout == "0" {
		return <-result
	}

	parsedTimeout, _ := strconv.ParseFloat(timeout, 64)
	select {
	case value := <-result:
		return value
	case <-time.After(time.Duration(parsedTimeout * float64(time.Second))):
		c.store.mu.Lock()
		c.store.waiters[key] = c.store.waiters[key][1:]
		c.store.mu.Unlock()
		return resp.WriteNullArray()
	}
}
