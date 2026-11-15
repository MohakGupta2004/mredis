package services

import (
	"strconv"
	"strings"
	"time"

	"github.com/MohakGupta2004/mredis/internal/resp"
	"github.com/MohakGupta2004/mredis/internal/utils"
)

func handleEcho(args []string) string {
	data := args[0]
	return resp.WriteBulkString(data)
}

func handlePing() string {
	return resp.WriteSimpleString("PONG")
}

func handleSet(args []string, store *StoreType) string {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := args[0]
	value := args[1]

	expiry := time.Time{}
	if len(args) > 2 {
		ttlFormat := strings.ToLower(args[2])
		ttlDuration, _ := strconv.ParseInt(args[3], 10, 64)
		expiry = utils.GetExpectedExpiryFromCommand(ttlFormat, time.Duration(ttlDuration))
	}

	store.data[key] = Value{
		Value: value,
		End:   expiry,
	}
	return resp.WriteSimpleString("OK")
}

func handleGet(args []string, store *StoreType) string {
	key := args[0]
	store.mu.RLock()
	val, ok := store.data[key]
	store.mu.RUnlock()
	if !ok {
		return resp.WriteNullBulkString()
	}

	if !val.End.IsZero() && time.Now().After(val.End) {
		store.mu.Lock()
		delete(store.data, key)
		store.mu.Unlock()
		return resp.WriteNullBulkString()
	}

	return resp.WriteBulkString(val.Value)
}
