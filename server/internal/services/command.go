package services

import "github.com/MohakGupta2004/mredis/internal/resp"

func ExecuteCommand(command string, args []string, store *StoreType) string {
	cmd := newStringCommand(args, store)
	list := newListCommand(args, store)

	switch command {
	case "echo":
		return cmd.echo()
	case "ping":
		return cmd.ping()
	case "set":
		return cmd.set()
	case "get":
		return cmd.get()
	case "rpush":
		return list.rpush()
	case "lpush":
		return list.lpush()
	case "lrange":
		return list.lrange()
	case "llen":
		return list.llen()
	case "lpop":
		return list.lpop()
	case "blpop":
		return list.blpop()
	default:
		return resp.WriteError("ERR unknown command")
	}
}
