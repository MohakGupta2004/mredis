package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MohakGupta2004/mredis/internal/utils"
)

func ExecuteCommand(command string, args []string, store *StoreType) string {
	switch command {
	case "echo":
		data := args[4]
		return fmt.Sprintf("$%v\r\n%s\r\n", len(data), data)
	case "ping":
		return "+PONG\r\n"
	case "set":
		store.mu.Lock()
		defer store.mu.Unlock()
		key := args[4]
		value := args[6]
		if len(args) > 8 {
			ttlFormat := strings.ToLower(args[8])
			inputTtl := args[10]
			ttlDuration, _ := strconv.ParseInt(inputTtl, 10, 64)
			expectedDuration := utils.GetExpectedExpiryFromCommand(ttlFormat, time.Duration(ttlDuration))
			store.data[key] = Value{
				Value: value,
				End:   expectedDuration,
			}
		} else {
			store.data[key] = Value{
				Value: value,
				End:   time.Time{},
			}
		}
		return "+OK\r\n"
	case "get":
		store.mu.RLock()
		key := args[4]
		val, ok := store.data[key]
		store.mu.RUnlock()
		if !ok {
			return "$-1\r\n"
		}

		if !val.End.IsZero() && time.Now().After(val.End) {
			store.mu.Lock()
			delete(store.data, key)
			store.mu.Unlock()
			return "$-1\r\n"
		}

		return fmt.Sprintf("$%v\r\n%s\r\n", len(val.Value), val.Value)
	default:
		return "-ERR unknown command\r\n"
	}
}
