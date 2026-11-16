package services

import (
	"strconv"
	"strings"
	"time"

	"github.com/MohakGupta2004/mredis/internal/resp"
	"github.com/MohakGupta2004/mredis/internal/utils"
)

type stringCommand struct {
	args  []string
	store *StoreType
}

func newStringCommand(args []string, store *StoreType) *stringCommand {
	return &stringCommand{
		args:  args,
		store: store,
	}
}

func (c *stringCommand) echo() string {
	data := c.args[0]
	return resp.WriteBulkString(data)
}

func (c *stringCommand) ping() string {
	return resp.WriteSimpleString("PONG")
}

func (c *stringCommand) set() string {
	c.store.mu.Lock()
	defer c.store.mu.Unlock()
	key := utils.GetKey(c.args)
	value := utils.GetValue(c.args, 0)

	expiry := time.Time{}
	if len(c.args) > 2 {
		ttlFormat := strings.ToLower(utils.GetValue(c.args, 1))
		ttlDuration, _ := strconv.ParseInt(utils.GetValue(c.args, 2), 10, 64)
		expiry = utils.GetExpectedExpiryFromCommand(ttlFormat, time.Duration(ttlDuration))
	}

	c.store.data[key] = Value{
		Value: value,
		End:   expiry,
	}
	return resp.WriteSimpleString("OK")
}

func (c *stringCommand) get() string {
	key := utils.GetKey(c.args)
	c.store.mu.RLock()
	val, ok := c.store.data[key]
	c.store.mu.RUnlock()
	if !ok {
		return resp.WriteNullBulkString()
	}

	if !val.End.IsZero() && time.Now().After(val.End) {
		c.store.mu.Lock()
		delete(c.store.data, key)
		c.store.mu.Unlock()
		return resp.WriteNullBulkString()
	}

	return resp.WriteBulkString(val.Value)
}
