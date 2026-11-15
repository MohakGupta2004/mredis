package services

import "github.com/MohakGupta2004/mredis/internal/resp"

func ExecuteCommand(command string, args []string, store *StoreType) string {
	switch command {
	case "echo":
		return handleEcho(args)
	case "ping":
		return handlePing()
	case "set":
		return handleSet(args, store)
	case "get":
		return handleGet(args, store)
	case "rpush":
		return handleRPush(args, store)
	case "lpush":
		return handleLPush(args, store)
	case "lrange":
		return handleLRange(args, store)
	case "llen":
		return handleLLen(args, store)
	case "lpop":
		return handleLPop(args, store)
	case "blpop":
		return handleBLPop(args, store)
	default:
		return resp.WriteError("ERR unknown command")
	}
}
